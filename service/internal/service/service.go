package service

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

	err = s.handler.findVideo(videoName)
	if err != nil {
		s.eHandler.Handle(err)
		return
	}

	s.handler.findOrCreateDir(videoName)

	err = s.handler.segmentVideo(videoName)
	if err != nil {
		s.eHandler.Handle(err)
		return
	}

	err = s.handler.returnManifest(videoName)
	if err != nil {
		s.eHandler.Handle(err)
		return
	}

}

func New(handler Handler, eHandler ErrHandler) *Service {
	return &Service{
		handler:  handler,
		eHandler: eHandler,
	}
}
