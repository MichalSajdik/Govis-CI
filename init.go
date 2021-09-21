package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func init() {
	LoggerServerINFO("Start Init\n")
	loadEnvVarsFromCLI()
	configureServer()
	createResultsFolder()
	logConfigurationOfEnvVariables()
	LoggerServerINFO("End Init\n")
}

func loadEnvVarsFromCLI() {
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--api":
			i = i + 1
			Api = os.Args[i]
		case "--secret":
			i = i + 1
			GithubSecret = os.Args[i]
		case "--basic_auth":
			i = i + 1
			BasicAuthToken = os.Args[i]
		case "--hostname":
			i = i + 1
			HostnameEnv = os.Args[i]
		case "--script_rights":
			i = i + 1
			ScriptRightsEnv = os.Args[i]
		case "--port":
			i = i + 1
			PortEnv = os.Args[i]
		case "--secure_mode":
			i = i + 1
			EnvSecureMode = os.Args[i]
		case "--script_timeout":
			i = i + 1
			EnvScriptTimeout = os.Args[i]
		case "--help":
			printHelp()
		case "-h":
			printHelp()
		case "-v":
			VerboseLogging = true
		default:
			LoggerServerWARN("Unknown parameters in CLI: " + os.Args[i] + "\n")
		}

	}
}

func logConfigurationOfEnvVariables() {
	globalEnvMap := make(map[string]string)
	globalEnvMap[EnvNameApi] = Api
	globalEnvMap[EnvNameBasicAuthToken] = BasicAuthToken
	globalEnvMap[EnvNameGithubSecret] = GithubSecret
	for globalVariableName, globalVariableValue := range globalEnvMap {
		if "" == globalVariableValue {
			partialMessage := globalVariableName + " is not set as ENVIRONMENTAL VARIABLE!\n"
			LoggerServerINFO(partialMessage)
		}
	}
}

func configureServer() {
	///////////////////////////////////
	// Setting `router` and `Server` //
	//////////////////////////////////
	router := mux.NewRouter().StrictSlash(true)

	// html interface || govis-web
	router.HandleFunc("/", Web.IndexHandler)
	router.HandleFunc("/"+ServerLogFileName, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(GlobalPath, ServerLogFileName))
	})
	router.HandleFunc("/"+AccessLogFileName, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(GlobalPath, AccessLogFileName))
	})
	//router.HandleFunc("/", Web.IndexHandler)
	// handling POST messages from github
	router.HandleFunc(PayloadUrl, EventsProcessing).Methods("POST")

	// handling access to data(FTP) via web browser, everything in `govis-web` directory is accessible
	router.
		PathPrefix("/" + ResultsDir).
		Handler(http.StripPrefix("/"+ResultsDir, http.FileServer(http.Dir(filepath.Join(GlobalPath, ResultsDir)))))

	accessLog := OpenLogFileWithName(GlobalPath, AccessLogFileName)
	loggedRouter := handlers.CombinedLoggingHandler(&accessLog, router)
	Server = http.Server{
		Addr:        GetGovisAddress(),
		Handler:     loggedRouter,
		ReadTimeout: time.Duration(30) * time.Second,
	}
}

func createResultsFolder() {
	// creating "results" directory
	err := os.Mkdir(filepath.Join(GlobalPath, ResultsDir), os.ModePerm)
	if err != nil {
		LoggerServerINFO("Unable to create directory `results` >> " + err.Error())
	}
	// creating "testResults" directory
	err = os.Mkdir(filepath.Join(GlobalPath, TestResultsDir), os.ModePerm)
	if err != nil {
		LoggerServerINFO("Unable to create directory `restResults` >> " + err.Error())
	}
}

func printHelp() {
	help := `Example: ./Govis-CI
	[--api <api>]
	[--secret <github secret>]
	[--basic_auth <github personal token>]
	[--hostname <hostname/ip>]
	[--status_baseurl <hostname:port/ip:port>]
	[--script_rights <rights>]
	[--secure_mode <mode>]
	[--script_timeout <seconds>]
	[--port <port>]
	[-v]
	
	<api> == rest api for organization, default: https://api.github.com/
	<github secret>] == WebHook Secret
	<github personal token> == Personal_access_Api_token of user with set status access to the repository
	<hostname/ip> == hostname on which we accept webhooks (0.0.0.0 accepts all and is default)
	<hostname:port/ip:port> == hostname or ip with port number indicating where will user be redirected for results
	<rights> == if unset Govis can't run 'script' from yaml directly
	<mode> == if set to 'negative' then we don't use token for verification with github
	<port> == port to start Govis on - default : '8000'
	<seconds> == number representing time in seconds
	'-v' stands for verbose logging`
	LoggerServerINFO("Help was printed")
	fmt.Println(help)
	os.Exit(0)
}
