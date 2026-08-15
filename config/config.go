package config

import "time"

type EchoServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

var Config EchoServerConfig

type CronConfig struct {
	CronFrequency     time.Duration
	LastExecutionTime time.Time
}

var CronCfg CronConfig

func SetUpCronConfig() {
	CronCfg.CronFrequency = time.Duration(time.Second.Milliseconds())
	CronCfg.LastExecutionTime = time.Now()
}
