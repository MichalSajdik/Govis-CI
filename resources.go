package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var GovisVersionCommitId string

// PayloadUrl is set in webhooks settings on github
const PayloadUrl = "/payload"

const Port = "8000"

// GovisYaml is const name of source file for scripts to be run
const GovisYAML = ".govis.yaml"

// PushEvent contains const name from github for PushEvent
const PushEvent = "push"

// PullRequestEvent contains const name from github for PullRequestEvent
const PullRequestEvent = "pull_request"

// LogFileName is name of log file for each push/pull event
const LogFileName = "result.log"

// AccessLogFileName is name of log file that contains all attempts to communicate with our `Server`
const AccessLogFileName = "AccessLog.log"

// ServerLogFileName is name of log file that contains all other happenings on `Server`
const ServerLogFileName = "ServerLog.log"

// Server that we use for sharing results and accepting webhooks
var Server http.Server

// AccessibleFiles contains files allowed to view, `""` stands for `index.html`
// var AccessibleFiles = []string{"/",GovisWebUrl+"", GovisWebUrl+"/index2.html", GovisWebUrl+"/AccessLog.txt", GovisWebUrl+"/ServerLog.txt"}

// GlobalPath contains in what path to our folder
var GlobalPath, _ = os.Getwd()

// ServerLogFile is file for server logs
var ServerLogFile = func() os.File {
	locationOfLogFile := filepath.Join(GlobalPath, ServerLogFileName)
	logFile, err := os.OpenFile(locationOfLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	fmt.Println("LogFile at >> ", locationOfLogFile, "Printed from OpenLogFile() in helper.go")
	return *logFile
}()

// ResultsDir is name of folder, where results from `scripts` are saved
const ResultsDir = "results"

// TestResults is name of folder, where results from `scripts` are saved
const TestResultsDir = "testResults"

const EnvVerboseLogging = "GOVIS_LOG"
const EnvNameApi = "GOVIS_API"
const EnvNameBasicAuthToken = "GOVIS_BASIC_AUTH"
const EnvNameGithubSecret = "GOVIS_SECRET"
const EnvNameScriptRights = "GOVIS_SCRIPT_RIGHTS"

// Api is const string value for each company specific , defined as env variable, default is "https://api.github.com/"
var Api = func() string {
	tmpApi := os.Getenv(EnvNameApi)
	if tmpApi == "" {
		return "https://api.github.com"
	} else {
		return tmpApi
	}
}()

// Script Rights boolean is true if you set GOVIS_SCRIPT_RIGHTS env var, it allows govis to run script from .govis.yaml
// under user govis is started from
var ScriptRights = func() bool {
	if ScriptRightsEnv == "" {
		return false
	} else {
		return true
	}
}()

var EnvSecureMode = os.Getenv("GOVIS_SECURE_MODE")

func GetSecureMode() bool {
	if EnvSecureMode == "negative" {
		return false
	}
	return true
}

var VerboseLogging = func() bool {
	if os.Getenv(EnvVerboseLogging) == "" {
		return false
	}
	return true
}()

var ScriptRightsEnv = os.Getenv(EnvNameScriptRights)

// BasicAuthToken, described in readme.md
var BasicAuthToken = os.Getenv(EnvNameBasicAuthToken)

// GithubSecret, described in readme.md
var GithubSecret = os.Getenv(EnvNameGithubSecret)

var HostnameEnv = os.Getenv("GOVIS_HOSTNAME")
var PortEnv = os.Getenv("GOVIS_PORT")

var EnvScriptTimeout = os.Getenv("GOVIS_SCRIPT_TIMEOUT")

const DefaultScriptTimeout = 300

var ScriptTimeout = func() time.Duration {
	t := EnvScriptTimeout
	if t == "" {
		return time.Duration(DefaultScriptTimeout) * time.Second
	}
	timeout, _ := strconv.Atoi(t)
	return time.Duration(timeout) * time.Second
}()

var Web = GovisWeb{
	Name:          "User",
	Time:          time.Now().Format(time.Stamp),
	HostnameValue: HostnameEnv,
}

// GovisScriptStructYaml is struct for parsing `GovisYAML` file
type GovisScriptStructYaml struct {
	Script    []string `yaml:"script"`
	AbapMerge []string `yaml:"abapMerge"`
}

type statusCheck struct {
	State       string `json:"state"`
	TargetURL   string `json:"target_url"`
	Description string `json:"description"`
	Context     string `json:"context"`
}

// Method to set status to failure
func (s *statusCheck) Success() {
	s.State = "success"
	s.Description = "Build was successful."
}

// Method to set status to failure
func (s *statusCheck) Failure() {
	s.State = "failure"
	s.Description = "Build failed."
}

// Method to set status to pending
func (s *statusCheck) Pending() {
	s.State = "pending"
	s.Description = "Build is pending."
}

// Factory function that returns empty status
func NewStatusCheck(resultFolder string) *statusCheck {
	targetUrl := os.Getenv("GOVIS_STATUS_BASEURL")
	if targetUrl == "" {
		targetUrl = "http://" + GetGovisAddress()
	}
	s := new(statusCheck)
	s.TargetURL = targetUrl + resultFolder
	s.Context = "continuous-integration/govis"
	return s
}

func NewSuccessStatus(resultFolder string) *statusCheck {
	s := NewStatusCheck(resultFolder)
	s.Success()
	return s
}

func NewFailedStatus(resultFolder string) *statusCheck {
	s := NewStatusCheck(resultFolder)
	s.Failure()
	return s
}

func NewPendingStatus(resultFolder string) *statusCheck {
	s := NewStatusCheck(resultFolder)
	s.Pending()
	return s
}

type Commit struct {
	Owner       string
	RepoName    string
	IdOfCommit  string
	StatusesUrl string
}
