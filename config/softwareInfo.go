package config

import "time"

type SoftwareInfo struct {
	Name         string    `json:"name" yaml:"name" mapstructure:"name"`
	Version      string    `json:"version" yaml:"version" mapstructure:"version"`
	MachineModel string    `json:"machineModel" yaml:"machineModel" mapstructure:"machineModel"`
	SerialNumber string    `json:"serialNumber" yaml:"serialNumber" mapstructure:"serialNumber"`
	UpdateTime   time.Time `json:"updateTime" yaml:"updateTime" mapstructure:"updateTime"`
}
