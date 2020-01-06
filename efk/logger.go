package efk

import (
	"encoding/json"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/google/uuid"
	"github.com/library/envConfig"
	"github.com/library/models"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func NewLogger(env *envConfig.Env) *fluent.Fluent {
	fluentPort, err := strconv.Atoi(env.FluentPort)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
		}).Error("error in type conversion of fluent port")
	}
	logger, err := fluent.New(fluent.Config{
		FluentPort: fluentPort,
		FluentHost: env.FluentHost,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
		}).Error("error creating new EFK logger")
	}
	return logger
}

func LogError(logger *fluent.Fluent, tag, task string, err error, statusCode int) {
	loggerErr := models.EfkLogger{
		ID:         uuid.New().String(),
		Timestamp:  time.Now(),
		Task:       task,
		Error:      err.Error(),
		StatusCode: statusCode,
	}
	marshalledLog, _ := json.Marshal(loggerErr)
	var m map[string]interface{}
	err = json.Unmarshal(marshalledLog, &m)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
		}).Error("error in marshalling EFK log")
	}
	if err := logger.Post(tag, m); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
		}).Error("error posting EFK log")
	}
}
