package main

import (
	"time"

	"fmt"

	"path"
	"runtime"

	"github.com/yeoji/godub"
	"github.com/yeoji/godub/signals"
)

func main() {
	sine()
	pulse()
}

func sine() {
	g := signals.NewSineSignal(0.1)
	segment, _ := g.GenerateAudioSegment(1*time.Second, godub.Volume(0))
	fmt.Println(segment)
	fmt.Println(segment.DBFS())
	godub.NewExporter(path.Join(dataDirectory(), "tmp", "sine-c.wav")).WithDstFormat("wav").Export(segment)
}

func pulse() {
	g := signals.NewSquareSignal(100)
	segment, _ := g.GenerateAudioSegment(1*time.Second, godub.Volume(0))
	fmt.Println(segment)
	fmt.Println(segment.DBFS())
	godub.NewExporter(path.Join(dataDirectory(), "tmp", "pulse-c.wav")).WithDstFormat("wav").Export(segment)
}

func dataDirectory() string {
	_, file, _, _ := runtime.Caller(0)
	examplesDir := path.Dir(path.Dir(file))
	return path.Join(examplesDir, "data")
}
