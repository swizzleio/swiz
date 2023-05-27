package repo

import (
	"fmt"
	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/environment/model"
)

type EnclaveRepo struct {
	enclaves       map[string]*model.Enclave
	defaultEnclave string
}

func NewEnclaveRepo(config appconfig.AppConfig) *EnclaveRepo {
	retVal := &EnclaveRepo{
		enclaves:       map[string]*model.Enclave{},
		defaultEnclave: config.DefaultEnclave,
	}

	for _, enclave := range config.EnclaveDefinition {
		retVal.enclaves[enclave.Name] = &enclave
	}

	return retVal
}

func (r *EnclaveRepo) GetEnclave(name string) (*model.Enclave, error) {
	if name == "" {
		name = r.defaultEnclave
	}

	if enclave, ok := r.enclaves[name]; ok {
		return enclave, nil
	}

	return nil, fmt.Errorf("enclave %s not found", name)
}
