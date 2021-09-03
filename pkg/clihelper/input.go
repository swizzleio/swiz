package clihelper

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
	"log"
)

// GetOrPromptOptions pulls a value from a context. If it doesn't exist, then it prompts the user with a set of options
func GetOrPromptOptions(ctx *cli.Context, key string, promptMessage string, options map[string]string,
	backOption string) string {

	// Check to see if the value was passed in
	val := ctx.String(key)
	if val != "" {
		return val
	}

	// Build the question
	optionList := []string{}
	for k := range options {
		optionList = append(optionList, k)
	}

	// Check to see if there needs to be a back option
	if backOption != "" {
		optionList = append(optionList, backOption)
	}
	question := []*survey.Question{
		{
			Name: "input",
			Prompt: &survey.Select{
				Message:  promptMessage,
				Options:  optionList,
				PageSize: 35,
			},
		},
	}

	answers := struct {
		Input string
	}{}

	// Ask
	err := survey.Ask(question, &answers)
	if err != nil {
		log.Fatalf("launching prompt. %v", err)
	}

	return options[answers.Input]
}
