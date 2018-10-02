/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// Parallel command
// Permits the parallel execution of a set of commands.
//
// {"type":"sync", "name": "parallel", "commands": [{"type":...},{"type":...}]}

package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"
	"sync"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/nalej/installer/internal/pkg/workflow/handler"
)

// Parallel structure with the commands to be executed.
type Parallel struct {
	sync.Mutex
	entities.GenericSyncCommand
	Description     string             `json:"description"`
	MaxParallelism  int                `json:"maxParallelism"`
	Commands        []entities.Command `json:"commands"`
	commandHandler  handler.CommandHandler
	commandResults  map[string]entities.CommandResult
	executionErrors map[string]derrors.Error
	finishChannel   chan string
}

// ParallelFromJSON structure with helper RawMessage to parse
// target commands.
type ParallelFromJSON struct {
	entities.GenericCommand
	Description    string            `json:"description"`
	MaxParallelism int               `json:"maxParallelism"`
	Commands       []json.RawMessage `json:"commands"`
}

// ToParallel transforms a JSON group into a Parallel structure parsing the required commands.
func (pfj *ParallelFromJSON) ToParallel() (*Parallel, derrors.Error) {
	p := NewCmdParser()
	cmds := make([]entities.Command, 0)
	for _, toParse := range pfj.Commands {
		toAdd, err := p.ParseCommand(toParse)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, *toAdd)
	}
	return NewParallel(pfj.Description, pfj.MaxParallelism, cmds), nil
}

// NewParallel creates a new Parallel structure with a given description and associated commands.
func NewParallel(description string, maxParallelism int, cmds []entities.Command) *Parallel {
	return &Parallel{
		GenericSyncCommand: *entities.NewSyncCommand(entities.ParallelCmd),
		Description:        description,
		MaxParallelism:     maxParallelism,
		Commands:           cmds, commandHandler: handler.GetCommandHandler(),
		commandResults:  make(map[string]entities.CommandResult),
		executionErrors: make(map[string]derrors.Error),
		finishChannel:   make(chan string)}
}

// NewParallelFromJSON creates a new command from a raw json payload.
func NewParallelFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	pfj := &ParallelFromJSON{}
	if err := json.Unmarshal(raw, &pfj); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	toParallel, err := pfj.ToParallel()
	if err != nil {
		return nil, err
	}
	var r entities.Command = toParallel
	return &r, nil
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// Run the current command.
//   returns:
//     The CommandResult
//     An error if the command execution fails
func (p *Parallel) Run(workflowID string) (*entities.CommandResult, derrors.Error) {

	awaiting := len(p.Commands)
	initialLaunch := len(p.Commands)
	if p.MaxParallelism != 0 {
		initialLaunch = min(p.MaxParallelism, len(p.Commands))
	}
	log.Info().Str("cmd", p.CommandID).Str("description", p.Description).Msg("Executing parallel group")
	log.Debug().Int("initialLaunch", initialLaunch).Msg("MaxParallelism set")
	launched := 0
	for index := 0; index < initialLaunch; index++ {
		nextCommand := p.Commands[index]
		toExecute := nextCommand
		log.Debug().Str("Id", toExecute.ID()).Str("cmd", toExecute.String()).Msg("Launching goroutine for command execution")
		go func() {
			p.execOnBackground(workflowID, toExecute)
		}()
		launched++
	}

	log.Debug().Int("launched", launched).Int("total", awaiting).Msg("")
	failed := false
	for received := 0; received < awaiting && !failed; received++ {
		log.Debug().Int("launched", launched).Int("finished", received).Msg("")
		cmdID := <-p.finishChannel
		log.Debug().Str("cmdID", cmdID).Msg("Command finished")
		if p.executionWithErrors() {
			log.Debug().Msg("command returned error, aborting.")
			failed = true
		} else {
			if launched < awaiting {
				nextCommand := p.Commands[launched]
				toExecute := nextCommand
				log.Debug().Str("Id", toExecute.ID()).Str("cmd", toExecute.String()).Msg("Launching goroutine for command execution")
				go func() {
					p.execOnBackground(workflowID, toExecute)
				}()
				launched++
			}
		}
	}

	return p.buildCommandResult()
}

func (p *Parallel) executionWithErrors() bool {
	overallSuccess := true
	p.Lock()
	for _, r := range p.commandResults {
		overallSuccess = overallSuccess && r.Success
	}
	p.Unlock()
	log.Debug().Bool("overallSuccess ", overallSuccess).Int(" len ", len(p.executionErrors)).Msg("")
	return len(p.executionErrors) > 0 || !overallSuccess
}

