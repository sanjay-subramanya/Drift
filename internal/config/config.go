package config

type Config struct {
	BaseBranch string
	Ignore     []string
}

func Default() Config {
	return Config{
		BaseBranch: "main",
	}
}
