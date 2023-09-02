package main

import (
	"flag"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/exp/slog"
)

func main() {
	flag.Parse()
	baseDir := flag.Arg(0)
	pattern := filepath.Join(baseDir, "**", "*.wav")
	guard := make(chan struct{}, runtime.NumCPU())
	for _, wav_file := range Must(filepath.Glob(pattern)) {
		guard <- struct{}{}
		go func(wav_file string) {
			ogg_file := strings.TrimSuffix(wav_file, filepath.Ext(wav_file)) + ".ogg"
			err := exec.Command("ffmpeg", "-y", "-i", wav_file, ogg_file).Run()
			if err != nil {
				slog.Warn("Could not convert %s", wav_file)
			}
			<-guard
		}(wav_file)
	}
}

func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}
