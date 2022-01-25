package mill

import "github.com/AlecAivazis/survey/v2"

func PromptInput(prompt string) (ret string) {
	survey.AskOne(&survey.Input{Message: prompt}, &ret, survey.WithValidator(survey.Required))
	return
}

func PromptPassword(prompt string) (ret string) {
	survey.AskOne(&survey.Password{Message: prompt}, &ret, survey.WithValidator(survey.Required))
	return
}
