package health

// "auth/log"

type Service struct {
	// logger log.Logger
	state *State
}

func NewService(state *State) *Service {
	return &Service{
		// logger: logger.WithPrefix("service.health"),
		state: state,
	}
}

func (s *Service) GetHealth() Health {
	return Health{
		State:       *s.state,
		StateString: s.state.String(),
		// Report:      s.logger.GetReport(),
	}
}
