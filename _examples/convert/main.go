package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"strings"

	"github.com/iFaceless/godub/converter"
)

var (
	DefaultConfig = Config{
		"default",
		[]string{"mp3"},
		[]int{2},
		[]int{44100},
		[]int{192, 128, 96},
	}

	Configs = []Config{
		DefaultConfig,
	}
)

type Config struct {
	id          string
	formats     []string
	channels    []int
	sampleRates []int
	bitRates    []int
}

type Task struct {
	name        string
	srcFilePath string
	destDir     string
	numChannel  int
	dstFormat   string
	sampleRate  int
	bitRate     int
}

func (t *Task) String() string {
	return fmt.Sprintf("Task<%s: channel-%d, sample rate-%d, bitrate-%d, src-'%s'>",
		t.name, t.numChannel, t.sampleRate, t.bitRate, t.srcFilePath)
}

func main() {
	for _, task := range collectTasks()[1:] {
		doWork(task)
	}
}

func collectTasks() (tasks []*Task) {
	for _, c := range Configs {
		for _, numChannel := range c.channels {
			for _, format := range c.formats {
				for _, sr := range c.sampleRates {
					for _, br := range c.bitRates {
						task := &Task{
							name:        c.id,
							srcFilePath: path.Join(dataDirectory(), "code-geass.mp3"),
							dstFormat:   format,
							destDir:     path.Join(tmpDataDirectory(), c.id),
							numChannel:  numChannel,
							sampleRate:  sr,
							bitRate:     br,
						}
						tasks = append(tasks, task)
					}
				}
			}
		}
	}
	return
}

func doWork(task *Task) {
	_, fileName := path.Split(task.srcFilePath)
	fileName = strings.Replace(fileName, path.Ext(fileName), "", len(fileName))
	destFile := path.Join(task.destDir,
		fmt.Sprintf("%s-%s-c_%d-sr_%d-br_%dk.%s",
			fileName,
			task.name,
			task.numChannel,
			task.sampleRate,
			task.bitRate,
			task.dstFormat))

	w, _ := os.Create(destFile)

	fmt.Println("Running:", task)

	// 配置下导出参数
	err := converter.NewConverter(w).
		WithBitRate(fmt.Sprintf("%dk", task.bitRate)).
		WithDstFormat(task.dstFormat).
		WithChannels(task.numChannel).
		WithSampleRate(task.sampleRate).
		Convert(task.srcFilePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Task completed, saved to '%s'\n", destFile)
}

func tmpDataDirectory() string {
	return path.Join(dataDirectory(), "tmp")
}

func dataDirectory() string {
	_, file, _, _ := runtime.Caller(0)
	examplesDir := path.Dir(path.Dir(file))
	return path.Join(examplesDir, "data")
}
