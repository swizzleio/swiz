package fileutil

import "gopkg.in/yaml.v3"

func YamlFromLocation[T any](location string) (*T, error) {
	return YamlFromLocationWithBaseDir[T]("", location)
}

func YamlFromLocationWithBaseDir[T any](baseDir string, location string) (*T, error) {

	// Open URL
	data, err := OpenUrlWithBaseDir(baseDir, location)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML into StackConfig
	out := new(T)
	err = yaml.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func YamlToLocation[T any](location string, data T) error {
	return YamlToLocationWithBaseDir("", location, data)
}

func YamlToLocationWithBaseDir[T any](baseDir string, location string, data T) error {

	// Marshal YAML into StackConfig
	out, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	// Write URL
	err = WriteUrlWithBaseDir(baseDir, location, out)
	if err != nil {
		return err
	}

	return nil
}
