/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// Enumeration of the supported commands.

package entities

// Exec command to execute local commands in the system.
const Exec = "exec"

// SCP command to copy data to remote hosts.
const SCP = "scp"

// SSH command to execute commands on remote hosts.
const SSH = "ssh"

// Logger command to add logging information to the workflow log.
const Logger = "logger"

// Fail command that aborts the workflow execution.
const Fail = "fail"

// Sleep commands that waits for a given ammount of time.
const Sleep = "sleep"

// GroupCmd command to sequentialize subsets of commands.
const GroupCmd = "group"

// ParallelCmd command to execute several commands in parallel.
const ParallelCmd = "parallel"

// TryCmd command that tries to execute a command, and in case of failure executes an alternative one.
const TryCmd = "try"

// ProcessCheck command to determine if a process is running on a given machine.
const ProcessCheck = "processCheck"


// RKEInstall command to launch the installation of a new cluster with RKE.
const RKEInstall = "rkeInstall"
const RKERemove = "rkeRemove"

const LaunchComponents = "launchComponents"

