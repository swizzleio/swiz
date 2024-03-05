package fileutil

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/yaml.v3"

	"github.com/swizzleio/swiz/mocks/pkg/fileutil"
)

type DummyStruct struct {
	Version      int               `yaml:"version"`
	Name         string            `yaml:"-"`
	RawName      string            `yaml:"-"`
	Order        int               `yaml:"-"`
	Parameters   map[string]string `yaml:"params"`
	TemplateFile string            `yaml:"template_file"`
}

func TestOpenWithBaseDir(t *testing.T) {
	tests := []struct {
		name        string
		baseDir     string
		location    string
		inputData   map[string]interface{}
		expectedErr bool
		mockOpenErr error
	}{
		{
			name:        "successful open",
			baseDir:     "",
			location:    "testdata/sample.yaml",
			inputData:   map[string]interface{}{"key": "value"},
			expectedErr: false,
		},
		{
			name:        "failed open due to wrong location",
			baseDir:     "",
			location:    "testdata/nonexistent.yaml",
			expectedErr: true,
			mockOpenErr: errors.New("dummy error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileUrlHelper := mockfileutil.FileUrlHelper{}
			helper := YamlHelp[DummyStruct]{
				f: &mockFileUrlHelper,
			}

			// Mock the method
			yamlData, _ := yaml.Marshal(tt.inputData)
			mockFileUrlHelper.On("OpenUrlWithBaseDir", tt.baseDir, tt.location).Return(yamlData, tt.mockOpenErr)

			// Test OpenWithBaseDir
			_, err := helper.OpenWithBaseDir(tt.baseDir, tt.location)
			assert.Equal(t, tt.expectedErr, err != nil)
		})
	}
}

func TestSaveWithBaseDir(t *testing.T) {
	tests := []struct {
		name        string
		baseDir     string
		location    string
		inputData   DummyStruct
		expectedErr bool
		mockSaveErr error
	}{
		{
			name:     "successful save",
			baseDir:  "",
			location: "testdata/sample.yaml",
			inputData: DummyStruct{
				Version: 1,
				Name:    "foo bar 123",
				RawName: "foobar",
				Order:   3,
				Parameters: map[string]string{
					"key": "value",
				},
				TemplateFile: "stuff.yaml",
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileUrlHelper := mockfileutil.FileUrlHelper{}
			helper := YamlHelp[DummyStruct]{
				f: &mockFileUrlHelper,
			}

			// Mock the method
			mockFileUrlHelper.On("WriteUrlWithBaseDir", tt.baseDir, tt.location, mock.Anything).Return(tt.mockSaveErr)

			// Test SaveWithBaseDir
			helper.Set(tt.inputData)
			err := helper.SaveWithBaseDir(tt.baseDir, tt.location)
			assert.Equal(t, tt.expectedErr, err != nil)
		})
	}
}

func TestSetFromB64(t *testing.T) {
	tests := []struct {
		name        string
		inputData   string
		expectedErr bool
	}{
		{
			name:        "valid base64",
			inputData:   base64.StdEncoding.EncodeToString([]byte("key: value")),
			expectedErr: false,
		},
		{
			name:        "invalid base64",
			inputData:   "invalid!!",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := YamlHelp[DummyStruct]{}

			// Test SetFromB64
			err := helper.SetFromB64(tt.inputData)
			assert.Equal(t, tt.expectedErr, err != nil)
		})
	}
}

func TestGetBase64(t *testing.T) {
	tests := []struct {
		name        string
		inputData   *DummyStruct
		output      string
		expectedErr bool
	}{
		{
			name: "successful base64 encoding",
			inputData: &DummyStruct{
				Version: 1,
				Name:    "foo bar 123",
				RawName: "foobar",
				Order:   3,
				Parameters: map[string]string{
					"key": "value",
				},
				TemplateFile: "stuff.yaml",
			},
			expectedErr: false,
			output:      "dmVyc2lvbjogMQpwYXJhbXM6CiAgICBrZXk6IHZhbHVlCnRlbXBsYXRlX2ZpbGU6IHN0dWZmLnlhbWwK",
		},
		{
			name:        "nil input data",
			inputData:   nil,
			expectedErr: true,
			output:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := YamlHelp[*DummyStruct]{}
			helper.Set(tt.inputData)

			// Test GetBase64
			out, err := helper.GetBase64()
			if !tt.expectedErr {
				assert.Equal(t, tt.output, out.Encoded)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestParseBase64(t *testing.T) {
	tests := []struct {
		name        string
		inputData   string
		expectedErr bool
	}{
		{
			name:        "valid base64",
			inputData:   base64.StdEncoding.EncodeToString([]byte("key: value")),
			expectedErr: false,
		},
		{
			name:        "invalid base64",
			inputData:   "invalid!!",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := YamlHelp[DummyStruct]{}

			// Test ParseBase64
			err := helper.ParseBase64(tt.inputData)
			assert.Equal(t, tt.expectedErr, err != nil)
		})
	}
}

func TestGetSignature(t *testing.T) {
	tests := []struct {
		name        string
		inputData   *DummyStruct
		expectedErr bool
	}{
		{
			name: "successful signature generation",
			inputData: &DummyStruct{
				Version: 1,
				Name:    "foo bar 123",
				RawName: "foobar",
				Order:   3,
				Parameters: map[string]string{
					"key": "value",
				},
				TemplateFile: "stuff.yaml",
			},
			expectedErr: false,
		},
		{
			name:        "nil input data",
			inputData:   nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := YamlHelp[*DummyStruct]{}
			helper.Set(tt.inputData)

			// Test GetSignature
			_, _, err := helper.GetSignature()
			assert.Equal(t, tt.expectedErr, err != nil)
		})
	}
}
