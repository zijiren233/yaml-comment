package yamlcomment_test

import (
	"os"
	"testing"

	yamlcomment "github.com/zijiren233/yaml-comment"
	"gopkg.in/yaml.v3"
)

type Model struct {
	String        string            `yaml:"string,omitempty" hc:"this is comment in head"`
	Map           map[string]string `yaml:"map,omitempty" lc:"this is comment in line"`
	Int           int               `yaml:"int,omitempty" fc:"this is comment in foot"`
	Slice         []string          `yaml:"slice,omitempty" hc:"this is comment in head" lc:"this is comment in line"`
	Float         float64           `yaml:"float,omitempty" hc:"this is comment in head" fc:"this is comment in foot"`
	Zero          string            `yaml:"zero,omitempty" lc:"this is comment in line" fc:"this is comment in foot"`
	CommaContinue string            `yaml:"comma_continue" hc:"this is comment in head" lc:"this is comment in line" fc:"this is comment in foot"`
}

func DefaultModel() *Model {
	return &Model{
		String: "default",
		Map: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		Int:   1,
		Slice: []string{"slice1", "slice2"},
		Float: 1.1,
		Zero:  "",
	}
}

func TestEncode(t *testing.T) {
	ce := yamlcomment.NewEncoder(yaml.NewEncoder(os.Stdout))
	ce.Encode(DefaultModel())
}

func TestEncodeMap(t *testing.T) {
	e := yaml.NewEncoder(os.Stdout)
	ce := yamlcomment.NewEncoder(e)
	ce.Encode(DefaultModel().Map)
}

func TestEncodeSlice(t *testing.T) {
	e := yaml.NewEncoder(os.Stdout)
	ce := yamlcomment.NewEncoder(e)
	ce.Encode(DefaultModel().Slice)
}

func TestMarshal(t *testing.T) {
	data, err := yamlcomment.Marshal(DefaultModel())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func TestUnMarshal(t *testing.T) {
	data, err := yamlcomment.Marshal(DefaultModel())
	if err != nil {
		t.Fatal(err)
	}
	model := new(Model)
	err = yaml.Unmarshal(data, model)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", model)
}
