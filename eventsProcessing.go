package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v27/github"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EventsProcessing continuous implementation for all POST messages individually implemented only for events from webhooks
func EventsProcessing(w http.ResponseWriter, r *http.Request) {
	LoggerServerINFO("Post message has come")
	event := r.Header.Get("X-Github-Event")
	payloadJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		LoggerServerWARN(fmt.Sprintf("payloadJSON >> %s\n", payloadJSON))
		LoggerServerERROR(fmt.Sprintf("error while converting payload context data >> %s", err.Error()))
	}
	defer r.Body.Close()

	LoggerServerINFO(fmt.Sprintf("header event >> %s\n", event))

	signature := r.Header.Get("X-Hub-Signature")
	err = github.ValidateSignature(r.Header.Get("X-Hub-Signature"), payloadJSON, []byte(GithubSecret))
	if err != nil {
		LoggerServerWARN(fmt.Sprintf("X-Hub-Signature >> %s\n", signature))
		LoggerServerWARN(fmt.Sprintf("Problem with validation of signature: %s\n", err))
		w.WriteHeader(403)
		return
	}
	LoggerServerINFO(fmt.Sprintf("Validation of signature was successful"))

	var htmlURL string
	var branchName string
	var commit Commit

	switch event {
	case PushEvent:
		htmlURL, branchName, commit, err = parsePushEvent(payloadJSON)

	case PullRequestEvent:
		htmlURL, branchName, commit, err = parsePullRequestEvent(payloadJSON)

	case "ping":
		LoggerServerINFO(fmt.Sprintf("New hook ping event"))
		w.WriteHeader(200)
		LoggerServerINFO(fmt.Sprintf("Response OK is being send\n"))
		return // sending response OK to git
	default:
		LoggerServerINFO(fmt.Sprintf("Unknown Event >> %s\n", event))
		w.WriteHeader(501) // 501 Not Implemented
		return
	}

	LoggerServerINFO(fmt.Sprintf("Owner: %s\n", commit.Owner))
	LoggerServerINFO(fmt.Sprintf("RepoName: %s\n", commit.RepoName))
	LoggerServerINFO(fmt.Sprintf("IdOfCommit: %s\n", commit.IdOfCommit))
	LoggerServerINFO(fmt.Sprintf("html_url >> %s\n", htmlURL))
	LoggerServerINFO(fmt.Sprintf("branch name >> %s\n", branchName))
	LoggerServerINFO(fmt.Sprintf("StatusesUrl >> %s\n", commit.StatusesUrl))

	directoryName := fmt.Sprintf("%s_%s_%s", commit.RepoName, time.Now().UTC().Truncate(time.Second), commit.IdOfCommit[:7])
	directoryName = strings.ReplaceAll(directoryName, " ", "_") //windows compatibility
	directoryName = strings.ReplaceAll(directoryName, ":", "_") //windows compatibility

	location := filepath.Join(GlobalPath, ResultsDir, directoryName)
	err = os.Mkdir(location, os.ModePerm)
	if err != nil {
		log.Fatalf("error creating directory: %v", err)
	}
	logFile := OpenLogFile(location)
	singleJob := SingleJobCreate(branchName, location, logFile, commit)

	gitAuthorizedURL, err := GetAuthorizedURL(htmlURL, os.Getenv("GOVIS_BASIC_AUTH"))
	if err != nil {
		LoggerServerWARN(fmt.Sprintf("Can't add token to clone url%s\n", err.Error()))
	}
	singleJob.GitClone(gitAuthorizedURL)

	go func() {
		err := singleJob.WorkWithRepository()
		if err != nil {
			LoggerServerWARN(fmt.Sprintf("Problem with task: %v >> %v\n", commit.IdOfCommit, err))
		}
	}()

	statusCode, message := Web.checkConfigurationOfEnvironment()
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(message))

	LoggerServerINFO(fmt.Sprintf("Response OK is being send\n"))
}

func parsePullRequestEvent(payloadJSON []byte) (string, string, Commit, error) {
	var payload github.PullRequestEvent
	err := json.Unmarshal(payloadJSON, &payload)
	if err != nil {
		LoggerServerINFO(fmt.Sprintf("LogError while parsing payload: %s\n", err.Error()))
	}

	branchName := payload.PullRequest.Head.GetRef()
	htmlURL := payload.PullRequest.Head.Repo.GetHTMLURL()
	var commit Commit
	commit.Owner = payload.PullRequest.Head.Repo.Owner.GetLogin()
	commit.RepoName = payload.PullRequest.Head.Repo.GetName()
	commit.IdOfCommit = payload.PullRequest.Head.GetSHA()
	commit.StatusesUrl = strings.ReplaceAll(payload.PullRequest.GetStatusesURL(), "{sha}", commit.IdOfCommit)

	return htmlURL, branchName, commit, err
}

func parsePushEvent(payloadJSON []byte) (string, string, Commit, error) {
	var payload github.PushEvent
	err := json.Unmarshal(payloadJSON, &payload)
	if err != nil {
		log.Println(err)
	}
	htmlURL := payload.Repo.GetHTMLURL()
	branchName := payload.GetRef()
	branchName = RemovePrefix("refs/heads/", branchName) // deletes prefix "refs/heads/"
	var commit Commit
	commit.Owner = payload.GetRepo().GetOwner().GetName()
	commit.RepoName = payload.GetRepo().GetName()
	commit.IdOfCommit = payload.HeadCommit.GetID()
	commit.StatusesUrl = strings.ReplaceAll(payload.Repo.GetStatusesURL(), "{sha}", commit.IdOfCommit)

	return htmlURL, branchName, commit, err
}
