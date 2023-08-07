package health

// "auth/log"

type Health struct {
	State       State  `json:"-"`
	StateString string `json:"state"`
	// Report      log.Report `json:"report"`
}
