package config

import "gopkg.in/yaml.v3"

type BigBlueButton struct {
	Secret                 string `yaml:"secret"`
	RecordingsPollInterval string `yaml:"recordingsPollInterval"`
}

var Bbb *BigBlueButton

func Load(value string) {
	var conf *BigBlueButton
	yaml.Unmarshal([]byte(value), &conf)
	Bbb = conf
}
