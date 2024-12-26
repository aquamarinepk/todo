package am

func DefOpts(log Logger, cfg *Config) []Option {
	return []Option{
		WithLog(log),
		WithCfg(cfg),
	}
}
