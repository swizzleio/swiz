package cmds

import (
	"fmt"
	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/environment/model"
	appcli "github.com/swizzleio/swiz/pkg/cli"
	"github.com/swizzleio/swiz/pkg/drivers/awswrap"
	"github.com/swizzleio/swiz/pkg/fileutil"
	"github.com/urfave/cli/v2"
	"path/filepath"
	"strings"
)

var EnvDefFileName = "env-def.yaml"

func init() {
	addSubCommand("config", &cli.Command{
		Name:   "generate",
		Usage:  "Generate the app config and stack config",
		Action: configGenCmd,
		Flags:  []cli.Flag{},
	})
}

func writeYaml[T any](location string, err error, data T) error {
	if err != nil {
		return err
	}
	ser := fileutil.NewYamlHelper[T]()
	err = ser.Set(data).Save(location)

	return err
}

func configGenCmd(ctx *cli.Context) error {
	// Parse AWS accounts
	awsAccts, err := getAwsConfig()
	if err != nil {
		return err
	}

	cfg, cfgErr := getCoreConfig()
	if cfgErr != nil {
		return cfgErr
	}

	// Parse params
	paramMap, paramErr := getParams(cfg)
	if paramErr != nil {
		return paramErr
	}

	// Parse stacks
	stacks, stackErr := getStacks(paramMap)
	if stackErr != nil {
		return stackErr
	}

	// Parse enclaves
	enclaves, enclaveNames, enclaveErr := getEnclaves(cfg, awsAccts, paramMap)
	if enclaveErr != nil {
		return enclaveErr
	}

	// Determine which is the default enclave
	defaultEnclave, defEncErr := getDefaultEnclave(enclaveNames)
	if defEncErr != nil {
		return defEncErr
	}

	envCfg := model.GenerateEnvironmentConfig(stacks, enclaves, defaultEnclave)

	cl.Info("Exporting files to %v\n", appconfig.DefaultOutLocation)
	fh := fileutil.NewFileHelper()

	fErr := fh.CreateDirIfNotExist(appconfig.DefaultOutLocation)
	if fErr != nil {
		return fErr
	}

	serErr := writeYaml[appconfig.AppConfig](fmt.Sprintf("%v/%v", appconfig.DefaultOutLocation, appconfig.DefaultFileName), nil, *cfg.AppConfig)
	serErr = writeYaml[model.EnvironmentConfig](fmt.Sprintf("%v/%v", appconfig.DefaultOutLocation, EnvDefFileName), serErr, envCfg)
	for _, stack := range stacks {
		serErr = writeYaml[model.StackConfig](model.GenerateFileName(stack.Name), serErr, stack)
	}

	return serErr
}

func getDefaultEnclave(enclaveNames []string) (string, error) {
	defaultEnclave := ""
	if len(enclaveNames) == 1 {
		return enclaveNames[0], nil
	}

	defaultEnclave, err := cl.AskOptions("Which enclave should be the default?", enclaveNames)
	if err != nil {
		return "", err
	}

	return defaultEnclave, nil
}

func getEnclaves(cfg *coreConfig, awsAccts []awswrap.AwsConfig, paramMap map[string]string) ([]model.Enclave, []string, error) {
	enclaves := []model.Enclave{}
	enclaveNames := []string{}
	for _, acct := range awsAccts {
		pv, err := cl.Ask(fmt.Sprintf("Name the enclave that AWS account %v will be part of (leave blank to ignore)", acct.Profile), false)
		if err != nil {
			return nil, nil, err
		}

		if strings.TrimSpace(pv) != "" {
			enclave := model.GenerateEnclave(acct, cfg.DomainName, paramMap)
			enclave.Name = pv

			enclaveNames = append(enclaveNames, pv)
			enclaves = append(enclaves, enclave)
		}
	}
	return enclaves, enclaveNames, nil
}

