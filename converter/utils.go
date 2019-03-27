package converter

import (
	"fmt"

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
	cmd := exec.Command("command", "-v", name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
