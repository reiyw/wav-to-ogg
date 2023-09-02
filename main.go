package main

import (
	"flag"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/exp/slog"
)

func main() {
	flag.Parse()
	baseDir := flag.Arg(0)
	guard := make(chan struct{}, runtime.NumCPU())
	filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) != ".wav" {
			return nil
		}
		guard <- struct{}{}
		go func(wav_file string) {
			ogg_file := strings.TrimSuffix(wav_file, filepath.Ext(wav_file)) + ".ogg"
			err := exec.Command("ffmpeg", "-y", "-i", wav_file, ogg_file).Run()
			if err != nil {
				slog.Warn("Could not convert %s", wav_file)
			} else {
				os.Remove(wav_file)
			}
			<-guard
		}(path)
		return nil
	})
}

func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}
