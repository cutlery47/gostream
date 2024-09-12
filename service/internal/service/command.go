package service

import (
	"fmt"
	"os/exec"
)

func findmedia(videoName, videoDir string) *exec.Cmd {
	return exec.Command("/bin/bash", "scripts/findmedia.sh", fmt.Sprintf("%v/%v.mp4", videoDir, videoName))
}

func mkdir(videoName, videoDir string) *exec.Cmd {
	return exec.Command("/bin/bash", "scripts/mkdir.sh", fmt.Sprintf("%v/segmented/%v", videoDir, videoName))
}

func segment(videoName, videoDir string) *exec.Cmd {
	return exec.Command("/bin/bash", "scripts/segment.sh", videoName, videoDir)
}
