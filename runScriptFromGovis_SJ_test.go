package main

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type TestSuiteForScriptSJ struct {
	suite.Suite
	TestSuiteBase
}

func (suiteSJ *TestSuiteForScriptSJ) SetupSuite() {
	commit := Commit{"Not_Used", "", "Not_Used", "Not_Used"}
	logFile := OpenLogFile(PathToDirectoryForTest)
	suiteSJ.singleJob = SingleJobCreate("Not_Used", PathToDirectoryForTest, logFile, commit)
}

func TestScriptSJTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuiteForScriptSJ))
}

func (suiteSJ *TestSuiteForScriptSJ) Test_RunScript_WithGoodScript() {
	if PlatformIsWindows() {
		return
	}
	script := "ls"
	err := suiteSJ.singleJob.RunScript(script, make([]string, 0))
	suiteSJ.Equal(nil, err)
}

func (suiteSJ *TestSuiteForScriptSJ) Test_RunScript_WithBadScript() {
	script := "asdasdasdasdasdd"
	err := suiteSJ.singleJob.RunScript(script, make([]string, 0))
	if err == nil {
		suiteSJ.FailNow("Error should be returned", "Error was not returned")
	}
}

func Test_ParseGovisYaml_AbapMerge_WithValidValues(t *testing.T) {
	logFile := OpenLogFile(PathToDirectoryForTest)
	yamlOfScripts := `abapMerge: 
- ls
- date`
	fmt.Println(yamlOfScripts)
	testFile, err := os.Create(filepath.Join(PathToDirectoryForTest, GovisYAML))
	if err != nil {
		t.Error("internal issue >> could not create testFile >> " + err.Error())
	}
	_, err = testFile.WriteString(yamlOfScripts)
	if err != nil {
		t.Error("internal issue >> could not write to testFile >> " + err.Error())
	}
	err = testFile.Close()
	if err != nil {
		t.Error("File " + GovisYAML + " was not closed >> " + err.Error())
	}
	singleJob := SingleJobCreate("space", PathToDirectoryForTest, logFile, Commit{
		Owner:      "space",
		RepoName:   "",
		IdOfCommit: "space",
	})
	// test
	scriptArray, err := singleJob.ParseGovis()
	if err != nil {
		t.Error("Unable to parse govis >> " + err.Error())
	}
	if len(scriptArray.AbapMerge) != 2 {
		t.Error(len(scriptArray.AbapMerge), ": ", scriptArray)
	}
	// cleanup
	err = os.Remove(filepath.Join(PathToDirectoryForTest, GovisYAML))
	if err != nil {
		t.Error("File " + GovisYAML + " was not deleted >> " + err.Error())
	}
}
