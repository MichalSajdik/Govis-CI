package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SingleJob struct {
	branchName           string
	location             string
	locationOfRepository string
	commit               Commit
	LoggingComponent     //logFile
}

// SingleJobCreate returns pointer to filled SingleJob struct
func SingleJobCreate(branchName, location string, logFile os.File, commit Commit) *SingleJob {
	singleJob := new(SingleJob)
	singleJob.branchName = branchName
	singleJob.location = location
	singleJob.locationOfRepository = filepath.Join(location, commit.RepoName)
	singleJob.commit = commit
	singleJob.logFile = logFile
	return singleJob
}

// workWithRepository basically second main
func (singleJob *SingleJob) WorkWithRepository() error {
	logFileResultsPath := strings.Replace(singleJob.logFile.Name(), GlobalPath, "", 1)
	//windows path replace
	resultLog := strings.Replace(logFileResultsPath, `\`, "/", -1)

	pendingStatusCheck := NewPendingStatus(resultLog)
	singleJob.GitSetStatusForCommit(*pendingStatusCheck)

	singleJob.GitCheckoutBranch()
	failedStatus := NewFailedStatus(resultLog)
	successStatusCheck := NewSuccessStatus(resultLog)
	govisScriptInArray, err := singleJob.ParseGovis()
	if err != nil {
		singleJob.LogWarn("LogError during parsing GOVIS: " + err.Error())
		singleJob.GitSetStatusForCommit(*failedStatus)
		return err
	}
	errFromGovisScript := singleJob.RunScriptFromGovis(govisScriptInArray)
	if errFromGovisScript != nil {
		singleJob.GitSetStatusForCommit(*failedStatus)
		return errFromGovisScript
	}

	singleJob.LogInfo("Success")
	// setting status for commit
	singleJob.GitSetStatusForCommit(*successStatusCheck)

	return nil
	// TODO implement automatic RemoveRepository
	//if PlatformIsWindows() {
	//	RemoveRepository(location)
	//}else{
	//	removeRepositoryBash(location)
	//}

}

func (singleJob *SingleJob) RunAbapMerge(AbapMerge []string) error {
	if len(AbapMerge) == 0 {
		singleJob.LogInfo("AbapMerge is empty")
		return nil
	}
	var abapMergeScriptBefore = []string{"podman run -d --name nwabap -v /sys/fs/cgroup/:/sys/fs/cgroup -v /opt/sap:/opt/sap --entrypoint /usr/lib/systemd/systemd --privileged -h vhcalnplci nwabap754_sapintern --system",
		"podman exec -i nwabap /opt/sap/boot.sh"}

	var abapMergeScriptAfter = []string{"podman rm -f nwabap"}

	var err, exitError error
	err = singleJob.RunScriptFromArray(abapMergeScriptBefore)
	if err != nil {
		exitError = err
		singleJob.LogWarn("During running abapMergeScriptBefore >>  Fail: " + err.Error())
	} else {
		var cmd = exec.Command("podman", "inspect", "nwabap", "--format", "{{ .NetworkSettings.IPAddress }}")
		var out, err = cmd.Output()
		if err != nil {
			exitError = err
			singleJob.LogWarn("After abapMergeScriptBefore (get nwabap IP) >>  Fail: " + err.Error())
		} else {
			err = singleJob.RunScriptInPodman(AbapMerge, "vhcalnplci:"+strings.TrimSpace(string(out)))
			if err != nil {
				exitError = err
				singleJob.LogWarn("During running AbapMerge >>  Fail: " + err.Error())
			}
		}
	}

	err = singleJob.RunScriptFromArray(abapMergeScriptAfter)
	if err != nil {
		if exitError == nil {
			exitError = err
		}
		singleJob.LogWarn("During running abapMergeScriptAfter >>  Fail: " + err.Error())
	}
	return exitError
}

// Call RemoveRepository function based on  platform
func (singleJob *SingleJob) RemoveRepository() {
	if PlatformIsWindows() {
		removeRepositoryPS(singleJob.location, singleJob.logFile)
	} else {
		removeRepositoryBash(singleJob.location, singleJob.logFile)
	}

}

func removeRepositoryPS(location string, logFile os.File) {
	// remove local files after working with them ;;; Remove-Item -LiteralPath "location" -Force -Recurse

	scriptForPowerShell := "Remove-Item -LiteralPath '" + location + "' -Force -Recurse"
	out, err := exec.Command(GetPathToPowerShell(), scriptForPowerShell).Output()
	if err != nil {
		LoggerServerINFO(fmt.Sprintf("Remove-Item -LiteralPath '%s' -Force -Recurse >> %s", location, err))
	} else {
		LoggerServerINFO(fmt.Sprintf("%s was deleted >> %s", location, out))
	}
}

func removeRepositoryBash(location string, logFile os.File) {
	cmd := exec.Command("rm", "-rf", location)
	out, err := cmd.Output()
	if err != nil {
		LoggerServerINFO(fmt.Sprintf("rm -rf %s >> %s", location, err))
	} else {
		LoggerServerINFO(fmt.Sprintf("%s was deleted >> %s", location, out))
	}
}
