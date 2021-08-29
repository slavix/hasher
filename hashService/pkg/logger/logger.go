package logger

import (
	"context"
	"fmt"
	formatter "github.com/fabienm/go-logrus-formatters"
	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"reflect"
)

var logger *logrus.Logger

func Init(serviceName string, level logrus.Level) *logrus.Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(level)

	gelfFmt := formatter.NewGelf(serviceName)
	logrusLogger.SetFormatter(gelfFmt)

	hook := graylog.NewGraylogHook("localhost:12201", map[string]interface{}{})
	logrusLogger.AddHook(hook)

	logger = logrusLogger

	return logrusLogger
}

func Panic(packageName string, funcName string, err error, msg string) {
	trace := errors.WithStack(err)
	logger.WithFields(logrus.Fields{
		"package":  packageName,
		"function": funcName,
		"error":    err,
		"trace":    trace,
	}).Panic(msg)
}

func Error(ctx context.Context, packageName string, funcName string, err error, msg string) {
	trace := errors.WithStack(err)
	logger.WithFields(logrus.Fields{
		"requestID": GetRequestID(ctx),
		"package":   packageName,
		"function":  funcName,
		"error":     err,
		"trace":     trace,
	}).Error(msg)
}

func Warn(ctx context.Context, packageName string, funcName string, err error, msg string) {
	trace := errors.WithStack(err)
	logger.WithFields(logrus.Fields{
		"requestID": GetRequestID(ctx),
		"package":   packageName,
		"function":  funcName,
		"error":     err,
		"trace":     trace,
	}).Warn(msg)
}

func Info(ctx context.Context, packageName string, funcName string, msg string, additionalInfo string) {
	logger.WithFields(logrus.Fields{
		"requestID": GetRequestID(ctx),
		"package":   packageName,
		"function":  funcName,
		"info":      additionalInfo,
	}).Info(msg)
}

func Debug(ctx context.Context, packageName string, funcName string, msg string, additionalInfo string) {
	logger.WithFields(logrus.Fields{
		"requestID": GetRequestID(ctx),
		"package":   packageName,
		"function":  funcName,
		"debugInfo": additionalInfo,
	}).Debug(msg)
}

func Trace(ctx context.Context, packageName string, funcName string, msg string, additionalInfo string) {
	logger.WithFields(logrus.Fields{
		"requestID": GetRequestID(ctx),
		"package":   packageName,
		"function":  funcName,
		"traceInfo": additionalInfo,
	}).Trace(msg)
}

func LogHandler(message string, a ...interface{}) {
	logger.WithFields(logrus.Fields{
		"package": "tcp",
		"info":    fmt.Sprintf(message, a...),
	}).Info("server tcp log")
}

func GetRequestID(ctx context.Context) string {
	requestIDCtxValue := ctx.Value("requestID")

	if reflect.TypeOf(requestIDCtxValue).Name() == "string" {
		requestID, ok := requestIDCtxValue.(string)
		if ok {
			return requestID
		}
	}

	logger.WithFields(logrus.Fields{
		"package":  "logger",
		"function": "getRequestID",
	}).Error("request doesn't have requestID")

	return ""
}
