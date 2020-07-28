package main

import (
	"fmt"
	"path"
	"runtime"

	"os"

	"time"

	"github.com/yeoji/godub"
	"github.com/yeoji/godub/converter"
	"github.com/caicloud/nirvana/log"
)

func main() {
	load()
	export()
	slice()
	reverse()
	properties()
	applyGain()
	appendSegment()
	repeat()
	overlay()
}

func load() {
	filePath := path.Join(dataDirectory(), "code-geass.mp3")
	// Load to get AudioSegment
	segment, _ := godub.NewLoader().Load(filePath)

	file, _ := os.Open(filePath)
	defer file.Close()

	// Load from io.Reader
	segment, _ = godub.NewLoader().Load(file)
	fmt.Println(segment)
}

func export() {
	filePath := path.Join(dataDirectory(), "code-geass.mp3")
	segment, _ := godub.NewLoader().Load(filePath)

	// Export to file
	toFilePath := path.Join(tmpDataDirectory(), "export-code-geass.mp3")
	err := godub.NewExporter(toFilePath).
		WithDstFormat("mp3").
		WithBitRate(converter.MP3BitRatePerfect).
		Export(segment)
	if err != nil {
		log.Fatal(err)
	}

	segment, _ = godub.NewLoader().Load(toFilePath)
	fmt.Println(segment)
}

type Config struct {
	name        string
	formats     []string
	channels    []int
	sampleRates []int
	bitRates    []int
}

func slice() {
	filePath := path.Join(dataDirectory(), "code-geass.mp3")
	segment, _ := godub.NewLoader().Load(filePath)

	slicedSegment, err := segment.Slice(0, 1*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sliced: %s\n", slicedSegment)

	slicedPath := path.Join(tmpDataDirectory(), "godub-slice-code-geass.mp3")
	err = godub.NewExporter(slicedPath).WithDstFormat("mp3").Export(slicedSegment)
	if err != nil {
		log.Fatal(err)
	}

	if r, err := godub.NewLoader().Load(slicedPath); err == nil {
		fmt.Printf("Loaded again: %s", r)
	}
}

func reverse() {
	filePath := path.Join(dataDirectory(), "code-geass.mp3")
	segment, _ := godub.NewLoader().Load(filePath)
	fmt.Printf("Loaded: %s\n", segment)

	newSeg, err := segment.Reverse()
	if err != nil {
		log.Fatal(err)
	}

	newSeg, err = newSeg.Reverse()
	if err != nil {
		log.Fatal(err)
	}

	reversedPath := path.Join(tmpDataDirectory(), "reverse-code-geass.mp3")
	godub.NewExporter(reversedPath).WithDstFormat("mp3").Export(newSeg)

	if r, err := godub.NewLoader().Load(reversedPath); err == nil {
		fmt.Printf("Reversed: %s", r)
	}
}

func properties() {
	filePath := path.Join(dataDirectory(), "code-geass.mp3")
	segment, _ := godub.NewLoader().Load(filePath)
	fmt.Printf("Loaded: %s\n", segment)

	start := time.Now()
	fmt.Printf("RMS: %f\n", segment.RMS())
	fmt.Printf("Calc RMS Elapsed: %s\n", time.Since(start))

	fmt.Printf("Max amplitude: %f\n", segment.MaxPossibleAmplitude())
	fmt.Printf("Max DBFS: %f\n", segment.MaxDBFS())
	fmt.Printf("DBFS: %s\n", segment.DBFS())
	fmt.Printf("Max: %f\n", segment.Max())
}

func applyGain() {
	filePath := path.Join(dataDirectory(), "code-geass.mp3")
	segment, _ := godub.NewLoader().Load(filePath)
	fmt.Printf("Loaded: %s\n", segment)

	newSeg, err := segment.ApplyGain(godub.Volume(10))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(newSeg.DBFS())

	newPth := path.Join(tmpDataDirectory(), "gain-code-geass.mp3")
	godub.NewExporter(newPth).WithDstFormat("mp3").Export(newSeg)
}

func appendSegment() {
	filePath := path.Join(dataDirectory(), "code-geass.mp3")
	segment, _ := godub.NewLoader().Load(filePath)
	fmt.Printf("Loaded: %s\n", segment)

	newSeg, err := segment.Append(segment, segment)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(newSeg)

	newPth := path.Join(tmpDataDirectory(), "append-code-geass.mp3")
	godub.NewExporter(newPth).WithDstFormat("mp3").WithBitRate(converter.MP3BitRatePerfect).Export(newSeg)
}

func repeat() {
	filePath := path.Join(dataDirectory(), "code-geass.mp3")
	segment, _ := godub.NewLoader().Load(filePath)
	fmt.Printf("Loaded: %s\n", segment)

	newSeg, err := segment.Repeat(10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(newSeg)

	newPth := path.Join(tmpDataDirectory(), "repeat-code-geass.mp3")
	godub.NewExporter(newPth).WithDstFormat("mp3").WithBitRate(converter.MP3BitRatePerfect).Export(newSeg)
}

func overlay() {
	segment, _ := godub.NewLoader().Load(path.Join(dataDirectory(), "code-geass.mp3"))
	fmt.Printf("Loaded: %s\n", segment)

	otherSegment, _ := godub.NewLoader().Load(path.Join(dataDirectory(), "ring.mp3"))

	overlaidSeg, err := segment.Overlay(otherSegment, &godub.OverlayConfig{LoopToEnd: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(overlaidSeg)

	destPath := path.Join(tmpDataDirectory(), "overlay-ring.wav")
	godub.NewExporter(destPath).WithDstFormat("wav").Export(overlaidSeg)
}

func tmpDataDirectory() string {
	return path.Join(dataDirectory(), "tmp")
}

func dataDirectory() string {
	_, file, _, _ := runtime.Caller(0)
	examplesDir := path.Dir(path.Dir(file))
	return path.Join(examplesDir, "data")
}
