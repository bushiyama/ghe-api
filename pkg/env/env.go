package env

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Variables ...
type Variables struct {
	Token  string `yaml:"Token"`
	Domain string `yaml:"Domain"`
	User   string `yaml:"User"`
	Repo   string `yaml:"Repo"`
}

// LoadVariables ...
func LoadVariables() (*Variables, error) {

	var buf []byte
	var err error

	buf, err = ioutil.ReadFile("./pkg/env/.env.yaml")

	if buf == nil || err != nil {
		return nil, errors.Wrap(err, "failed to load env variables")
	}

	var vars Variables
	err = yaml.Unmarshal(buf, &vars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load env variables")
	}
	return &vars, nil
}
