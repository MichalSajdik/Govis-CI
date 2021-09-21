package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

var PathToDirectoryForTest = filepath.Join(GlobalPath, TestResultsDir)

const urlProvided = "<git remote repository with token>" // not to be changed
const badUrlProvided = "<git-corp-link>:FXUBRQ-QE/Govis-C"
const nameOfTestDir = "fileForTest" // @fileForTest directory should be present in @urlProvided repository

type TestSuiteForGitWorkSJ struct {
	suite.Suite
	TestSuiteBase
	specialFolder string
}

func (suiteSJ *TestSuiteForGitWorkSJ) SetupSuite() {
	repoName, _ := GetRepoName(urlProvided)
	commit := Commit{"<user name>", repoName, "<repo name>", "<link to commit id>"}
	logFile := OpenLogFile(PathToDirectoryForTest)

	suiteSJ.specialFolder = "ToBeSpecified"
	suiteSJ.singleJob = SingleJobCreate("branchForTestingDontChange", "ToBeSpecified", logFile, commit)
}

func TestGitWorkSJTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuiteForGitWorkSJ))
}

// Test_GitClone tests these cases:
// repository should not exist before calling GitClone
// repository should exist after calling GitClone
func (suiteSJ *TestSuiteForGitWorkSJ) Test_GitClone_And_RemoveRepository() {

	// setting location for test files
	suiteSJ.specialFolder = "TestGitCloneAndRemoveRepository"
	suiteSJ.singleJob.location = filepath.Join(PathToDirectoryForTest, suiteSJ.specialFolder)
	suiteSJ.singleJob.locationOfRepository = filepath.Join(suiteSJ.singleJob.location, suiteSJ.singleJob.commit.RepoName)

	// repository should not exist before calling GitClone
	os.Mkdir(filepath.Join(GlobalPath, TestResultsDir, suiteSJ.specialFolder), os.ModePerm)
	fmt.Println(suiteSJ.singleJob.locationOfRepository)
	err := Ls(suiteSJ.singleJob.locationOfRepository)

	if err == nil {
		fmt.Println("location of existing repository >> " + suiteSJ.singleJob.location)
		suiteSJ.FailNow("Cloned repository already exists: %s", err.Error())
	}
	////////////////////////////////////////////
	suiteSJ.singleJob.GitClone(urlProvided)
	////////////////////////////////////////////
	// repository should exist after calling GitClone
	err = Ls(suiteSJ.singleJob.locationOfRepository)

	suiteSJ.Equal(nil, err, "Repository was not cloned")

	// cleanup
	suiteSJ.singleJob.RemoveRepository()
	err = Ls(suiteSJ.singleJob.locationOfRepository)
	if err == nil {
		suiteSJ.FailNow("Repository was not removed: %s", err.Error())
	}
}

// Test_GitClone_WithBadUrl tests these cases:
// calling GitClone on wrong url should not work
func (suiteSJ *TestSuiteForGitWorkSJ) Test_GitClone_WithBadUrl() {
	suiteSJ.singleJob.GitClone(badUrlProvided)
	err := Ls(PathToDirectoryForTest + suiteSJ.singleJob.commit.RepoName)
	if err == nil {
		suiteSJ.FailNow("Git clone was successful but it should not be")
	}
}

func (suiteSJ *TestSuiteForGitWorkSJ) Test_GitCheckout_ExistingBranch() {
	// setUp
	specialFolder := "TestGitCheckoutExistingBranch"
	os.Mkdir(filepath.Join(GlobalPath, TestResultsDir, specialFolder), os.ModePerm)

	// setting location for test files
	suiteSJ.singleJob.location = filepath.Join(PathToDirectoryForTest, specialFolder)
	suiteSJ.singleJob.locationOfRepository = filepath.Join(suiteSJ.singleJob.location, suiteSJ.singleJob.commit.RepoName)

	suiteSJ.singleJob.GitClone(urlProvided)

	// test
	suiteSJ.singleJob.GitCheckoutBranch()

	err := Ls(filepath.Join(suiteSJ.singleJob.locationOfRepository, nameOfTestDir))
	suiteSJ.Equal(nil, err, "Git checkout test was unsuccessful")
	// cleanUp
	suiteSJ.singleJob.RemoveRepository()
}

func (suiteSJ *TestSuiteForGitWorkSJ) Test_GitCheckout_NonExistentBranch() {
	// setUp
	suiteSJ.singleJob.location = filepath.Join(PathToDirectoryForTest, suiteSJ.singleJob.commit.RepoName)
	suiteSJ.singleJob.locationOfRepository = filepath.Join(suiteSJ.singleJob.location, suiteSJ.singleJob.commit.RepoName)

	suiteSJ.singleJob.GitClone(urlProvided)

	// test
	suiteSJ.singleJob.GitCheckoutBranch()
	err := Ls(filepath.Join(suiteSJ.singleJob.locationOfRepository, nameOfTestDir)) //Ls(location + "/" + nameOfTestDir)
	if err == nil {
		suiteSJ.FailNow("Git checkout test was unsuccessful: %s", err)
	}
	// cleanUp
	suiteSJ.singleJob.RemoveRepository()
}

func (suiteSJ *TestSuiteForGitWorkSJ) Test_GitSetStatusForCommit_ExistingCommit() {
	// SetUp
	resultLog := strings.Replace(TestResultsDir, `\`, "/", -1)
	var state string
	var objs []map[string]*json.RawMessage
	authorizationHeader := fmt.Sprintf("authorization: Basic %s", base64.StdEncoding.EncodeToString([]byte(BasicAuthToken)))

	// Test
	pendingStatusCheck := NewPendingStatus(resultLog)
	suiteSJ.singleJob.GitSetStatusForCommit(*pendingStatusCheck)

	// didn't work with http.DefaultClient.Do(req)
	bodyJSON, err := exec.Command("curl", "-H", authorizationHeader, "<company git link>/api/v3/repos/I516366/testRepositoryForGovis/commits/3190638c9a3b278c2f100ab16e4fb09aa83093bd/statuses").Output()
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal([]byte(bodyJSON), &objs); err != nil {
		panic(err)
	}
	state = fmt.Sprintf("%s", objs[0]["state"])
	suiteSJ.Equal(pendingStatusCheck.State, state[2:len(state)-1])

	failedStatusCheck := NewFailedStatus(resultLog)
	suiteSJ.singleJob.GitSetStatusForCommit(*failedStatusCheck)
	bodyJSON, err = exec.Command("curl", "-H", authorizationHeader, "<company git link>/api/v3/repos/I516366/testRepositoryForGovis/commits/3190638c9a3b278c2f100ab16e4fb09aa83093bd/statuses").Output()
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal([]byte(bodyJSON), &objs); err != nil {
		log.Fatal(err)
	}
	state = fmt.Sprintf("%s", objs[0]["state"])
	suiteSJ.Equal(failedStatusCheck.State, state[2:len(state)-1])

	successStatusCheck := NewSuccessStatus(resultLog)
	suiteSJ.singleJob.GitSetStatusForCommit(*successStatusCheck)
	bodyJSON, err = exec.Command("curl", "-H", authorizationHeader, "<company git link>/api/v3/repos/I516366/testRepositoryForGovis/commits/3190638c9a3b278c2f100ab16e4fb09aa83093bd/statuses").Output()
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal([]byte(bodyJSON), &objs); err != nil {
		log.Fatal(err)
	}
	state = fmt.Sprintf("%s", objs[0]["state"])
	suiteSJ.Equal(successStatusCheck.State, state[2:len(state)-1])
}
