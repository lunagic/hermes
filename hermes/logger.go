package hermes

import (
	"log"

	"github.com/fatih/color"
)

// Logger for the howler engine.
type Logger interface {
	Log(message string)
	Error(message string)
	Info(message string)
	Success(message string)
}

// loggerLogger is a instance of logger that logs using a logger.
type loggerLogger struct {
	Logger *log.Logger
}

// Log a message.
func (logger loggerLogger) Log(message string) {
	bold := color.New(color.Bold)
	logger.log(bold.Sprint(message))
}

// Info level log.
func (logger loggerLogger) Info(message string) {
	blue := color.New(color.FgCyan)
	logger.log(blue.Sprint(message))
}

// Error level log.
func (logger loggerLogger) Error(message string) {
	red := color.New(color.FgRed)
	logger.log(red.Sprint(message))
}

// Success level log.
func (logger loggerLogger) Success(message string) {
	green := color.New(color.FgGreen)
	logger.log(green.Sprint(message))
}

func (logger loggerLogger) log(message string) {
	logger.Logger.Print(message)
}
