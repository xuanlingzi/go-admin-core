package config

type FilePath struct {
	Temp string `json:"temp" yaml:"temp"`
}

var FilePathConfig = new(FilePath)
