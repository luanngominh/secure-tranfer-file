package config

import "crypto/rsa"

type Config struct {
	Port        string
	Address     string
	PrivateKey  string
	PublicKey   string
	StoragePath string
}

var (
	Cfg        = &Config{}
	PrivateKey *rsa.PrivateKey
)
