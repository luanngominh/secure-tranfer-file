package config

type Config struct {
	Port       string
	Address    string
	PrivateKey string
	PublicKey  string
}

var (
	Cfg = &Config{}
)
