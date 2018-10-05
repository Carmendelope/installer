/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package workflow

import (
	"bytes"
	"encoding/json"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"regexp"
	"strings"
	"text/template"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/commands"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
)

type rawWorkflow struct {
	Description string            `json:"description"`
	Commands    []json.RawMessage `json:"commands"`
}

// Parser structure with the required parameters.
type Parser struct {
	cmdParser commands.CmdParser
}

// NewParser creates a new parser.
func NewParser() *Parser {
	return &Parser{*commands.NewCmdParser()}
}

// ReadWorkflow reads a workflow from a file, parsing the data and applying the template.
//   params:
//     filePath The path of the file with the workflow.
//     name The name of the workflow.
//     params The template parameters.
//   returns:
//     A Workflow structure.
//     An error if the workflow cannot be generated.
func (p *Parser) ReadWorkflow(filePath string, name string, params Parameters) (*Workflow, derrors.Error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, derrors.NewUnavailableError(errors.CannotReadWorkflowFile, err)
	}
	return p.ParseWorkflow(string(content), name, params)
}

// ParseWorkflowWithRequest processes an incoming parser workflow request.
func (p *Parser) ParseWorkflowWithRequest(request *ParseWorkflowRequest) (*Workflow, derrors.Error) {
	return p.ParseWorkflow(request.Template, request.Name, request.Parameters)
}

// ParseWorkflow reads a workflow from a string, parsing the data and applying the template.
//   params:
//     content The template content with the workflow.
//     name The name of the workflow.
//     params The template parameters.
//   returns:
//     A Workflow structure.
//     An error if the workflow cannot be generated.
func (p *Parser) ParseWorkflow(content string, name string, params Parameters) (*Workflow, derrors.Error) {
	ft := template.New("Workflow: " + name).Funcs(template.FuncMap{
		"joinStringArray": func(elements []string) string {
			return "\"" + strings.Join(elements, "\",\"") + "\""
		},
	})
	commentsRegex := regexp.MustCompile("(?m)[\r\n]+^[[:blank:]]*//.*$")
	// remove comments stating with //
	templateToParse := commentsRegex.ReplaceAllString(content, "")
	ft, err := ft.Parse(templateToParse)
	if err != nil {
		return nil, derrors.NewInternalError(errors.CannotParseTemplate, err)
	}
	log.Debug().Str("template", ft.Name()).Msg("Executing template")
	// output buffer for the JSON content
	buf := new(bytes.Buffer)
	err = ft.Execute(buf, params)
	if err != nil {
		return nil, derrors.NewInternalError(errors.CannotApplyTemplate, err)
	}
	jsonPayload := buf.String()
	return p.ParseJSON(jsonPayload, name)
}

// ParseJSON reads a workflow from a JSON string, parsing the data and applying the template.
//   params:
//     jsonPayload The JSON content with the workflow.
//     name The name of the workflow.
//   returns:
//     A Workflow structure.
//     An error if the workflow cannot be generated.
func (p *Parser) ParseJSON(jsonPayload string, name string) (*Workflow, derrors.Error) {
	passwordRegex := regexp.MustCompile("\"password\":\".*\",")
	privateKeyRegex := regexp.MustCompile("\"privateKey\":\".*\"")
	redactedJSON := passwordRegex.ReplaceAllString(jsonPayload, "\"password\":\"REDACTED\",")
	redactedJSON = privateKeyRegex.ReplaceAllString(redactedJSON, "\"privateKey\":\"REDACTED\"")
	redactedJSON = strings.Replace(redactedJSON, "\n", "", -1)
	redactedJSON = strings.Replace(redactedJSON, "\t", "", -1)
	log.Debug().Str("redactedJSON", redactedJSON).Msg("Workflow to be parsed")

	var aux rawWorkflow
	if err := json.Unmarshal([]byte(jsonPayload), &aux); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err).WithParams(jsonPayload)
	}

	result := make([]entities.Command, 0)
	for index, raw := range aux.Commands {
		toShow := string(raw)
		toShow = passwordRegex.ReplaceAllString(toShow, "\"password\":\"REDACTED\",")
		toShow = privateKeyRegex.ReplaceAllString(toShow, "\"privateKey\":\"REDACTED\"")
		log.Debug().Int("index", index).Str("cmd", toShow).Msg("processing cmd")
		cmd, err := p.cmdParser.ParseCommand(raw)
		if err != nil {
			return nil, err
		}
		result = append(result, *cmd)
	}

	return NewWorkflow(name, aux.Description, result), nil
}
