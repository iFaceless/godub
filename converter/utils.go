package converter

import (
	"fmt"
	"runtime"

	"os/exec"
)

func GetEncoderName() string {
	if IsCommandAvailable(FFMPEGEncoder) {
		return FFMPEGEncoder
	} else {
		panic(fmt.Sprintf("command `%s` not found", FFMPEGEncoder))
	}
}

func IsCommandAvailable(name string) bool {
	cmd := exec.Command(getCheckCommand(), name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func getOS() string {
	return runtime.GOOS
}

func getCheckCommand() string {
	if getOS() == "windows" {
		return "where"
	} else {
		return "which"
	}
}
