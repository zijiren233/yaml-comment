package yamlcomment_test

import (
	"os"
	"testing"

	yamlcomment "github.com/zijiren233/yaml-comment"
	"gopkg.in/yaml.v3"
)

type Model struct {
	String        string            `yaml:"string,omitempty,hc=This is a head comment in string,lc=This is a line comment in string,fc=This is a foot comment in string"`
	Map           map[string]string `yaml:"map,omitempty,hc=This is a head comment in map,lc=This is a line comment in map,fc=This is a foot comment in map"`
	Int           int               `yaml:"int,omitempty,hc=This is a head comment in int,lc=This is a line comment in int,fc=This is a foot comment in int"`
	Slice         []string          `yaml:"slice,omitempty,hc=This is a head comment in slice,lc=This is a line comment in slice,fc=This is a foot comment in slice"`
	Float         float64           `yaml:"float,omitempty,hc=This is a head comment in float,lc=This is a line comment in float,fc=This is a foot comment in float"`
	Zero          string            `yaml:"zero,omitempty,hc=This is a head comment in zero,lc=This is a line comment in zero,fc=This is a foot comment in zero"`
	CommaContinue string            `yaml:"comma_continue,hc=This is a head comment, in comma continue"`
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
