package environment

import (
	"fmt"
	"strings"
)

type ParamStore struct {
	params map[string]string
}

func NewParamStore(params map[string]string) *ParamStore {
	if params == nil {
		params = map[string]string{}
	}

	return &ParamStore{
		params: params,
	}
}

func (s ParamStore) getParam(paramName string) string {
	return s.params[s.cleanParam(paramName)]
}

func (s ParamStore) getParams(paramNames map[string]string) map[string]string {
	params := map[string]string{}
	for k, v := range paramNames {
		if s.isReplaceParam(v) {
		  params[k] = s.getParam(v)
		} else {
			params[k] = v
		}
	}
	return params
}

func (s *ParamStore) setParam(stackName string, paramName string, paramValue string) {
	if stackName != "" {
		paramName = fmt.Sprintf("%v.%v", stackName, paramName)
	}
	s.params[paramName] = paramValue
}

func (s *ParamStore) setParams(stackName string, params map[string]string) {
	for k, v := range params {
		s.setParam(stackName, k, v)
	}
}

func (s ParamStore) cleanParam(paramName string) string {
	if s.isReplaceParam(paramName) {
		// Strip prefix and suffix
		paramName = paramName[2 : len(paramName)-2]
		return paramName
	}

	return paramName
}

func (s ParamStore) isReplaceParam(paramName string) bool {
	return strings.HasPrefix(paramName, "{{") && strings.HasSuffix(paramName, "}}")
}