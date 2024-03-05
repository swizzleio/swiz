package repo

import (
	"fmt"

	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/environment/model"
	"github.com/swizzleio/swiz/pkg/errtype"
	"github.com/swizzleio/swiz/pkg/fileutil"
)

type EnvironmentRepo struct {
	envCfg   map[string]*model.EnvironmentConfig
	config   appconfig.AppConfig
	serEnv   fileutil.SerializeHelper[model.EnvironmentConfig] // TODO: May want to rethink the template
	serStack fileutil.SerializeHelper[model.StackConfig]
	openUrl  fileutil.FileUrlHelper
}

func NewEnvironmentRepo(config appconfig.AppConfig) *EnvironmentRepo {
	return &EnvironmentRepo{
		envCfg:   map[string]*model.EnvironmentConfig{},
		config:   config,
		serEnv:   fileutil.NewYamlHelper[model.EnvironmentConfig](),
		serStack: fileutil.NewYamlHelper[model.StackConfig](),
		openUrl:  fileutil.NewFileUrlHelper(),
	}
}

func (r *EnvironmentRepo) Bootstrap() error {
	errList := errtype.ErrList{}

	// Bootstrap environment from YAML
	for _, envDef := range r.config.EnvDefinition {
		yamlData, err := r.serEnv.OpenWithBaseDir(r.config.BaseDir, envDef.EnvDefFile)
		if err != nil {
			// There is an error in the environment definition, this should not be fatal
			errList.Add(err)
		} else {

		}

		yamlData.EnvDefName = envDef.Name
		r.envCfg[envDef.Name] = yamlData
	}

	return errList.ErrOrNil()

}

func (r *EnvironmentRepo) GetEnvironmentByDef(envDef string) (*model.EnvironmentConfig, error) {
	if envDef == "" {
		envDef = r.config.DefaultEnv
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
			stack, err := r.serStack.OpenWithBaseDir(r.config.BaseDir, stackCfg.ConfigFile)
			if err != nil {
				// TODO: Check if error is due to handlebars incompatibility
				// Unlike environment, a stack error is fatal
				return nil, err
			}

			templateFile, err := r.openUrl.UrlWithBaseDir(r.config.BaseDir, stack.TemplateFile)
			if err != nil {
				return nil, err
			}

			stack.TemplateFile = templateFile
			stack.Name = stackCfg.Name
			stack.RawName = stackCfg.Name
			stack.Order = stackCfg.Order
			envCfg.Stacks[stackCfg.Name] = stack
		}
	}

	return envCfg, nil
}
