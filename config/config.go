package config

type EchoServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

var Config EchoServerConfig
