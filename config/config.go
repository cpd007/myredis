package config

import (
	"time"

	"github.com/cpd007/myredis/constants"
)

type RedisServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

var Config RedisServerConfig

type CronConfig struct {
	CronFrequency     time.Duration
	LastExecutionTime time.Time
}

var CronCfg CronConfig

func SetupCronConfig() {
	CronCfg.CronFrequency = time.Duration(time.Second.Milliseconds())
	CronCfg.LastExecutionTime = time.Now()
}

type EvictionConfig struct {
	EvictionStrategy constants.EvictionStrategy
	MaxKeyLimit      int
}

var EvictionCfg EvictionConfig

func SetupEvictionConfig() {
	EvictionCfg = EvictionConfig{
		EvictionStrategy: constants.EvictFist,
		MaxKeyLimit:      5,
	}
}
