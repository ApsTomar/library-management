package efk

import (
	"encoding/json"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/library/envConfig"
	"github.com/library/models"
	"strconv"
	"time"
)

func NewLogger(env *envConfig.Env) *fluent.Fluent {
	fluentPort, err := strconv.Atoi(env.FluentPort)
	if err != nil {
		glog.Fatal(err)
	}
	logger, err := fluent.New(fluent.Config{
		FluentPort: fluentPort,
		FluentHost: env.FluentHost,
	})
	if err != nil {
		glog.Fatal(err)
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
		glog.Error(err)
	}
	if err := logger.Post(tag, m); err != nil {
		glog.Error(err)
	}
}
