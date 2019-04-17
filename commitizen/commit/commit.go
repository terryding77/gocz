package commit

import (
	"bytes"
	"fmt"
	"os/exec"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"
	"github.com/terryding77/gocz/commitizen/config"
	"github.com/terryding77/gocz/commitizen/validator"
)

// Commitment checks git-commit-message
type Commitment interface {
	Setup(config.Config) error
	Read() (string, error)
	Execute(*Args) error
}

type basicCommitment struct {
	cfg       config.Config
	commitMsg string
}

// New returns a empty Commitment
func New() Commitment {
	return &basicCommitment{}
}

func (cmt *basicCommitment) Setup(c config.Config) error {
	cmt.cfg = c
	return nil
}

func render(tempStr string, replys map[string]interface{}) (string, error) {
	tmpl, err := template.New("commitMsg").Parse(tempStr)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, replys)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func skipThisQuestion(question map[string]interface{}, res map[string]interface{}) bool {
	if skipByNotConfirmKey, ok := question["skip_if_not_confirm"].(string); ok {
		if vbool, ok := res[skipByNotConfirmKey].(bool); ok && !vbool {
			return true
		} else {
			key := question["key"].(string)
			logrus.Debugf("skip question %s by choose 'No' for question %s", key, skipByNotConfirmKey)
		}
	}
	return false
}
func getAskOptions(question map[string]interface{}) (opts []survey.AskOpt) {
	// getValidator
	if validStr, ok := question["valid"].(string); ok {
		validFunc, err := validator.New(validStr)
		if err != nil {
			logrus.Warn("get valid func error")
		}
		if validFunc != nil {
			opts = append(opts, survey.WithValidator(validFunc))
		}
	}
	return opts
}
func interactive(questions interface{}) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	for _, q := range questions.([]interface{}) {
		question := q.(map[string]interface{})
		key := question["key"].(string)
		message := question["message"].(string)
		help := "no help info"
		if ihelp, ok := question["help"].(string); ok {
			help = ihelp
		}
		askOpts := getAskOptions(question)
		switch question["type"].(string) {
		case "select":
			ops := question["option"].([]interface{})
			options := []string{}
			optionMap := map[string]string{}
			for _, op := range ops {
				o := op.(map[string]interface{})
				hint := o["hint"].(string)
				options = append(options, hint)
				optionMap[hint] = o["value"].(string)
			}
			prompt := &survey.Select{
				Message: message,
				Options: options,
				Help:    help,
			}
			var hintValue string
			if err := survey.AskOne(prompt, &hintValue, askOpts...); err != nil {
				return nil, err
			}
			res[key] = optionMap[hintValue]
		case "input":
			// TODO maybe some variable in message
			prompt := &survey.Input{
				Message: message,
				Help:    help,
			}
			value := question["default"].(string)
			if !skipThisQuestion(question, res) {
				if err := survey.AskOne(prompt, &value, askOpts...); err != nil {
					return nil, err
				}
			}
			res[key] = value
		case "confirm":
			// TODO maybe some variable in message
			prompt := &survey.Confirm{
				Message: message,
				Help:    help,
			}
			value := question["default"].(bool)
			if err := survey.AskOne(prompt, &value, askOpts...); err != nil {
				return nil, err
			}
			res[key] = value
		}
	}
	return res, nil
}

func (cmt *basicCommitment) Read() (string, error) {
	replys, err := interactive(cmt.cfg.Get("commit.question"))
	if err != nil {
		return "", err
	}
	cmt.commitMsg, err = render(cmt.cfg.GetString("commit.template"), replys)
	if err != nil {
		return "", err
	}
	return cmt.commitMsg, nil
}

func (cmt *basicCommitment) Execute(args *Args) error {
	rawArgs := args.Combination(cmt.commitMsg)

	logrus.Debug(rawArgs)

	cmd := exec.Command("git", rawArgs...)
	out, err := cmd.CombinedOutput()
	logrus.Debugf("combined out:\n%s\n", string(out))
	if err != nil {
		return fmt.Errorf("cmd.Run() failed with reason: %s", err)
	}

	return nil
}
