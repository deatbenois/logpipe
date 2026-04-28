package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/formatter"
	"github.com/yourorg/logpipe/internal/highlight"
	"github.com/yourorg/logpipe/internal/output"
	"github.com/yourorg/logpipe/internal/ratelimit"
	"github.com/yourorg/logpipe/internal/sampler"
	"github.com/yourorg/logpipe/internal/source"
	"github.com/yourorg/logpipe/internal/transform"
	"github.com/spf13/pflag"
)

const version = "0.1.0"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		filterExprs  []string
		transformOps []string
		formatFlag   = formatter.DefaultFormat()
		sampleRate   = sampler.NewRateFlag()
		rlFlag       = ratelimit.NewRateFlag()
		maxLines     int
		colorize     bool
		showVersion  bool
		quiet        bool
	)

	flags := pflag.NewFlagSet("logpipe", pflag.ContinueOnError)
	flags.StringArrayVarP(&filterExprs, "filter", "f", nil, "filter expression (e.g. 'level=error', 'msg~timeout', 'traceId'")
	flags.StringArrayVarP(&transformOps, "transform", "t", nil, "transform operation (e.g. 'rename:msg:message', 'drop:password', 'redact:token')")
	flags.VarP(&formatFlag, "format", "o", "output format: raw, json, pretty (default: pretty)")
	flags.VarP(sampleRate, "sample", "s", "sample rate between 0 and 1 (e.g. 0.1 for 10%)")
	flags.VarP(rlFlag, "rate-limit", "r", "max lines per second (e.g. 100, 0 for unlimited)")
	flags.IntVarP(&maxLines, "max-lines", "n", 0, "stop after emitting this many lines (0 for unlimited)")
	flags.BoolVarP(&colorize, "color", "c", true, "colorize output based on log level")
	flags.BoolVarP(&quiet, "quiet", "q", false, "suppress all output (useful for testing filters)")
	flags.BoolVarP(&showVersion, "version", "v", false, "print version and exit")

	if err := flags.Parse(os.Args[1:]); err != nil {
		return err
	}

	if showVersion {
		fmt.Printf("logpipe %s\n", version)
		return nil
	}

	// Build pipeline components.
	rules, err := filter.ParseRules(filterExprs)
	if err != nil {
		return fmt.Errorf("invalid filter: %w", err)
	}
	ops, err := transform.ParseOps(transformOps)
	if err != nil {
		return fmt.Errorf("invalid transform: %w", err)
	}

	f := filter.New(rules)
	tr := transform.New(ops)
	hl := highlight.New(colorize)
	rl := ratelimit.New(rlFlag.Rate())
	sp := sampler.New(sampleRate.Rate())

	var out *output.Writer
	if quiet {
		out = output.New(io.Discard, "")
	} else {
		out = output.NewStdout()
	}
	if maxLines > 0 {
		out = output.NewLimiter(out, maxLines)
	}
	fmt := formatter.New(formatFlag.Format(), hl)

	// Build sources from remaining args (files) or stdin.
	args := flags.Args()
	var sources []*source.Source
	if len(args) == 0 {
		sources = append(sources, source.NewFromStdin())
	} else {
		for _, path := range args {
			s, err := source.NewFromFile(path)
			if err != nil {
				return fmt.Errorf("open %s: %w", path, err)
			}
			sources = append(sources, s)
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	lines := source.TailAll(ctx, sources...)
	for entry := range lines {
		if !sp.Sample(entry.Line) {
			continue
		}
		line := tr.Apply(entry.Line)
		if !f.Match(line) {
			continue
		}
		if err := rl.Wait(ctx); err != nil {
			break
		}
		if err := fmt.Write(out, entry.Source, line); err != nil {
			if errors.Is(err, output.ErrLimitReached) {
				break
			}
			return err
		}
	}

	return out.Close()
}
