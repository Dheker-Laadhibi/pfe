/*

	Package config provides configuration settings and initialization functions related to logging.

	Functions:
	- InitLog(): Configures the logging settings, including log level and formatting.

	Dependencies:
	- "labs/utils": Custom package for utility functions.
	- "github.com/antonfisher/nested-logrus-formatter": Logrus formatter for nested log entries.
	- "github.com/sirupsen/logrus": Structured logger for Go.

	Environment Variables:
	- LOG_LEVEL: Log level setting for controlling the verbosity of log output.

	Usage:
	- Call InitLog() at the beginning of the program to configure the logging settings.

	Note:
	- The function uses the nested-logrus-formatter for enhanced log formatting.
	- The log level is determined by the LOG_LEVEL environment variable, defaulting to INFO if not specified.

	Last update :
	01/02/2024 10:22

*/

package config

import (
	"labs/utils"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

// InitLog configures the logging settings, including log level and formatting.
func InitLog() {

	logLevel, err := utils.GetStringEnv("LOG_LEVEL")
	if err != nil {
		logrus.Fatal("Failed to load LOG_LEVEL from env file: ", err)
	}

	logrus.SetLevel(getLoggerLevel(logLevel))
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&nested.Formatter{
		HideKeys:        false,
		FieldsOrder:     []string{"component", "category"},
		TimestampFormat: "2006-01-02 15:04:05",
		ShowFullLevel:   false,
		CallerFirst:     false,
	})
}

// getLoggerLevel returns the logrus.Level based on the provided log level string.
func getLoggerLevel(value string) logrus.Level {
	switch value {
	case "DEBUG":
		return logrus.DebugLevel
	case "TRACE":
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}
