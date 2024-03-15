package utils

import (
	"os"

	yaml "github.com/goccy/go-yaml"
)

func LoadYaml(structData any, fileName string) error {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return yaml.UnmarshalWithOptions(bytes, structData, yaml.UseJSONUnmarshaler())
}
