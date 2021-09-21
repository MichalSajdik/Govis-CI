package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type GovisWeb struct {
	Name                string
	Time                string
	Api                 string
	GithubSecretColor   string
	BasicAuthTokenColor string
	HostnameValue       string
	GovisAddressValue   string
	GitIsInstalledColor string
	GovisCanWriteColor  string
	GovisSeesAPIColor   string
	GovisVersionValue   string
}

func (web *GovisWeb) IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate, err := template.New("index.html").Parse(indexTempString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	path := r.URL.Path
	LoggerServerINFO("Request for >> " + path)
	if path == "" {
		path = "index.html"
	} else {
		path = RemovePrefix("/", path)
		if path == "" {
			path = "index.html"
		}
	}
	LoggerServerINFO("FileName of web is >> " + path)

	web.checkConfigurationOfEnvironment()

	web.updateTime()
	if err := indexTemplate.ExecuteTemplate(w, path, web); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (web *GovisWeb) updateTime() {
	web.Time = time.Now().Format(time.Stamp)
}

func getVersion() string {
	if GovisVersionCommitId == ""{
		return "GOVIS_1.0"
	}
	return "GOVIS_1.0" + "_commitID_" + GovisVersionCommitId
}

func (web *GovisWeb) checkConfigurationOfEnvironment() (int, string) {
	web.GovisVersionValue = getVersion()
	web.Api = Api
	returnStatusCode := 200
	message := "GovisReport:\n"

	resp, err := http.Get(Api)
	if err != nil {
		errMessage := "Can't connect to Api >> " + err.Error() + "\n"
		LoggerServerWARN(errMessage)
		web.GovisSeesAPIColor = "#FF0000"
		returnStatusCode = 500
		message = message + errMessage
	} else {
		if resp.StatusCode != 200 {
			errMessage := "Can't connect to Api\n"
			LoggerServerWARN(errMessage)
			web.GovisSeesAPIColor = "#FF0000"
			returnStatusCode = 500
			message = message + errMessage
		} else {
			web.GovisSeesAPIColor = "#00FF00"
		}
	}

	if BasicAuthToken == "" {
		errMessage := "BasicAuthToken is not set\n"
		LoggerServerWARN(errMessage)
		web.BasicAuthTokenColor = "#FF0000"
		returnStatusCode = 500
		message = message + errMessage
	} else {
		web.BasicAuthTokenColor = "#00FF00"
	}

	if GithubSecret == "" {
		errMessage := "GithubSecret is not set\n"
		LoggerServerWARN(errMessage)
		web.GithubSecretColor = "#FF0000"
		returnStatusCode = 500
		message = message + errMessage
	} else {
		web.GithubSecretColor = "#00FF00"
	}

	web.HostnameValue = GetHostname()
	web.GovisAddressValue = GetGovisAddress()

	if !GovisCanWrite() {
		errMessage := "Govis can't write files to results folder\n"
		LoggerServerWARN(errMessage)
		web.GovisCanWriteColor = "#FF0000"
		returnStatusCode = 500
		message = message + errMessage
	} else {
		web.GovisCanWriteColor = "#00FF00"
	}

	if !GitIsInstalled() {
		errMessage := "Git is not installed\n"
		LoggerServerWARN(errMessage)
		web.GitIsInstalledColor = "#FF0000"
		returnStatusCode = 500
		message = message + errMessage
	} else {
		web.GitIsInstalledColor = "#00FF00"
	}

	return returnStatusCode, message
}

func GitIsInstalled() bool {
	cmd := exec.Command("git", "version")
	out, err := cmd.Output()
	if err != nil {
		LoggerServerWARN(fmt.Sprintf("Git does not exist >> %s", err.Error()))
		return false
	} else {
		LoggerServerINFO(fmt.Sprintf("Git does exist >> %s", out))
		return true
	}

}

func GovisCanWrite() bool {

	testData := []byte("test data for: >> \n\t\t>> govis\n")
	filePath := filepath.Join(GlobalPath, ResultsDir, "test.data")
	err := ioutil.WriteFile(filePath, testData, os.ModePerm)
	if err != nil {
		LoggerServerWARN("Unable to create file `test.data` we dont have right's to create files >> " + err.Error())
		return false
	}
	os.Remove(filePath)
	return true
}

const indexTempString string = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
		<title>Govis-Web</title>
		<style>
		.welcome { color: darkblue; }
		</style>
</head>
<body>
Version : {{.GovisVersionValue}}
<br>
<div class="welcome"><b>Welcome {{.Name}}</b> , Load time is {{.Time}}</div>
<span style="color: {{.GovisSeesAPIColor}}; ">GovisSeesAPI</span>
|
<span style="color: {{.BasicAuthTokenColor}}; ">BasicAuthToken</span>
|
<span style="color: {{.GovisCanWriteColor}}; ">GovisCanWrite</span>
|
<span style="color: {{.GithubSecretColor}}; ">GithubSecret</span>
|
<span style="color: {{.GitIsInstalledColor}}; ">GitIsInstalled</span>
<br>
<span style="color: Black; ">Hostname: {{.HostnameValue}}</span>
<br>
<span style="color: Black; ">Address: {{.GovisAddressValue}}</span>
<br>
<span style="color: Black; ">API: {{.Api}}</span>

<br>
<br>
<a href="results/">Results folder</a><br>
<a href="AccessLog.log">AccessLog.log</a><br>
<a href="ServerLog.log">ServerLog.log</a><br>
` + `
<h2>Readme >></h2><br>
# Setup:<br>
<br>
## Server:<br>
<br>
1. Create user, and folder he owns<br>
2. Move Govis binary from releases to that folder<br>
3. Configuration<br>
4. Set Environment Variables:<br>
<br>
- "GOVIS_SECRET" - WebHook Secret<br>
- "GOVIS_BASIC_AUTH" - Personal_access_Api_token of user with set status access to the repository<br>
- "GOVIS_PORT" - port to start Govis on - default : "8000"<br>
- "GOVIS_SCRIPT_RIGHTS" - if unset Govis can't run "script" from yaml directly, only from "abapMerge:" inside podman<br>
- "GOVIS_SECURE_MODE" - if set to "negative" then we don't use token for verification with github<br>
- "GOVIS_STATUS_BASEURL" - url with port for accessing govis-ci server<br>
- "GOVIS_SCRIPT_TIMEOUT" - number of seconds to timeout each script run - default : 300 <br>
- "GOVIS_HOSTNAME" - hostname on which we accept webhooks (0.0.0.0 accepts all and is default)<br>
- "GOVIS_LOG" - if set to any string then we use verbose logging otherwise we don't use it<br>
- "GOVIS_API" - rest api of <br>
    - SAP : <compnay git url>/api/v3<br>
    - Normal GitHub (default) : https://api.github.com/<br>
    <br>
5. Run binary in service/nohup/screen/tmux<br>
<br>
## Github:<br>
<br>
1. Create new webhook in github repository that targets "IP:PORT/payload/" , and has "GOVIS_SECRET" as secret, and runs on pull request event.<br>
2. Create ".govis.yaml" file with following structure:<br>
<br>
"<br>
script:<br>
- test command<br>
- another test command<br>
"<br>
3. Create a PR with this file.<br>
4. After result status gets back go to Settings - Branches, and add that status to required<br>
<br>
- you can view results at "http:/IP:PORT/govis-ci/results/", but only when program is running<br>

</body>
</html>`