func getStacks(params map[string]string) ([]model.StackConfig, error) {
	stacks := []model.StackConfig{}
	for {
		templateFile, err := cl.AskAutocomplete("Enter the filename of the IaC template for your stack (leave blank to exit):", false,
			func(toComplete string) []string {
				files, _ := filepath.Glob(toComplete + "*")
				return files
			})
		if err != nil {
			return nil, err
		}

		if strings.TrimSpace(templateFile) == "" {
			break
		}

		stackName, sErr := cl.Ask("Enter the name of the stack:", true)
		if sErr != nil {
			return nil, sErr
		}

		stacks = append(stacks, model.GenerateStackConfig(stackName, templateFile, params))
	}
	return stacks, nil
}

func getParams(cfg *coreConfig) (map[string]string, error) {
	globalParamList := strings.Split(cfg.GlobalParams, ",")
	paramMap := map[string]string{}
	for _, gp := range globalParamList {
		pv, err := cl.Ask(fmt.Sprintf("Provide a value for the global parameter %v", gp), false)
		if err != nil {
			return nil, err
		}

		paramMap[gp] = pv
	}
	return paramMap, nil
}

type coreConfig struct {
	EnvName      string
	GlobalParams string
	DomainName   string
	AppConfig    *appconfig.AppConfig
}

func getCoreConfig() (*coreConfig, error) {
	qs := []appcli.AskManyOpts{
		{
			Key:           "EnvName",
			Message:       "Provide a name of your environment",
			Required:      true,
			TransformMode: appcli.TransformModeCamelCase,
		},
		{
			Key:           "DomainName",
			Message:       "What domain name do you want to use for this environment?",
			TransformMode: appcli.TransformModeTrimSpace,
		},

		{
			Key:           "GlobalParams",
			Message:       "Specify a comma seperated list for any global parameters (i.e. LogLevel,VpcId)",
			TransformMode: appcli.TransformModeTrimSpace,
		},
	}

	resp, err := cl.AskMany(qs)
	if err != nil {
		return nil, err
	}

	answers := coreConfig{
		EnvName:      resp["EnvName"],
		DomainName:   resp["DomainName"],
		GlobalParams: resp["GlobalParams"],
	}

	answers.AppConfig = appConfigMgr.GenFromEnv(appconfig.EnvDef{
		Name:       answers.EnvName,
		EnvDefFile: fmt.Sprintf("%v/%v", appconfig.DefaultOutLocation, EnvDefFileName),
	})

	return &answers, nil
}

func getAwsConfig() ([]awswrap.AwsConfig, error) {
	cl.Info("Scanning for AWS accounts...\n")
	awsCfg, err := awswrap.NewAwsConfigManage()
	if err != nil {
		return nil, err
	}

	awsAccts, aErr := awsCfg.GetAllOrgAccounts()
	if aErr != nil {
		// Get the default account
		awsAcct, cErr := awsCfg.GetDefaultConfig()
		if cErr == nil {
			awsAccts = append(awsAccts, *awsAcct)
		}
	}

	// If there are no accounts, prompt for a default account
	if len(awsAccts) == 0 {
		qs := []appcli.AskManyOpts{
			{
				Key:           "Name",
				Message:       "Provide the name of your AWS account",
				Required:      true,
				TransformMode: appcli.TransformModeTrimSpace,
			},
			{
				Key:           "AccountId",
				Message:       "Enter the AWS account id",
				TransformMode: appcli.TransformModeTrimSpace,
			},

			{
				Key:           "Region",
				Message:       "What region do you want to use for this account?",
				TransformMode: appcli.TransformModeTrimSpace,
			},
		}

		resp, err := cl.AskMany(qs)
		if err != nil {
			return nil, err
		}

		acct := awswrap.AwsConfig{
			Profile:   resp["Name"],
			AccountId: resp["AccountId"],
			Region:    resp["Region"],
		}

		awsAccts = append(awsAccts, acct)
	}
	return awsAccts, nil
}
