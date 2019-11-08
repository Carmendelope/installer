/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package errors

// InvalidEntity message indicating that the associated entity cannot be validated.
const InvalidEntity = "invalid entity, check mandatory fields"

// MarshalError message to indicate errors with the json.Marshal operation.
const MarshalError = "cannot marshal structure"

// UnmarshalError message to indicate errors with the json.Unmarshal operation.
const UnmarshalError = "cannot unmarshal structure"

// OpFail message to indicate that a complex operation has failed.
const OpFail = "operation failed"

// MissingRESTParameter message to indicate that a required parameter is missing.
const MissingRESTParameter = "missing rest parameter"

// IOError message to indicate that an I/O has occurred reading/writing from file, socket, etc.
const IOError = "I/O error"

// HTTPConnectionError message to indicate that the communication with an external entity has failed.
const HTTPConnectionError = "HTTP connection error"

// SSHConnectionError message to indicate that the communication with an external entity using SSH has failed.
const SSHConnectionError = "SSH connection error"

// Templates

// CannotParseTemplate error to indicate that the template file contains invalid syntax.
const CannotParseTemplate = "cannot parse workflow template file"

const CannotParseRKETemplate = "cannot parse RKE cluster template file"

// CannotApplyTemplate error to indicate that the template cannot be applied with the given parameters.
const CannotApplyTemplate = "cannot apply template parameters"

//Parameters

// CannotParseParameters error to indicate that the parameters input file cannot be read.
const CannotParseParameters = "cannot parse parameters file"

// Workflows

// CannotWriteWorkflowFile to indicate that the output workflow cannot be generated.
const CannotWriteWorkflowFile = "cannot write workflow file"

// CannotReadWorkflowFile error to indicate that the workflow file cannot be read.
const CannotReadWorkflowFile = "cannot read workflow file"

// WorkflowWithoutCommands error to indicate that the specified workflow does not contain any command.
const WorkflowWithoutCommands = "attempting to execute workflow without commands"

// InvalidWorkflowState error to indicate that the workflow current state does not follow the expected transitions.
const InvalidWorkflowState = "invalid workflow state"

// WorkflowDoesNotExists error to indicate the target workflow does not exists.
const WorkflowDoesNotExists = "workflow does not exists"

// WorkflowAlreadyExists error to indicate the target workflow is already in the system.
const WorkflowAlreadyExists = "workflow already exists"

// WorkflowExecutionFailed error to indicate that the execution of the workflow failed.
const WorkflowExecutionFailed = "workflow execution failed"

// Commands

// UnsupportedCommandType error to indicate that the selected command type is not supported and cannot be executed.
const UnsupportedCommandType = "unsupported command type"

// UnsupportedCommand error to indicate that the selected command is not supported.
const UnsupportedCommand = "unsupported command"

// CannotExecuteSyncCommand to indicate that the synchronous command execution failed.
const CannotExecuteSyncCommand = "cannot execute synchronous command"

// InvalidCommandIndex to indicate that the command to be executed is not defined in the workflow.
const InvalidCommandIndex = "command index outside of bounds of the current workflow"

// DuplicatedIDCommand to indicate that the command already exists.
const DuplicatedIDCommand = "command id is duplicated"

// NotExistCommand to indicate that the target command does not exists.
const NotExistCommand = "command id does not exist"

// InvalidCommandParameters to indicate that the command was expecting a set of parameters that are not present.
const InvalidCommandParameters = "missing/invalid command parameters"

// Assets

// FileDoesNotExist to indicate that a file does not exist in the given path.
const FileDoesNotExist = "target file does not exists"

// ExpectingFile to indicate that the target path is a directory.
const ExpectingFile = "target path is a directory, expecting file"

// CorePackageInvalidName indicates that the core package name does not match the expected format.
const CorePackageInvalidName = "core packages should have the core- prefix, and the be compressed as tar.gz"

// Parameters

// ParameterDoesNotExists error to indicate that the requested parameter does not exists.
const ParameterDoesNotExists = "requested parameter does not exists"

// InvalidNodeConf error to indicate that the nodes file is not valid.
const InvalidNodeConf = "nodes.json must define a set of masters and valid credentials"

// InvalidNumMaster error to indicate that the number of master nodes is not supported.
const InvalidNumMaster = "invalid number of master nodes, expecting 1 or 3"

const InvalidYAML = "invalid YAML file"

const InvalidNumberOfTokens = "invalid number of teleport tokens returned"
