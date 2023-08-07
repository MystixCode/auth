package root

// "auth/log"

type Service struct {
	// log log.Logger
}

func NewService() *Service {
	return &Service{
		// log: log,
	}
}

func (s *Service) GetRoot() Root {
	return Root{
		Hello: "world",
	}
}
