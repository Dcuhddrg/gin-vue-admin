package config

type Config struct {
	EncryptionKey string `mapstructure:"encryption-key" json:"encryptionKey" yaml:"encryption-key"`
}

var ConfigData Config
