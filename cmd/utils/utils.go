package utils

import (
	"encoding/json"
	"os"
)

type ConfigJSON struct {
	Network  string `json:"network"`
	Port     int    `json:"port"`
	HostName string `json:"hostname"`
}

// reading .json file
func (c *ConfigJSON) readJSON(path string) (err error) {
	var jsonFile []byte
	jsonFile, err = os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonFile, c)
	return err
}

func (c *ConfigJSON) Init(path string) error {
	return c.readJSON(path)
}
