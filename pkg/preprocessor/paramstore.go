package preprocessor

import (
	"fmt"
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

func (s *ParamStore) GetParam(paramName string) string {
	return s.params[CleanTemplateParam(paramName)]
}

func (s *ParamStore) GetParams(paramNames map[string]string) map[string]string {
	params := map[string]string{}
	for k, v := range paramNames {
		if IsTemplateReplaceParam(v) {
			params[k] = s.GetParam(v)
		} else {
			params[k] = v
		}
	}
	return params
}

func (s *ParamStore) SetParam(stackName string, paramName string, paramValue string) {
	if stackName != "" {
		paramName = fmt.Sprintf("%v.%v", stackName, paramName)
	}
	s.params[paramName] = paramValue
}

func (s *ParamStore) SetParams(stackName string, params map[string]string) {
	for k, v := range params {
		s.SetParam(stackName, k, v)
	}
}
