package main

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type TestSuiteForParseGovis struct {
	suite.Suite
	TestSuiteBase
	goodYamlScript string
	badYamlScript  string
	testFile       *os.File
}

func (suiteSJ *TestSuiteForParseGovis) SetupSuite() {
	suiteSJ.badYamlScript = "{ \"script\":[\"go run main.go\",\"go run main2.go\",\"go run main3.go\", \"go run main4.go\"}"
	suiteSJ.goodYamlScript = `script: 
- go run main.go
- go run main2.go
- go run main3.go
- go run main4.go`

	commit := Commit{"Not_Used", "", "Not_Used", "Not_Used"}
	logFile := OpenLogFile(PathToDirectoryForTest)
	suiteSJ.singleJob = SingleJobCreate("Not_Used", PathToDirectoryForTest, logFile, commit)
}

func (suiteSJ *TestSuiteForParseGovis) SetupTest() {
	testFile, err := os.Create(filepath.Join(PathToDirectoryForTest, GovisYAML))
	if err != nil {
		suiteSJ.FailNow("internal issue >> could not create testFile >> " + err.Error())
	}
	suiteSJ.testFile = testFile
}

func (suiteSJ *TestSuiteForParseGovis) TearDownTest() {
	err := suiteSJ.testFile.Close()
	if err != nil {
		suiteSJ.FailNow("File " + GovisYAML + " was not closed >> " + err.Error())
	}

	err = os.Remove(filepath.Join(PathToDirectoryForTest, GovisYAML))
	if err != nil {
		suiteSJ.FailNow("File " + GovisYAML + " was not deleted >> " + err.Error())
	}
}

func TestParseGovisSJTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuiteForParseGovis))
}

func (suiteSJ *TestSuiteForParseGovis) Test_ParseGovisYaml_WithValidValues() {

	//suiteSJ.singleJob.locationOfRepository = PathToDirectoryForTest

	_, err := suiteSJ.testFile.WriteString(suiteSJ.goodYamlScript)
	if err != nil {
		suiteSJ.FailNow("internal issue >> could not write to testFile >> " + err.Error())
	}

	// test
	scriptArray, err := suiteSJ.singleJob.ParseGovis()
	suiteSJ.Equal(nil, err, "Unable to parse govis")
	suiteSJ.Equal(4, len(scriptArray.Script), fmt.Sprintf("%v%v%v", len(scriptArray.Script), ": ", scriptArray))

}

func (suiteSJ *TestSuiteForParseGovis) Test_ParseGovis_WithInvalidValues() {
	_, err := suiteSJ.testFile.WriteString(suiteSJ.badYamlScript)
	if err != nil {
		suiteSJ.FailNow("internal issue >> could not write to testFile >> " + err.Error())
	}

	// test
	scriptArray, err := suiteSJ.singleJob.ParseGovis()
	if err == nil {
		suiteSJ.FailNow(fmt.Sprintf("%v%v%v", len(scriptArray.Script), ": ", scriptArray))
	}

}
