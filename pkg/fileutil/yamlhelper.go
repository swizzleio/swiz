// Package fileutil provides utilities for working with files and YAML data.
package fileutil

import (
	"encoding/base64"
	"github.com/swizzleio/swiz/pkg/security"
	"gopkg.in/yaml.v3"
)

// SerializeHelper is an interface that defines methods for working with YAML data.
// When using this interface, the type T must be a pointer to a struct. Currently, this is type constrained to any
// since go does not yet support struct type constraints in interfaces . There is a [proposal] to add this feature to
// the language.
// Breaking the YAGNI principle, this interface is defined to allow for future extensibility. For example, if we want
// to add a method to this interface that allows for the YAML data to be encrypted, we can do so without breaking
// backwards compatibility.
// Future near term functionality is going to also include the ability to safely parse freeform YAML data.
//
// [proposal]: https://github.com/golang/go/issues/51259
type SerializeHelper[T any] interface {
	Open(location string) (*T, error)
	OpenWithBaseDir(baseDir string, location string) (*T, error)
	Save(location string) error
	SaveWithBaseDir(baseDir string, location string) error
	Set(data T) SerializeHelper[T]
	SetFromB64(data string) error
	Get() T
	GetBase64() (*Base64Resp, error)
	ParseBase64(data string) error
	GetSignature() (sig string, wordList string, err error)
}

// Base64Resp is a struct that contains the base64 encoded data, the word list used to generate the signature,
type Base64Resp struct {
	Encoded   string
	WordList  string
	Signature string
}

// YamlHelp is a struct that implements the YamlHelper interface.
type YamlHelp[T any] struct {
	Yaml T             // The YAML data
	f    FileUrlHelper // Helper for file URLs
	fh   FileHelper    // Helper for files
}

// NewYamlHelper creates a new instance of YamlHelper with default settings.
func NewYamlHelper[T any]() SerializeHelper[T] {
	return &YamlHelp[T]{
		f: NewFileUrlHelper(),
	}
}

// Open reads and parses the YAML data from the specified location.
func (y YamlHelp[T]) Open(location string) (*T, error) {
	return y.OpenWithBaseDir("", location)
}

// OpenWithBaseDir reads and parses the YAML data from the specified location, using a base directory.
func (y YamlHelp[T]) OpenWithBaseDir(baseDir string, location string) (*T, error) {
	// Open URL
	data, err := y.f.OpenUrlWithBaseDir(baseDir, location)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML into T
	out := new(T)
	err = yaml.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// Save writes the YAML data to the specified location.
func (y YamlHelp[T]) Save(location string) error {
	return y.SaveWithBaseDir("", location)
}

// SaveWithBaseDir writes the YAML data to the specified location, using a base directory.
func (y YamlHelp[T]) SaveWithBaseDir(baseDir string, location string) error {
	// Marshal T into YAML
	out, err := yaml.Marshal(y.Yaml)
	if err != nil {
		return err
	}

	// Write URL
	err = y.f.WriteUrlWithBaseDir(baseDir, location, out)
	if err != nil {
		return err
	}

	return nil
}

// Set sets the YAML data of the YamlHelper.
func (y YamlHelp[T]) Set(data T) SerializeHelper[T] {
	y.Yaml = data

	return &y
}

// SetFromB64 parses the YAML data from the specified base64 encoded string.
func (y YamlHelp[T]) SetFromB64(data string) error {
	// Decode base64
	b64 := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(b64, []byte(data))
	if err != nil {
		return err
	}

	// Unmarshal YAML into T
	return yaml.Unmarshal(b64[:n], &y.Yaml)
}

// Get returns the YAML data stored in the YamlHelper.
func (y YamlHelp[T]) Get() T {
	return y.Yaml
}

// GetBase64 returns the YAML data stored in the YamlHelper as a base64 encoded string.
func (y YamlHelp[T]) GetBase64() (*Base64Resp, error) {
	// Marshal YAML
	out, err := yaml.Marshal(y.Yaml)
	if err != nil {
		return nil, err
	}

	// Get encoding and signature
	retVal := &Base64Resp{}
	retVal.Encoded = base64.StdEncoding.EncodeToString(out)
	retVal.Signature, retVal.WordList = security.GetSha256AndWordList(retVal.Encoded)

	return retVal, nil
}

// ParseBase64 parses the YAML data from the specified base64 encoded string.
func (y YamlHelp[T]) ParseBase64(data string) error {
	// Decode base64
	b64 := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(b64, []byte(data))
	if err != nil {
		return err
	}

	// Unmarshal YAML into T
	out := new(T)
	err = yaml.Unmarshal(b64[:n], &out)
	if err != nil {
		return err
	}

	return nil
}

// GetSignature returns the signature and word list of the YAML data stored in the YamlHelper.
func (y YamlHelp[T]) GetSignature() (sig string, wordList string, err error) {
	// Marshal YAML
	var out []byte
	out, err = yaml.Marshal(y.Yaml)
	if err != nil {
		return
	}

	// Get encoding and signature
	sig, wordList = security.GetSha256AndWordList(string(out))

	return
}
