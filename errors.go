package main

import (
	"errors"
	"fmt"
)

var (
	WriteToFileError  = errors.New("Failed to write to the file")
	HTTPGetError      = errors.New("Failed to retrieve data")
	ReadResponseError = errors.New("Failed to read from a HTTP response")
)

func formatErr(errName error, jobName string, err error) error {
	return fmt.Errorf("Error:%e, Job Name:%s, Description: %e", HTTPGetError, jobName, err)
}
