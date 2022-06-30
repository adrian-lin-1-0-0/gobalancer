package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config interface {
	Get() ConfigData
}

type ConfigData struct {
	Listener []Listener `json:"listeners"`
	Instance []Instance `json:"instances"`
}

type Listener struct {
	Port              int    `json:"port"`
	SSL               bool   `json:"ssl"`
	SSLCertificate    string `json:"ssl_certificate"`
	SSLCertificateKey string `json:"ssl_certificate_key"`
	Upstream          int    `json:"upstream"`
	HealthCheck       int    `json:"healthcheck"`
	Algo              string `json:"algo"`
}

type Instance struct {
	Addr string `json:"addr"`
}

func (c *ConfigData) Get() ConfigData {
	return *c
}

func Load() (Config, error) {
	jsonFile, err := os.Open("gobalancer.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	serviceConfig := &ConfigData{}

	err = json.Unmarshal(byteValue, serviceConfig)
	if err != nil {
		return nil, err
	}
	return serviceConfig, nil
}
