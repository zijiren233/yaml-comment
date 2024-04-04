package yamlcomment_test

import (
	"os"
	"testing"

	yamlcomment "github.com/zijiren233/yaml-comment"
	"gopkg.in/yaml.v3"
)

type Model struct {
	Map    map[string]*Model2 `yaml:"map,omitempty" lc:"this is model map comment in line"`
	String string             `yaml:"string,omitempty" hc:"this is model string comment in head"`
}

type Model2 struct {
	String string            `yaml:"string,omitempty" hc:"this is model2 string comment in head"`
	Map    map[string]string `yaml:"map,omitempty" lc:"this is model2 map comment in line"`
	Data   []byte            `yaml:"data,omitempty" hc:"this is model2 data comment in head" lc:"this is model2 data comment in line" fc:"this is model2 data comment in foot"`
	Array  [3]string         `yaml:"array,omitempty" hc:"this is model2 array comment in head" lc:"this is model2 array comment in line" fc:"this is model2 array comment in foot"`
}

func DefaultModel() *Model {
	return &Model{
		Map: map[string]*Model2{
			"1": DefaultModel2(),
			"2": DefaultModel2(),
		},
		String: "some string",
	}
}

func DefaultModel2() *Model2 {
	return &Model2{
		String: "default",
		Data:   []byte("default"),
		Array:  [3]string{"default1", "default2", "default3"},
	}
}

func TestEncode(t *testing.T) {
	ce := yamlcomment.NewEncoder(yaml.NewEncoder(os.Stdout))
	ce.Encode(DefaultModel())
}
