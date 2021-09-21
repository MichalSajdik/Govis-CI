package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

//ParseGovis parses .govis.yaml
// @return array of strings representing scripts in govis file
func (singleJob *SingleJob) ParseGovis() (GovisScriptStructYaml, error) {

	govisScriptInArray, err := parseGovisYaml(singleJob)
	if err != nil {
		singleJob.LogWarn("LogError during parsing GOVIS " + GovisYAML)
	} else {
		singleJob.LogInfo("Govis parsed from " + GovisYAML)
	}
	singleJob.LogInfo(fmt.Sprintf("GovisParse >>  %v", govisScriptInArray))
	return govisScriptInArray, err
}

//parseGovisYamlBash parses govis file in @repoName directory of .yaml format
// @return array of strings representing scripts in govis file
func parseGovisYaml(singleJob *SingleJob) (GovisScriptStructYaml, error) {
	locationOfRepository := singleJob.locationOfRepository
	var govisScript []byte
	var govisScriptInArray GovisScriptStructYaml
	var err error
	if locationOfRepository == "" {
		govisScript, err = ioutil.ReadFile(GovisYAML)
		if err != nil {
			singleJob.LogWarn(fmt.Sprintf("LogError during reading %s: %s", GovisYAML, err))
			return govisScriptInArray, err
		}
	} else {
		govisScript, err = ioutil.ReadFile(filepath.Join(locationOfRepository, GovisYAML))
		if err != nil {
			singleJob.LogWarn(fmt.Sprintf("LogError during reading %s: %s", GovisYAML, err))
			return govisScriptInArray, err
		}
	}

	err = yaml.Unmarshal(govisScript, &govisScriptInArray)
	if err != nil {
		singleJob.LogWarn(fmt.Sprintf("yaml is >> %s", govisScript))
		singleJob.LogWarn(fmt.Sprintf("LogError while parsing yaml >> %s", err))
	}

	return govisScriptInArray, err
}
