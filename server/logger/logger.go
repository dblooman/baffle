package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Fields is for passing semi-structured data that doesn't already have a type to the logger
type Fields map[string]interface{}

// Fielder is a interface that allows any type to be converted into log fields
type Fielder interface {
	ToFields() Fields
}

var (
	log = logrus.New()
)

func init() {
	log.Formatter = new(logrus.JSONFormatter)
	log.Out = os.Stderr

}

// Info logs with the given message & addition fields at the INFO Level
// This doesn't log to sentry
func Info(message string, data Fields) {
	log.WithFields(logrus.Fields(data)).Info(message)
}

// Error logs an error with the error message & converts the message to fields if Fielder is implemented
func Error(err error) {
	log.WithField("error", errorToField(err)).Error(err.Error())
}

// Fatal logs errors in the same way as Error then flushes the errors and calls os.Exit(1)
func Fatal(err error) {
	// this doesn't use log.Fatal because we need to flush the hook before exiting
	log.WithField("error", errorToField(err)).Error(err.Error())
	os.Exit(1)
}

func errorToField(err error) interface{} {
	if err == nil {
		return nil
	}

	if fielder, ok := err.(Fielder); ok {
		return fielder.ToFields()
	}

	return err.Error()
}