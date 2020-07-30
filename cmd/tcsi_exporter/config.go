package main

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type ConfigSpecification struct {
	Http struct {
		Addr string `envconfig:"HTTP_ADDR" default:":9115"`
	}
	Collector struct {
		Token string `envconfig:"COLLECTOR_TOKEN" required:"true"`
	}
}

func ConfigFromEnv() ConfigSpecification {
	config := ConfigSpecification{}
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return config
}
