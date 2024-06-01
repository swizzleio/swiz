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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "Location to output file to",
				DefaultText: "",
			},
			&cli.BoolFlag{
				Name:  "force",
				Value: false,
				Usage: "Force overwrite of files",
			},
		},
	})
}

func writeYaml[T any](location string, err error, data T) error {
	if err != nil {
		return err
	}
	ser := fileutil.NewYamlHelper[T](appFs)
	err = ser.Set(data).Save(location)

	return err
}

func configGenCmd(ctx *cli.Context) error {
	output := ctx.String("output")
	force := ctx.Bool("force")
	appConfigOut := appconfig.DefaultLocation
	appConfigOutDir := appconfig.DefaultSwizDir
	configOut := fmt.Sprintf("%v/%v", appconfig.DefaultOutLocation, EnvDefFileName)
	configOutDir := appconfig.DefaultOutLocation

	if output != "" {
		if output[len(output)-1] != '/' {
			output += "/"
		}

		appConfigOut = fmt.Sprintf("%v/%v", output, appconfig.DefaultFileName)
		configOut = fmt.Sprintf("%v%v", output, EnvDefFileName)
		appConfigOutDir = output
		configOutDir = output
	}

	fh := fileutil.NewFileHelper(appFs)

	// Check to see if file exists
	if !force {
		fileList := []string{appConfigOut, configOut}
		for _, fileName := range fileList {
			exists, err := fh.FileExists(fileName)
			if err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("file %s already exists, use --force to overwrite", fileName)
			}
		}
	}

	// Parse AWS accounts
	awsAccts, err := getAwsConfig()
	if err != nil {
		return err
	}

	cfg, cfgErr := getCoreConfig(configOut)
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

	cl.Info("Exporting app config to %v\n", appConfigOut)
	cl.Info("Exporting environment definition to %v\n", configOut)
	if output != "" {
		cl.Info("The app config file should be place in the %v directory\n", appconfig.DefaultSwizDir)
	}

	fErr := fh.CreateDirIfNotExist(appConfigOutDir)
	if fErr != nil {
		return fErr
	}
	fErr = fh.CreateDirIfNotExist(configOutDir)
	if fErr != nil {
		return fErr
	}

	serErr := writeYaml[appconfig.AppConfig](appConfigOut, nil, *cfg.AppConfig)
	if serErr != nil {
		return serErr
	}
	serErr = writeYaml[model.EnvironmentConfig](configOut, serErr, envCfg)
	if serErr != nil {
		return serErr
	}

	for _, stack := range stacks {
		stackSerErr := writeYaml[model.StackConfig](model.GenerateFileName(stack.Name), serErr, stack)
		if stackSerErr != nil {
			serErr = stackSerErr
		}
	}

	return serErr
}

func getDefaultEnclave(enclaveNames []string) (string, error) {
	defaultEnclave := ""
	if len(enclaveNames) == 1 {
		return enclaveNames[0], nil
	}

	defaultEnclave, err := cl.AskOptions("Which enclave should be the default", enclaveNames)
	if err != nil {
		return "", err
	}

	return defaultEnclave, nil
}

func getEnclaves(cfg *coreConfig, awsAccts []awswrap.AwsConfig, paramMap map[string]string) ([]model.Enclave, []string, error) {
	enclaves := []model.Enclave{}
	enclaveNames := []string{}
	counter := 0
	for _, acct := range awsAccts {
		pv, err := cl.Ask(fmt.Sprintf("Name the enclave that AWS account %v will be part of. An enclave refers to production, development, test environments (leave blank to ignore)", acct.Profile), false)
		if err != nil {
			return nil, nil, err
		}

		if pv == "" {
			// Generate enclave name, if one isn't provided it's NameMe
			if len(awsAccts) > 1 {
				pv = fmt.Sprintf("%s-%v", model.DefaultEnclaveName, counter)
			} else {
				pv = model.DefaultEnclaveName
			}
		}
		counter++

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
	// Check to make sure we actually have a parameter
	if (len(globalParamList) > 0) && globalParamList[0] != "" {
		for _, gp := range globalParamList {
			pv, err := cl.Ask(fmt.Sprintf("Provide a value for the global parameter %v", gp), false)
			if err != nil {
				return nil, err
			}

			paramMap[gp] = pv
		}
	}
	return paramMap, nil
}

type coreConfig struct {
	EnvName      string
	GlobalParams string
	DomainName   string
	AppConfig    *appconfig.AppConfig
}

func getCoreConfig(envDefLoc string) (*coreConfig, error) {
	qs := []appcli.AskManyOpts{
		{
			Key:           "DomainName",
			Message:       "What domain name do you want to use for this environment",
			TransformMode: appcli.TransformModeTrimSpace,
		},
		{
			Key:           "EnvName",
			Message:       "What to you want to name this environment (leave blank to ignore)",
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
		EnvDefFile: envDefLoc,
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
				Message:       "What region do you want to use for this account",
				TransformMode: appcli.TransformModeTrimSpace,
			},
		}

		resp, askErr := cl.AskMany(qs)
		if askErr != nil {
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
