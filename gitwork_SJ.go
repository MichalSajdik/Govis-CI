package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// GitClone does a simple git clone command for @htmlUrl that's provided and saves it in our repo
// git config --global core.compression 0 >> may be needed to be set
func (singleJob *SingleJob) GitClone(htmlURL string) {
	os.Setenv("GIT_TERMINAL_PROMPT", "0")
	// try `"--depth", "1",` add to cmd if error // clones only master
	cmd := exec.Command("git", "clone", htmlURL)
	cmd.Dir = singleJob.location
	out, err := cmd.Output()
	os.Unsetenv("GIT_TERMINAL_PROMPT")

	if err != nil {
		singleJob.LogWarn(fmt.Sprintf("git clone was unsuccessful\n"))
		return
	}

	singleJob.LogInfo(fmt.Sprintf("Repository was cloned >> %s\n", out))

	err = os.Chmod(singleJob.location+"/"+singleJob.commit.RepoName, 0777)
	if err != nil {
		singleJob.LogError("LogError during git clone" + err.Error())
	}
}

// GitCheckoutBranch shortcut for "git checkout @branchName" command
func (singleJob *SingleJob) GitCheckoutBranch() {
	script := "git checkout " + singleJob.branchName + "@at=" + singleJob.locationOfRepository // for logger only

	cmd := exec.Command("git", "checkout", singleJob.branchName)
	cmd.Dir = singleJob.locationOfRepository
	singleJob.LogInfo(fmt.Sprintf("Location >> %s ", singleJob.locationOfRepository))
	out, err := cmd.Output()
	if err != nil {
		singleJob.LogWarn(fmt.Sprintf("git checkout %s >> %s >> with command >> %s", singleJob.branchName, err, script))
	} else {
		singleJob.LogInfo(fmt.Sprintf("%s was checkout-ed >> %s >> with command >> %s", singleJob.branchName, out, script))
	} // printing output from git checkout 'branchName'
}

//GitSetStatusForCommit uses REST api to set status for given @idOfCommit
func (singleJob *SingleJob) GitSetStatusForCommit(postMessage statusCheck) {
	LoggerServerINFO(fmt.Sprintf("Setting Status For Commit: %s\n", singleJob.commit.IdOfCommit))
	jsonStatusCheckCreation, _ := json.Marshal(postMessage)
	LoggerServerINFO(fmt.Sprintf("%s", jsonStatusCheckCreation))

	urlToRepository := singleJob.commit.StatusesUrl
	payloadInString := fmt.Sprintf(`{"state": "%s",
                           "target_url": "%s" ,
                           "description": "%s",
                           "context": "%s"}`, postMessage.State, postMessage.TargetURL, postMessage.Description, postMessage.Context)

	payload := strings.NewReader(payloadInString)

	req, _ := http.NewRequest("POST", urlToRepository, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(BasicAuthToken)))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		LoggerServerWARN(fmt.Sprintf("Set commit status failed: %s\n", err))
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	LoggerServerINFO(fmt.Sprintf("curl %s\n", urlToRepository))
	LoggerServerINFO(fmt.Sprintf("Response >> %v\n", res))
	LoggerServerINFO(fmt.Sprintf("Body of Response >> %s\n", string(body)))
}
