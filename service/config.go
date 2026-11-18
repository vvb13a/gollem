package service

// Config holds app-level settings that are not stored in the database. Add
// deployment knobs here and populate them from env vars in main.go.
type Config struct {
	AppName string
}

func DefaultConfig() *Config {
	return &Config{
		AppName: "gollem",
	}
}
