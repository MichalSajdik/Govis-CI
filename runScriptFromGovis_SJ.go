package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func (singleJob *SingleJob) RunScriptFromGovis(govisScriptInArray GovisScriptStructYaml) error {
	var err, exitError error
	if !ScriptRights {
		return errors.New("No rights to run direct script")
	}
	err = singleJob.runScriptOnPlatform(govisScriptInArray.Script)

	if err != nil {
		exitError = err
		singleJob.LogInfo("Fail: " + err.Error())
	}
	err = singleJob.RunAbapMerge(govisScriptInArray.AbapMerge)
	if err != nil {
		exitError = err
		singleJob.LogInfo("Fail: " + err.Error())
	}
	return exitError
}

func (singleJob *SingleJob) runScriptOnPlatform(Script []string) error {
	var err error
	if PlatformIsWindows() {
		err = singleJob.RunScriptFromArrayInWindows(Script)
	} else {
		err = singleJob.RunScriptFromArray(Script)
	}
	return err
}

// RunScriptFromArray runs scripts from @govisScriptInArray with powershell.exe in @location repository
// returns first error during run
func (singleJob *SingleJob) RunScriptFromArray(script []string) error {
	var err error
	var exitErr error
	execArgs := []string{"-c", "PLACEHODLER"}

	nameOfCommand := "/bin/sh"
	for _, cmd := range script {
		singleJob.LogInfo(fmt.Sprintf("Starting run for >> %s", cmd))
		execArgs[1] = "(" + cmd + ")"
		err = singleJob.RunScript(nameOfCommand, execArgs)
		if (err != nil) && (exitErr == nil) {
			exitErr = err
		}

	}
	singleJob.LogInfo("All scripts in array finished")
	return exitErr
}

func (singleJob *SingleJob) RunScriptFromArrayInWindows(script []string) error {
	var err error
	var exitErr error
	for i := 0; i < len(script); i++ {
		singleJob.LogInfo(fmt.Sprintf("Starting run for >> %s", script[i]))
		words := strings.Fields(script[i])

		err = singleJob.RunScript(words[0], words[1:])

		if (err != nil) && (exitErr == nil) {
			exitErr = err
		}
	}

	singleJob.LogInfo("All scripts in array finished")

	return exitErr
}

//runs podman container, with infinit sleep, then runs all scripts in it
//closes container
func (singleJob *SingleJob) RunScriptInPodman(script []string, link_host string) error {
	var exitErr error
	nameOfCommand := "podman"
	repoDir := singleJob.location + "/" + singleJob.commit.RepoName
	argsOfCommand := []string{"run", "-d", "--add-host", link_host, "-v", "/opt/sap:/opt/sap:ro",
		"-v", repoDir + ":" + repoDir, "--workdir", repoDir,
		"--name", singleJob.branchName, "username/abapops", "sleep", "infinity"}

	singleJob.LogInfo(fmt.Sprintf("Starting podman to run commands in: %s %s\n", nameOfCommand, argsOfCommand))
	startIsolation := exec.Command(nameOfCommand, argsOfCommand...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	startIsolation.Stdout = &out
	startIsolation.Stderr = &stderr
	err := startIsolation.Run()
	if err != nil {
		singleJob.LogWarn(fmt.Sprint(err) + ": " + stderr.String())
		singleJob.LogWarn("Podman starting failed with name " + singleJob.branchName + out.String())
		singleJob.LogWarn("#########")
		return err
	} else {
		singleJob.LogInfo(fmt.Sprint(err) + ": " + stderr.String())
		singleJob.LogInfo("Started podman with name " + singleJob.branchName + out.String())
		singleJob.LogInfo("#########")
	}

	execArgs := []string{"exec", "-i", singleJob.branchName, "/bin/sh", "-c", "PLACEHODLER"}
	connectionRC := "(. <filename>.sh.rc; "
	for _, cmd := range script {
		singleJob.LogInfo(fmt.Sprintf("Starting run for >> %s", cmd))
		execArgs[5] = connectionRC + cmd + ")"
		err = singleJob.RunScript(nameOfCommand, execArgs)
		if (err != nil) && (exitErr == nil) {
			exitErr = err
		}

	}
	singleJob.LogInfo("All scripts in array finished")

	//exit

	//TODO:zombies might stay in podman, service to check them and kill them
	exitIsolation := exec.Command("podman", "rm", singleJob.branchName, "-f")

	err = exitIsolation.Run()
	if err != nil {
		singleJob.LogWarn(fmt.Sprint(err) + ": " + stderr.String())
		singleJob.LogWarn("Podman starting failed with name " + singleJob.branchName + out.String())
		singleJob.LogWarn("#########")
	} else {
		singleJob.LogInfo(fmt.Sprint(err) + ": " + stderr.String())
		singleJob.LogInfo("Started podman with name " + singleJob.branchName + out.String())
		singleJob.LogInfo("#########")
	}
	return exitErr

}

// RunScript calls one of the responsible functions for running script based on platform
func (singleJob *SingleJob) RunScript(nameOfCommand string, argsOfCommand []string) error {
	script := "'" + strings.Join(argsOfCommand, "'") + "'"
	singleJob.LogInfo(fmt.Sprintf("App: %s; Script:%s\n", nameOfCommand, script))
	cmd := exec.Command(nameOfCommand, argsOfCommand...)

	cmd.Dir = singleJob.locationOfRepository

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return err
	}
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case <-time.After(ScriptTimeout):
		// Timeout happened first, kill the process and print a message.
		cmd.Process.Kill()
		singleJob.LogInfo("Command timed out")
		err = errors.New("command timed out")
	case err = <-done:
		if err != nil {
			singleJob.LogInfo(fmt.Sprint(err) + ": " + stderr.String())
			singleJob.LogInfo("Result: " + out.String())
			singleJob.LogInfo("#########")
		} else {
			singleJob.LogInfo(fmt.Sprint(err) + ": " + stderr.String())
			singleJob.LogInfo("Result: " + out.String())
			singleJob.LogInfo("#########")
		}
	}
	return err
}
