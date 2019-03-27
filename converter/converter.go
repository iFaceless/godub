package converter

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"path"

	"github.com/iFaceless/godub/utils"
	"github.com/tink-ab/tempfile"
)

var (
	DefaultCodecs = map[string]string{
		"ogg": "libvorbis",
	}
	ValidCoverExtensions = utils.NewSet(".png", ".jpg", ".jpeg", ".bmp", ".tif", ".tiff")
	ValidID3TagVersions  = utils.NewSet(3, 4)
	FileExtAlias         = map[string]string{
		"wave": "wav",
	}
)

const (
	FFMPEGEncoder = "ffmpeg"
)

const (
	MP3BitRateEconomy  = "64k"
	MP3BitRateStandard = "128k"
	MP3BitRateGood     = "192k"
	MP3BitRatePerfect  = "320k"

	M4ABitRateEconomy  = "64k"
	M4ABitRateStandard = "128k"
	M4ABitRateGood     = "160k"
	M4ABitRatePerfect  = "256k"
)

type Converter struct {
	w             io.Writer
	channels      int
	dstFormat     string
	bitRate       string
	codec         string
	coverPath     string
	tags          map[string]string
	id3TagVersion int
	sampleRate    int
	params        []string
	cmd           *exec.Cmd

	// It's a temp file
	srcFilename string
}

func NewConverter(w io.Writer) *Converter {
	return &Converter{
		w:             w,
		dstFormat:     "mp3",
		params:        make([]string, 0),
		id3TagVersion: 4,
		cmd: exec.Command(
			GetEncoderName(),
			// Always overwrite existing files
			"-y",
		),
	}
}

func (c *Converter) WithWriter(w io.Writer) *Converter {
	c.w = w
	return c
}

func (c *Converter) WithDstFormat(f string) *Converter {
	if f == "" {
		return c
	}

	c.dstFormat = c.fixFormat(f)
	return c
}

func (c *Converter) fixFormat(f string) string {
	f = strings.TrimSpace(strings.ToLower(f))
	if v, ok := FileExtAlias[f]; ok {
		return v
	} else {
		return f
	}
}

func (c *Converter) WithChannels(v int) *Converter {
	if v == 1 || v == 2 {
		c.channels = v
	}
	return c
}

func (c *Converter) WithBitRate(rate string) *Converter {
	if rate == "" {
		return c
	}

	c.bitRate = rate
	return c
}

func (c *Converter) WithSampleRate(rate int) *Converter {
	if rate == 0 {
		return c
	}
	c.sampleRate = rate
	return c
}

func (c *Converter) WithParams(p ...string) *Converter {
	if len(p) == 0 {
		return c
	}

	c.params = p
	return c
}

func (c *Converter) WithCodec(codec string) *Converter {
	if codec == "" {
		return c
	}

	c.codec = codec
	return c
}

func (c *Converter) WithCover(coverPath string) *Converter {
	if coverPath == "" {
		return c
	}

	c.coverPath = coverPath
	return c
}

func (c *Converter) WithTags(tags map[string]string) *Converter {
	if len(tags) == 0 {
		return c
	}

	c.tags = tags
	return c
}

func (c *Converter) WithID3TagVersion(v int) *Converter {
	if v == 0 {
		return c
	}

	c.id3TagVersion = v
	return c
}

func (c *Converter) DstFormat() string {
	return c.dstFormat
}

func (c *Converter) Convert(src interface{}) error {
	switch src := src.(type) {
	case io.Reader:
		file, err := ioutil.TempFile("", "src-file")
		if err != nil {
			return err
		}
		defer os.Remove(file.Name())

		_, err = io.Copy(file, src)
		if err != nil {
			return err
		}
		c.srcFilename = file.Name()
	case string: // Treat it as source filename
		c.srcFilename = src
	default:
		return fmt.Errorf("conversion error, expected `io.Reader` or file path to original audio")
	}

	return c.doConvert()
}

func (c *Converter) doConvert() error {
	dstFile, err := tempfile.TempFile("", "dst", "."+c.dstFormat)
	if err != nil {
		return err
	}
	defer os.Remove(dstFile.Name())

	c.extendCmdArgs("-i", c.srcFilename)
	c.extendCodecFormatArgs()
	c.extendChannelArgs()

	err = c.extendCoverArgs()
	if err != nil {
		return err
	}

	c.extendBitRateArgs()
	c.extendSampleRateArgs()

	err = c.extendTagsArgs()
	if err != nil {
		return err
	}

	c.extendExtraArgs()
	c.extendCmdArgs(dstFile.Name())

	err = c.cmd.Run()
	if err != nil {
		return EncodeError(fmt.Sprintf("encoding failed: %s", err))
	}

	// Copy to dst writer
	buf, err := ioutil.ReadFile(dstFile.Name())
	c.w.Write(buf)

	return err
}

func (c *Converter) extendCodecFormatArgs() {
	// Set codec
	if c.codec == "" {
		c.codec = DefaultCodecs[c.dstFormat]
	}
	if c.codec != "" {
		c.extendCmdArgs("-acodec", c.codec)
	}
}

func (c *Converter) extendChannelArgs() {
	if c.channels == 1 || c.channels == 2 {
		c.extendCmdArgs("-ac", fmt.Sprintf("%d", c.channels))
	}
}

func (c *Converter) extendCoverArgs() error {
	if c.coverPath != "" {
		coverExt := path.Ext(c.coverPath)
		if !ValidCoverExtensions.Has(coverExt) {
			return InvalidCoverError(fmt.Sprintf("Cover image format '%s' not supported", coverExt))
		}

		if c.dstFormat != "mp3" {
			return InvalidCoverError("Cover images are only allowed for MP3 files")
		}

		c.extendCmdArgs("-i", c.coverPath, "-map", "0", "-map", "1", "-c:v", "mjpeg")
	}
	return nil
}

func (c *Converter) extendBitRateArgs() {
	if c.bitRate != "" {
		c.extendCmdArgs("-b:a", c.bitRate)
	}
}

func (c *Converter) extendSampleRateArgs() {
	if c.sampleRate > 0 {
		c.extendCmdArgs("-ar", fmt.Sprintf("%d", c.sampleRate))
	}
}

func (c *Converter) extendTagsArgs() error {
	if len(c.tags) > 0 {
		for key, value := range c.tags {
			c.extendCmdArgs("-metadata", fmt.Sprintf("%s=%s", key, value))
		}

		if c.dstFormat == "mp3" {
			if ValidID3TagVersions.Has(c.id3TagVersion) {
				c.extendCmdArgs("-id3v2_version", fmt.Sprintf("%d", c.id3TagVersion))
			} else {
				return InvalidID3TagVersionError(
					fmt.Sprintf("id3v2_version '%d' is not allowed", c.id3TagVersion))
			}
		}
	}

	return nil
}

func (c *Converter) extendExtraArgs() {
	if runtime.GOOS == "darwin" && c.codec == "mp3" {
		c.extendCmdArgs("-write_xing", "0")
	}

	c.extendCmdArgs(c.params...)
}

func (c *Converter) extendCmdArgs(args ...string) {
	c.cmd.Args = append(c.cmd.Args, args...)
}