func (p *Parallel) execOnBackground(workflowID string, cmd entities.Command) {
	err := p.commandHandler.AddCommand(cmd.ID(), p.ParallelCallback, p.logCallback)
	if err != nil {
		log.Warn().Str("err", err.DebugReport()).Msg("error on exec")
		p.Lock()
		p.executionErrors[cmd.ID()] = err
		p.Unlock()
		return
	}

	p.commandHandler.AddLogEntry(p.CommandID, "Executing on Parallel: "+cmd.ID())
	if cmd.Type() == entities.SyncCommandType {
		log.Debug().Str("cmd", cmd.String()).Msg("SYNC")
		result, err := cmd.(entities.SyncCommand).Run(workflowID)
		err = p.commandHandler.FinishCommand(cmd.ID(), result, err)
		if err != nil {
			log.Warn().Str("id", cmd.ID()).Str("err", err.DebugReport()).Msg("error executing sync command on parallel group")
			p.Lock()
			p.executionErrors[cmd.ID()] = err
			p.Unlock()
			// The only error returned by finish command is if the command is not found.
			p.finishChannel <- cmd.ID()
		}
	} else {
		log.Debug().Str("cmd", cmd.String()).Msg("ASYNC")
		err := cmd.(entities.AsyncCommand).Run(workflowID)
		if err != nil {
			//If the execution return errors, the executor call to the commandHandler with the error.
			err = p.commandHandler.FinishCommand(cmd.ID(), nil, err)
			if err != nil {
				log.Warn().Str("id", cmd.ID()).Str("err", err.DebugReport()).Msg("error executing async command on parallel group")
				p.Lock()
				p.executionErrors[cmd.ID()] = err
				p.Unlock()
				// The only error returned by finish command is if the command is not found.
				p.finishChannel <- cmd.ID()
			}
		}
	}
}

func (p *Parallel) logCallback(cmdID string, logEntry string) {
	//p.commandHandler.AddLogEntry(p.CommandID, fmt.Sprintf("[%s] %s", cmdID, logEntry))
	p.commandHandler.AddLogEntry(p.CommandID, fmt.Sprintf("[%s] %s", p.Description, logEntry))
}

// ParallelCallback function to be called when one of the commands being executed in parallel finishes.
func (p *Parallel) ParallelCallback(cmdID string, result *entities.CommandResult, error derrors.Error) {
	log.Debug().Str("cmdID", cmdID).Msg("received callback from parallel command ")
	if result != nil {
		p.Lock()
		p.commandResults[cmdID] = *result
		p.Unlock()
	}

	if error != nil {
		p.Lock()
		p.executionErrors[cmdID] = error
		p.Unlock()
		log.Debug().Str("cmdID", cmdID).Str("err", error.DebugReport()).Msg("Parallel command failed")
	}
	p.finishChannel <- cmdID
}

func (p *Parallel) buildCommandResult() (*entities.CommandResult, derrors.Error) {
	log.Debug().Msg("Build final command result")
	if len(p.executionErrors) > 0 {
		log.Warn().Int("numErrors", len(p.executionErrors)).Msg("Execution errors")
		toReturn := derrors.NewInternalError(errors.WorkflowExecutionFailed)
		for commandID, err := range p.executionErrors {
			toReturn = toReturn.WithParams(commandID).WithParams(err)
		}
		log.Warn().Msg(toReturn.DebugReport())
		return nil, toReturn
	}

	overallSuccess := true

	var overallOutput bytes.Buffer
	for key, value := range p.commandResults {
		overallSuccess = overallSuccess && value.Success
		overallOutput.WriteString("Output of " + key + "\n" + value.Output + "\n")
	}

	return entities.NewCommandResultNoShow(overallSuccess, overallOutput.String(), nil), nil

}

// String obtains a string representation
func (p *Parallel) String() string {
	cmdNames := make([]string, 0)
	for _, cmd := range p.Commands {
		cmdNames = append(cmdNames, cmd.Name())
	}
	if p.MaxParallelism == 0 {
		return fmt.Sprintf("SYNC Parallel %s actions: %s",
			p.Description, strings.Join(cmdNames, ", "))
	}
	return fmt.Sprintf("SYNC Parallel %s maxParallelism: %d actions: %s",
		p.Description, p.MaxParallelism, strings.Join(cmdNames, ", "))
}

// PrettyPrint returns a simple space indexed string.
func (p *Parallel) PrettyPrint(identation int) string {
	cmdNames := make([]string, 0)
	for _, cmd := range p.Commands {
		cmdNames = append(cmdNames, cmd.PrettyPrint(identation+2))
	}
	return fmt.Sprintf("%sSYNC Parallel %s maxParallelism: %d actions:\n%s\n",
		strings.Repeat(" ", identation), p.Description, p.MaxParallelism, strings.Join(cmdNames, "\n"))
}

// UserString returns a simple string representation of the command for the user.
func (p *Parallel) UserString() string {
	cmdNames := make([]string, 0)
	for _, cmd := range p.Commands {
		cmdNames = append(cmdNames, cmd.Name())
	}
	return fmt.Sprintf("Parallel group: %s actions: %s", p.Description, strings.Join(cmdNames, ", "))
}
