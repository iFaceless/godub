godub
==================

[godub](https://github.com/iFaceless) lets you manipulate audio in an easy and elegant way. It's deeply inspired by the excellent project [pydub](https://github.com/jiaaro/pydub).

# Why

There are some audio packages in the Go world, but we believe that [pydub](https://github.com/jiaaro/pydub) provides a better way to do stuff to audios. Therefore, here we have [godub](https://github.com/iFaceless)! However, not all the features of [pydub](https://github.com/jiaaro/pydub) are supported.

# Features

- Load audio file, supports mp3/m4a/wav...
- Export/Convert audio with custom config.
- Slice an audio.
- Concatenate audios.
- Repeat an audio.
- Overlay with other audios.
- Reverse an audio.
- ...

# Quickstart

## Load

```go
func main() {
	filePath := path.Join(dataDirectory(), "code-geass.mp3")
	segment, _ := godub.NewLoader().Load(filePath)
	fmt.Println(segment)

	buf, _ := ioutil.ReadAll()

	// Load from buffer, load also accepts `io.Reader` type.
	segment, _ = godub.NewLoader().Load(bytes.NewReader(buf))
	fmt.Println(segment)
}
```

## Export

```go
func main() {
	filePath := path.Join(dataDirectory(), "code-geass.mp3")
	segment, _ := godub.NewLoader().Load(filePath)

	// Export as a mp3 file.
	toFilePath := path.Join(dataDirectory(), "converted", "code-geass.mp3")
	err := godub.NewExporter(toFilePath).
		WithDstFormat("mp3").
		WithBitRate("128k").
		Export(segment)
	if err != nil {
		log.Fatal(err)
	}

	// Let's check it.
	segment, _ = godub.NewLoader().Load(toFilePath)
	fmt.Println(segment)
}
```

## Convert

```go
func main() {
	filePath := path.Join(dataDirectory(), "code-geass.mp3")
	toFilePath := path.Join(dataDirectory(), "converted", "code-geass.m4a")
	w, _ := os.Create(toFilePath)

	err := audio.NewConverter(w).
		WithBitRate("64k").
		WithDstFormat("m4a").
		Convert(filePath)
	if err != nil {
		log.Fatal(err)
	}

	segment, _ := godub.NewLoader().Load(toFilePath)
	fmt.Println(segment)
}
```

## Slice

```go
func slice() {
	filePath := path.Join(dataDirectory(), "code-geass.mp3")
	segment, _ := godub.NewLoader().Load(filePath)

	// Slice from 0 sec to 100 secs.
	slicedSegment, _ := segment.Slice(0, 100*time.Second)

	// Then export it as mp3 file.
	slicedPath := path.Join(dataDirectory(), "slice", "code-geass.mp3")
	godub.NewExporter(slicedPath).WithDstFormat("mp3").Export(slicedSegment)
}
```

## Concatenate

```go
func main() {
	filePath := path.Join(dataDirectory(), "code-geass.mp3")
	segment, _ := godub.NewLoader().Load(filePath)

	// Yep, you can append more than one segment.
	newSeg, err := segment.Append(segment, segment)
	if err != nil {
		log.Fatal(err)
	}

	// Save the newly created audio segment as mp3 file.
	newPth := path.Join(dataDirectory(), "append", "code-geass.mp3")
	godub.NewExporter(newPth).WithDstFormat("mp3").WithBitRate(audio.MP3BitRatePerfect).Export(newSeg)
}
```

## Overlay

```go
func main() {
	segment, _ := godub.NewLoader().Load(path.Join(dataDirectory(), "code-geass.mp3"))
	otherSegment, _ := godub.NewLoader().Load(path.Join(dataDirectory(), "ring.mp3"))

	overlaidSeg, err := segment.Overlay(otherSegment, &godub.OverlayConfig{LoopToEnd: true})
	if err != nil {
		log.Fatal(err)
	}

	godub.NewExporter(destPath).WithDstFormat("wav").Export(path.Join(tmpDataDirectory(), "overlay-ring.wav"))
}
```

# Dependency

[godub](https://github.com/iFaceless/godub)  uses [ffmpeg](https://ffmpeg.org/ffmpeg.html) as its backend to support encoding, decoding and conversion.

# References
1. [go-binary-pack](https://github.com/roman-kachanovsky/go-binary-pack)
1. [Python struct](https://docs.python.org/3/library/struct.html)
1. [Digital Audio - Creating a WAV (RIFF) file](http://www.topherlee.com/software/pcm-tut-wavformat.html)
1. [ffmpeg tutorial](http://keycorner.org/pub/text/doc/ffmpeg-tutorial.htm)
1. [Python: manipulate raw audio data with `audioop`](https://docs.python.org/2/library/audioop.html)
1. [ffmpeg Documentation](https://ffmpeg.org/ffmpeg.html)

# Similar Projects
1. [cryptix/wav](https://github.com/cryptix/wav)
1. [mdlayher/waveform](https://github.com/mdlayher/waveform)
1. [go-audio](https://github.com/go-audio)

# TODO
- [ ] Audio effects.
- [ ] As always, test test test!

# License
[godub](https://github.com/iFaceless/godub) is licensed under the [MIT license](./LICENSE.md). Please feel free and have fun~
