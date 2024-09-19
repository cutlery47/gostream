package utils

import (
	"os/exec"
	"strconv"
)

// creates directory if one doesn't exits
func MKDir(dir string) *exec.Cmd {
	return exec.Command("/bin/bash", "scripts/mkdir.sh", dir)
}

// returns 0 if video exitsts, else 1
func FindVideo(path string) *exec.Cmd {
	return exec.Command("/bin/bash", "scripts/find.sh", path)
}

func SegmentVideoAndCreateManifest(vidPath, manPath, chunkPath string, segTime int) *exec.Cmd {
	return exec.Command("/bin/bash", "scripts/segment.sh", vidPath, manPath, chunkPath, strconv.Itoa(segTime))
}
