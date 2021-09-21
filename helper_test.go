package main

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

const locationOfPowerShell = "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"

//func TestSomething(t *testing.T) {
//
//	Ls("/")
//	fmt.Println()
//	t.LogError("I'm in a bad mood.")
//}

type TestSuiteForResources struct {
	suite.Suite
}

func (t *TestSuiteForResources) SetupSuite() {
}

func TestResourcesTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuiteForResources))
}

func (t *TestSuiteForResources) TestGetPathToPowerShell() {
	if PlatformIsWindows() {
		testString := GetPathToPowerShell()
		if locationOfPowerShell != testString {
			t.FailNow("Wrong location of powershell.exe")
		}
	}
}

func (t *TestSuiteForResources) TestGetAuthorizedURL() {
	html := "https://test/I516366/expectedResult"
	result, _ := GetAuthorizedURL(html, "randomtoken")

	expectedResult := "https://randomtoken@test/I516366/expectedResult"
	if expectedResult != result {
		t.FailNow(expectedResult + "!=" + result)
	}
}

func (t *TestSuiteForResources) TestGetRepoName() {
	htmlToExtractRepoNameFrom := "corp/I516366/expectedResult"
	result, _ := GetRepoName(htmlToExtractRepoNameFrom)

	expectedResult := "expectedResult"
	if expectedResult != result {
		t.FailNow("GetRepoName() extracted wrong string from " + htmlToExtractRepoNameFrom + " was extracted: " + result + " was expecting: " + expectedResult)
	}
}

func (t *TestSuiteForResources) TestGetRepoName2() {
	htmlToExtractRepoNameFrom := "corp/I516366/"
	_, err := GetRepoName(htmlToExtractRepoNameFrom)
	if err == nil {
		t.FailNow(htmlToExtractRepoNameFrom + "contains wrong format >> it should not end at '/' char")
	}
}

func (t *TestSuiteForResources) Test_LsRepo_FileExists() {
	err := Ls(GlobalPath)
	if err != nil {
		t.FailNow("No files in GlobalPath")
	}
}

func (t *TestSuiteForResources) Test_LsRepo_FileNotExist() {
	err := Ls("asdsdsdasd")
	if err == nil {
		t.FailNow("Ls function is not reliable")
	}

}

func (t *TestSuiteForResources) Test_GetGovisAddress() {
	HostnameEnv = "testurl"

	envAddress := GetGovisAddress()
	if envAddress != "testurl:8000" {
		t.FailNow(envAddress + "!=" + "testurl:8000")
	}

	PortEnv = "666"

	envAddress = GetGovisAddress()
	if envAddress != "testurl:666" {
		t.FailNow(envAddress + "!=" + "testurl:666")
	}
}

func TestPass(t *testing.T) {

	// pass
}
