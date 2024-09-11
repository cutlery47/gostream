package service

import (
	"fmt"
	"log"
	"os/exec"
)

const videoDir = "vids"

type Service struct {
	handler  Handler
	eHandler ErrHandler
}

func (s *Service) Run() {
	videoName, err := s.handler.getVideoName()
	if err != nil {
		s.eHandler.Handle(err)
		return
	}

	// check if requested mediafile exists in the first place
	_, err = exec.Command("/bin/bash", "scripts/findmedia.sh", fmt.Sprintf("%v/%v.mp4", videoDir, videoName)).Output()
	if err != nil {
		log.Println("Find: ", err)
		return
	}

	// create a directory for segmented file if it doesnt exist yet
	_, err = exec.Command("/bin/bash", "scripts/mkdir.sh", fmt.Sprintf("%v/segmented/%v", videoDir, videoName)).Output()
	// if file hasn't already been segmented -- we segment it
	if err != nil {
		// segment the media file
		cmd := exec.Command(
			"ffmpeg", "-i", fmt.Sprintf("%v/%v.mp4", videoDir, videoName),
			"-bsf:v", "h264_mp4toannexb",
			"-codec", "copy",
			"-hls_list_size", "0",
			fmt.Sprintf("%v/segmented/%v/%v.m3u8", videoDir, videoName, videoName))

		if _, err := cmd.Output(); err != nil {
			log.Println("Segmentation error:", err)
		}
	}

}

func New(handler Handler, eHandler ErrHandler) *Service {
	return &Service{
		handler:  handler,
		eHandler: eHandler,
	}
}
