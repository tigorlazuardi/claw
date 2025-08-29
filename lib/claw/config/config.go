package config

type Config struct {
	Download  Download  `koanf:"download"`
	Scheduler Scheduler `koanf:"scheduler"`
	Webhooks  Webhooks  `koanf:"webhooks"`
}

func DefaultConfig() Config {
	return Config{
		Download:  DefaultDownload(),
		Scheduler: DefaultScheduler(),
	}
}
