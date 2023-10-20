package yamlcomment_test

import (
	"os"
	"testing"

	yamlcomment "github.com/zijiren233/yaml-comment"
	"gopkg.in/yaml.v3"
)

type Model struct {
	Map    map[string]*Model2 `yaml:"map,omitempty" lc:"this is comment in line"`
	String string             `yaml:"string,omitempty" hc:"this is comment in head"`
}

type Model2 struct {
	String string            `yaml:"string,omitempty" hc:"this is comment in head"`
	Map    map[string]string `yaml:"map,omitempty" lc:"this is comment in line"`
}

func DefaultModel() *Model {
	return &Model{
		Map: map[string]*Model2{
			"1": DefaultModel2(),
			"2": DefaultModel2(),
		},
	}
}

func DefaultModel2() *Model2 {
	return &Model2{
		String: "default",
	}
}

func TestEncode(t *testing.T) {
	ce := yamlcomment.NewEncoder(yaml.NewEncoder(os.Stdout))
	ce.Encode(DefaultModel())
}
