package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const LogInfoPrefix = "INFO: "
const LogWarnPrefix = "WARN: "
const LogErrorPrefix = "ERR : "

type LoggingComponent struct {
	logFile os.File
}

func (lc *LoggingComponent) LogInfo(message string) {
	if !VerboseLogging {
		return
	}
	logger(lc.logFile, message, LogInfoPrefix)
}

func (lc *LoggingComponent) LogWarn(message string) {
	logger(lc.logFile, message, LogWarnPrefix)
	logger(*os.Stderr, message, LogWarnPrefix)
}

func (lc *LoggingComponent) LogError(message string) {
	logger(lc.logFile, message, LogErrorPrefix)
	logger(*os.Stderr, message, LogErrorPrefix)
	os.Exit(1)
}

// LoggerINFO is used for parallel logging to more files
// @file destination file for log
// @message string to be written
// todo add lock so other process is not able to switch output
func logger(file os.File, message, prefix string) {
	log.SetOutput(&file)
	log.SetPrefix(prefix)
	log.Println(message)
}

func loggerServer(message, prefix string) {
	logger(ServerLogFile, message, prefix)
}

func LoggerServerINFO(message string) {
	if !VerboseLogging {
		return
	}
	loggerServer(message, LogInfoPrefix)
}

func LoggerServerWARN(message string) {
	loggerServer(message, LogWarnPrefix)
	logger(*os.Stderr, message, LogWarnPrefix)
}

func LoggerServerERROR(message string) {
	loggerServer(message, LogErrorPrefix)
	logger(*os.Stderr, message, LogErrorPrefix)
	os.Exit(1)
}

func OpenLogFileWithName(location, name string) os.File {
	// preparing file for logging
	locationOfLogFile := filepath.Join(location, name)
	logFile, err := os.OpenFile(locationOfLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	fmt.Println("LogFile at >> ", locationOfLogFile, "Printed from OpenLogFile() in helper.go")
	return *logFile
}

// OpenLogFile opens/creates @logFileName file for logging
// return os.file
func OpenLogFile(location string) os.File {
	return OpenLogFileWithName(location, LogFileName)
}
