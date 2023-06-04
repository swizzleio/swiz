package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/environment/model"
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

func convertToCamelCase(s string) string {
	s = strings.TrimSpace(s)
	words := strings.Split(s, " ")

	// Convert the first word to lowercase
	words[0] = strings.ToLower(words[0])

	// Convert subsequent words to title case
	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}

	// Join the words together
	return strings.Join(words, "")
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

	fmt.Printf("Exporting files to %v\n", appconfig.DefaultOutLocation)
	fErr := fileutil.CreateDirIfNotExist(appconfig.DefaultOutLocation)
	if fErr != nil {
		return fErr
	}

	fErr = fileutil.YamlToLocation(fmt.Sprintf("%v/%v", appconfig.DefaultOutLocation, appconfig.DefaultFileName), cfg.AppConfig)
	if fErr != nil {
		return fErr
	}
	fErr = fileutil.YamlToLocation(fmt.Sprintf("%v/%v", appconfig.DefaultOutLocation, EnvDefFileName), envCfg)
	if fErr != nil {
		return fErr
	}

	for _, stack := range stacks {
		fErr = fileutil.YamlToLocation(model.GenerateFileName(stack.Name), stack)
		if fErr != nil {
			return fErr
		}
	}

	return nil
}

func getDefaultEnclave(enclaveNames []string) (string, error) {
	defaultEnclave := ""
	if len(enclaveNames) == 1 {
		return enclaveNames[0], nil
	}

	prompt := &survey.Select{
		Message: "Which enclave should be the default?",
		Options: enclaveNames,
	}
	if sErr := survey.AskOne(prompt, &defaultEnclave); sErr != nil {
		return "", sErr
	}
	return defaultEnclave, nil
}

func getEnclaves(cfg *coreConfig, awsAccts []awswrap.AwsConfig, paramMap map[string]string) ([]model.Enclave, []string, error) {
	enclaves := []model.Enclave{}
	enclaveNames := []string{}
	for _, acct := range awsAccts {

		pv := ""
		prompt := &survey.Input{
			Message: fmt.Sprintf("Name the enclave that AWS account %v will be part of (leave blank to ignore)", acct.Name),
		}
		if sErr := survey.AskOne(prompt, &pv); sErr != nil {
			return nil, nil, sErr
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
		templateFile := ""

		prompt := &survey.Input{
			Message: "Enter the filename of the IaC template for your stack (leave blank to exit):",
			Suggest: func(toComplete string) []string {
				files, _ := filepath.Glob(toComplete + "*")
				return files
			},
		}

		if sErr := survey.AskOne(prompt, &templateFile); sErr != nil {
			return nil, sErr
		}

		if strings.TrimSpace(templateFile) == "" {
			break
		}

		stackName := ""

		prompt = &survey.Input{
			Message: "Enter the name of the stack:",
		}
		if sErr := survey.AskOne(prompt, &stackName, survey.WithValidator(survey.Required)); sErr != nil {
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
		pv := ""
		prompt := &survey.Input{
			Message: fmt.Sprintf("Provide a value for the global parameter %v", gp),
		}
		if sErr := survey.AskOne(prompt, &pv); sErr != nil {
			return nil, sErr
		}
		paramMap[gp] = pv
	}
	return paramMap, nil
}

type coreConfig struct {
	EnvName      string
	GlobalParams string
	DomainName   string
	AppConfig    appconfig.AppConfig
}

func getCoreConfig() (*coreConfig, error) {
	qs := []*survey.Question{
		{
			Name:      "envName",
			Prompt:    &survey.Input{Message: "Provide a name of your environment"},
			Transform: survey.TransformString(convertToCamelCase),
			Validate:  survey.Required,
		},
		{
			Name:      "domainName",
			Prompt:    &survey.Input{Message: "What domain name do you want to use for this environment?"},
			Transform: survey.TransformString(strings.TrimSpace),
		},
		{
			Name:      "globalParams",
			Prompt:    &survey.Input{Message: "Specify a comma seperated list for any global parameters (i.e. LogLevel,VpcId)"},
			Transform: survey.TransformString(strings.TrimSpace),
		},
	}

	answers := coreConfig{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return nil, err
	}

	answers.AppConfig = appconfig.Generate(appconfig.EnvDef{
		Name:       answers.EnvName,
		EnvDefFile: fmt.Sprintf("%v/%v", appconfig.DefaultOutLocation, EnvDefFileName),
	})

	return &answers, nil
}

func getAwsConfig() ([]awswrap.AwsConfig, error) {
	fmt.Printf("Scanning for AWS accounts...\n")
	awsAccts, aErr := awswrap.GetAllOrgAccounts()
	if aErr != nil {
		// Get the default account
		awsAcct, aErr := awswrap.GetDefaultConfig()
		if aErr == nil {
			awsAccts = append(awsAccts, *awsAcct)
		}
	}

	// If there are no accounts, prompt for a default account
	if len(awsAccts) == 0 {
		awsCfg := awswrap.AwsConfig{}
		awsQs := []*survey.Question{
			{
				Name:      "name",
				Prompt:    &survey.Input{Message: "Provide a name of your AWS account"},
				Transform: survey.TransformString(strings.TrimSpace),
				Validate:  survey.Required,
			},

			{
				Name:      "accountId",
				Prompt:    &survey.Input{Message: "Enter the AWS account id"},
				Transform: survey.TransformString(strings.TrimSpace),
			},
			{
				Name:      "region",
				Prompt:    &survey.Input{Message: "What region do you want to use for this account?"},
				Transform: survey.TransformString(strings.TrimSpace),
			},
		}

		sErr := survey.Ask(awsQs, &awsCfg)
		if sErr != nil {
			return nil, sErr
		}

		awsAccts = append(awsAccts, awsCfg)
	}
	return awsAccts, nil
}
