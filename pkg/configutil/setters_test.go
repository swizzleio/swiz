package configutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetZero(t *testing.T) {
	// Testing different types
	assert.Equal(t, int(0), getZero[int](), "Zero for int is not as expected")
	assert.Equal(t, float64(0), getZero[float64](), "Zero for float64 is not as expected")
	assert.Equal(t, "", getZero[string](), "Zero for string is not as expected")
}

func TestSetOrDefault(t *testing.T) {
	// Testing int
	assert.Equal(t, 5, SetOrDefault[int](0, 5), "Expected to return default value for int zero")
	assert.Equal(t, 10, SetOrDefault[int](10, 5), "Expected to return value itself for int non-zero")

	// Testing float
	assert.Equal(t, 5.5, SetOrDefault[float64](0, 5.5), "Expected to return default value for float zero")
	assert.Equal(t, 10.1, SetOrDefault[float64](10.1, 5.5), "Expected to return value itself for float non-zero")

	// Testing string
	assert.Equal(t, "default", SetOrDefault[string]("", "default"), "Expected to return default value for string zero")
	assert.Equal(t, "value", SetOrDefault[string]("value", "default"), "Expected to return value itself for string non-zero")
}

func TestFlagOrConfig(t *testing.T) {
	// Testing int
	flagInt := 10
	configValInt := 20
	assert.Equal(t, 20, FlagOrConfig[int](flagInt, &configValInt), "Expected to return config value for int")
	assert.Equal(t, 10, FlagOrConfig[int](flagInt, nil), "Expected to return flag value for int when config is nil")

	// Testing float
	flagFloat := 10.5
	configValFloat := 20.5
	assert.Equal(t, 20.5, FlagOrConfig[float64](flagFloat, &configValFloat), "Expected to return config value for float")
	assert.Equal(t, 10.5, FlagOrConfig[float64](flagFloat, nil), "Expected to return flag value for float when config is nil")

	// Testing string
	flagString := "flag"
	configValString := "config"
	assert.Equal(t, "config", FlagOrConfig[string](flagString, &configValString), "Expected to return config value for string")
	assert.Equal(t, "flag", FlagOrConfig[string](flagString, nil), "Expected to return flag value for string when config is nil")
}
