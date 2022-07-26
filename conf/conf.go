package conf

import (
	"encoding/json"
	"io/ioutil"

	"github.com/guowenshuai/ieth/types"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func LoadYaml(file string, config *types.Config) error {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		logrus.Printf("yamlFile.Get err   #%v ", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		logrus.Errorf("read config err: %s", err.Error())
		return err
	}
	d, _ := json.Marshal(config)
	logrus.Printf("%s", string(d))
	return nil
}
