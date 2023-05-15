package repo

import (
	"fmt"

	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/environment/model"
	"github.com/swizzleio/swiz/pkg/errtype"
	"github.com/swizzleio/swiz/pkg/fileutil"
)

type EnvironmentRepo struct {
	envCfg      map[string]*model.EnvironmentConfig
	baseDir     string
	defaultName string
}

func NewEnvironmentRepo(config appconfig.AppConfig) (*EnvironmentRepo, error) {
	retVal := &EnvironmentRepo{
		envCfg:  map[string]*model.EnvironmentConfig{},
		baseDir: config.BaseDir,
	}

	errList := errtype.ErrList{}

	// Bootstrap environment from YAML
	for _, envDef := range config.EnvDefinition {
		yamlData, err := fileutil.YamlFromLocationWithBaseDir[model.EnvironmentConfig](config.BaseDir, envDef.EnvDefFile)
		if err != nil {
			// There is an error in the environment definition, this should not be fatal
			errList.Add(err)
		} else {
			yamlData.EnvDefName = envDef.Name
			retVal.envCfg[envDef.Name] = yamlData
			if envDef.Default {
				retVal.defaultName = envDef.Name
			}
		}
	}

	return retVal, errList.ErrOrNil()
}

func (r *EnvironmentRepo) GetEnvironmentByDef(envDef string) (*model.EnvironmentConfig, error) {
	if envDef == "" {
		envDef = r.defaultName
	}

	// Check to see if there is an environment with the name in the config list
	envCfg, ok := r.envCfg[envDef]
	if !ok {
		return nil, fmt.Errorf("environment %s not found", envDef)
	}

	// Populate stack definition
	if envCfg.Stacks == nil {
		envCfg.Stacks = map[string]*model.StackConfig{}

		// Load stack files
		for _, stackCfg := range envCfg.StackCfgDef {
			stack, err := fileutil.YamlFromLocationWithBaseDir[model.StackConfig](r.baseDir, stackCfg.ConfigFile)
			if err != nil {
				// TODO: Check if error is due to handlebars incomptability
				// Unlike environment, a stack error is fatal
				return nil, err
			}

			stack.Name = stackCfg.Name
			envCfg.Stacks[stackCfg.Name] = stack
		}
	}

	return envCfg, nil
}
