package main

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type TestSuiteBase struct {
	logFile   os.File
	singleJob *SingleJob
}

type TestSuiteForSingleJob struct {
	suite.Suite
	TestSuiteBase
}

func (suiteSJ *TestSuiteForSingleJob) SetupSuite() {
	commit := Commit{"Not_Used", "Not_Used", "Not_Used", "Not_Used"}
	logFile := OpenLogFile(PathToDirectoryForTest)
	suiteSJ.singleJob = SingleJobCreate("Not_Used", PathToDirectoryForTest, logFile, commit)
}

func TestSingleJobTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuiteForSingleJob))
}

func (suiteSJ *TestSuiteForSingleJob) Test_WorkWithRepository_WithBadScript() {
	var err error
	err = suiteSJ.singleJob.WorkWithRepository()
	if err == nil {
		suiteSJ.FailNow("No LogError without govis yaml")
	}
}
