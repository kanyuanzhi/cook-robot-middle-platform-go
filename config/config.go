package config

type Config struct {
	System             System             `yaml:"system"`
	Postgres           Postgres           `yaml:"postgres"`
	Sqlite             Sqlite             `yaml:"sqlite"`
	CommandRPC         CommandRPC         `yaml:"commandRPC"`
	SoftwareUpdaterRPC SoftwareUpdaterRPC `yaml:"softwareUpdaterRPC"`
}

type System struct {
	Port           int    `yaml:"port"`
	SuccessCode    int    `yaml:"successCode"`
	SuccessMessage string `yaml:"successMessage"`
	ErrorCode      int    `yaml:"errorCode"`
	ErrorMessage   string `yaml:"errorMessage"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Sqlite struct {
	Path     string `yaml:"path"`
	Database string `yaml:"database"`
}

type CommandRPC struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type SoftwareUpdaterRPC struct {
	Host                   string `yaml:"host"`
	Port                   int    `yaml:"port"`
	FileServerPort         int    `yaml:"fileServerPort"`
	SavePath               string `yaml:"savePath"`
	UnzipPath              string `yaml:"unzipPath"`
	UIFolderName           string `yaml:"uiFolderName"`
	MiddlePlatformFilename string `yaml:"middlePlatformFilename"`
	ControllerFilename     string `yaml:"controllerFilename"`
}
