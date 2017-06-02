package cfg

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	MpdHost  string `json:"mpdhost"`
	MpdPort  uint   `json:"mpdport"`
	HttpPort uint   `json:"httpport"`
	WebRoot  string `json:"webroot"`
	StartMpd bool   `json:"startmpd"`
	LogLevel string `json:"loglevel"`
}

func (c *Config) Load(file string) error {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(f, c)
}
