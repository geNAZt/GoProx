package config

import (
	"os"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Config MainConfig

type IPAndPortConfig struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type MainConfig struct {
	Listener      IPAndPortConfig `yaml:"listener"`
	DefaultServer IPAndPortConfig `yaml:"default"`
	UseTCP        bool            `yaml:"useTcp"`
}

func init() {
	_, err := os.Stat("config.yml")
	if err != nil {
		// File does not exist
		Config = MainConfig{
			Listener: IPAndPortConfig{
				IP:   "0.0.0.0",
				Port: 19132,
			},
			DefaultServer: IPAndPortConfig{
				IP:   "127.0.0.1",
				Port: 19133,
			},
			UseTCP: false,
		}

		// convert object to yaml string
		yamlContent, err := yaml.Marshal(&Config)
		if err != nil {
			panic(err)
		}

		// create new config file
		file, err := os.Create("config.yml" )
		if err != nil {
			panic(err)
		}

		// write yaml content to file
		_, err = file.Write(yamlContent)
		if err != nil {
			panic(err)
		}

		file.Close()
	} else {
		// open file for reading
		file, err := os.Open("config.yml" )
		if err != nil {
			panic(err)
		}

		// Read all data
		data, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		// read yaml content to object
		err = yaml.Unmarshal(data, &Config)
		if err != nil {
			panic(err)
		}

		file.Close()
	}
}
