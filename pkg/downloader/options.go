package downloader

import (
	"io"
)

func WithOutputDir(dir string) Option {
	return func(dl *Downloader) {
		dl.outputDir = dir
	}
}

func WithYTFormat(format string) Option {
	return func(dl *Downloader) {
		dl.ytFormat = format
	}
}

func WithYTArguments(args []string) Option {
	return func(dl *Downloader) {
		dl.ytArgs = args
	}
}

func WithSCArguments(args []string) Option {
	return func(dl *Downloader) {
		dl.scArgs = args
	}
}

func WithWriters(out, err io.Writer) Option {
	return func(dl *Downloader) {
		dl.outWriter = out
		dl.errWriter = err
	}
}
