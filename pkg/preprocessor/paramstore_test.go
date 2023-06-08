package preprocessor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewParamStore(t *testing.T) {
	store := NewParamStore(nil)
	assert.NotNil(t, store, "Expected NewParamStore to return a non-nil ParamStore")
	assert.NotNil(t, store.params, "Expected NewParamStore to initialize an empty map")

	initialParams := map[string]string{"param1": "value1"}
	store = NewParamStore(initialParams)
	assert.NotNil(t, store, "Expected NewParamStore to return a non-nil ParamStore")
	assert.Equal(t, initialParams, store.params, "Expected NewParamStore to copy the initial params")
}

func TestParamStore_GetParam(t *testing.T) {
	store := NewParamStore(map[string]string{"param1": "value1"})
	value := store.GetParam("param1")
	assert.Equal(t, "value1", value, "Expected GetParam to return the correct value")
}

func TestParamStore_GetParams(t *testing.T) {
	store := NewParamStore(map[string]string{
		"param1": "value1",
		"param2": "value2",
	})

	paramNames := map[string]string{
		"key1": "{{param1}}",
		"key2": "staticValue",
		"key3": "{{param2}}",
	}

	expectedResult := map[string]string{
		"key1": "value1",
		"key2": "staticValue",
		"key3": "value2",
	}

	result := store.GetParams(paramNames)
	assert.Equal(t, expectedResult, result, "Expected GetParams to return the correct value")
}

func TestParamStore_SetParam(t *testing.T) {
	store := NewParamStore(nil)
	store.SetParam("stack1", "param1", "value1")

	assert.Equal(t, "value1", store.params["stack1.param1"], "Expected SetParam to correctly set the parameter value")
}

func TestParamStore_SetParams(t *testing.T) {
	store := NewParamStore(nil)
	params := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}

	store.SetParams("stack1", params)

	expectedResult := map[string]string{
		"stack1.param1": "value1",
		"stack1.param2": "value2",
	}

	assert.Equal(t, expectedResult, store.params, "Expected SetParams to correctly set the parameters")
}
