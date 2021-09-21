package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// GetPathToPowerShell is used to find powershell.exe and return path in string + checks errors
func GetPathToPowerShell() string {
	// get path to powershell.exe
	pathToPowerShell, err := exec.LookPath("powershell.exe") //%SystemRoot%\system32\WindowsPowerShell\v1.0\powershell.exe
	if err != nil {
		LoggerServerINFO(fmt.Sprintf("Unable to find powershell.exe >> %s\n", err.Error()))
	}
	// print powershell.exe location
	scriptForPowerShell := "echo " + "'PowerShell location >> " + pathToPowerShell + "'"
	out, err := exec.Command(pathToPowerShell, scriptForPowerShell).Output()
	if err != nil {
		LoggerServerINFO(fmt.Sprintf("%s\n", err))
	} else {
		LoggerServerINFO(fmt.Sprintf("%s\n", out))
	}
	return pathToPowerShell
}

// GetAuthorizedURL takes @htmlURL and adds @token for verification
func GetAuthorizedURL(htmlURL string, token string) (string, error) {
	if GetSecureMode() == false {
		return htmlURL, nil
	}
	if len(htmlURL) < 1 {
		LoggerServerWARN(fmt.Sprintf(htmlURL + " >> html url too short"))
		return htmlURL, errors.New("htmlurl too short")
	}
	removeHTTPS := strings.Replace(htmlURL, "https://", "https://"+token+"@", 1)
	return removeHTTPS, nil
}

// GetRepoName takes @htmlUrl url to repository and extracts @repoName repository name
func GetRepoName(htmlUrl string) (string, error) {
	htmlUrlSplit := strings.Split(htmlUrl, "/")
	repoName := strings.TrimSuffix(htmlUrlSplit[len(htmlUrlSplit)-1], ".git")
	LoggerServerINFO(fmt.Sprintf("Repository name >> %s\n", repoName))
	if len(repoName) < 1 {
		LoggerServerINFO(fmt.Sprintf(repoName + " >> repository name too short"))
		return repoName, errors.New("short name for repoName")
	}
	if len(repoName) == 1 && (repoName[0] == '/' || repoName[0] == '\\') {
		LoggerServerINFO(fmt.Sprintf(repoName + " >> is empty"))
		return repoName, errors.New("empty name for repoName")
	}
	return repoName, nil
}

// PrintReadmeFile used only for debug
// it will print README.md file in given repository name from @htmlUrl
// @deprecated
func PrintReadmeFile(htmlUrl string) {
	repoName, _ := GetRepoName(htmlUrl)
	scriptForPowerShell := "Set-Location " + repoName + " ; cat README.md"
	govisScript, err := exec.Command(GetPathToPowerShell(), scriptForPowerShell).Output()
	if err != nil {
		LoggerServerINFO(fmt.Sprintf("Unable to find %s >> %s\n", "README.md", err.Error()))
	} else {
		LoggerServerINFO(fmt.Sprintf("Parsing %s  >> %s\n", "README.md", govisScript))
	}
	LoggerServerINFO(fmt.Sprintf("Parsing %s  >> %s\n", "script", scriptForPowerShell))
}

// Ls is just debug/testing tool
func Ls(location string) error {
	var err error
	if PlatformIsWindows() {
		err = LsPS(location)
	} else {
		err = LsBash(location)
	}
	return err
}

// LsBash is just debug/testing tool
func LsBash(location string) error {

	cmd := exec.Command("ls", location)

	LoggerServerINFO(fmt.Sprintf("dir >>>>>>>>>>> %s\n", GlobalPath))
	LoggerServerINFO(fmt.Sprintf("cmd.Dir >>>>>>>>>>> %s\n", cmd.Dir))
	cmd.Dir = GlobalPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		LoggerServerINFO(fmt.Sprintf("cmd.Run() failed with %s\n", err.Error()))
	}
	LoggerServerINFO(fmt.Sprintf("combined out:\n%s\n", string(out)))
	return err
}

// LsPS is just debug/testing tool
func LsPS(location string) error {
	script := "ls " + location
	out, err := exec.Command(GetPathToPowerShell(), script).Output()
	if err != nil {
		LoggerServerINFO(fmt.Sprintf("LogError with script %s >> %s\n", script, err.Error()))
	} else {
		LoggerServerINFO(fmt.Sprintf("Finished %s  >> %s\n", script, out))
	}
	return err
}

// PlatformIsWindows helper func to learn on what platform we run
func PlatformIsWindows() bool {
	platform := runtime.GOOS // windows environmental variable
	LoggerServerINFO(fmt.Sprintf("platform>>%s%s", platform, "<<"))
	return strings.Contains(strings.ToLower(platform), "windows")
}

// Returns hostname:port, env vars host> hostname command: env var port > default port
func GetGovisAddress() string {
	govisPort := func() string {
		if len(PortEnv) > 1 {
			return PortEnv
		} else {
			return Port
		}
	}()
	govisHostname := GetHostname()
	return govisHostname + ":" + govisPort
}

func GetHostname() string {
	if len(HostnameEnv) > 1 {
		return HostnameEnv
	} else {
		hostnameOS, _ := os.Hostname()
		return string(hostnameOS)
	}
}

func RemovePrefix(prefix, originalString string) string {
	n := len(prefix)
	if len(originalString) <= len(prefix) {
		return ""
	}
	newString := originalString[n:]
	return newString
}
