package source

import (
	"context"
	"sync"
)

// FanIn merges multiple Line channels into a single channel.
// The returned channel is closed once all input channels are drained.
func FanIn(ctx context.Context, channels ...<-chan Line) <-chan Line {
	out := make(chan Line)
	var wg sync.WaitGroup

	forward := func(ch <-chan Line) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case line, ok := <-ch:
				if !ok {
					return
				}
				select {
				case out <- line:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go forward(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// TailAll starts tailing all provided Readers and fans their output into one channel.
func TailAll(ctx context.Context, readers []*Reader) <-chan Line {
	channels := make([]<-chan Line, len(readers))
	for i, r := range readers {
		channels[i] = r.Tail(ctx)
	}
	return FanIn(ctx, channels...)
}
