package utils

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func NewLogger() {
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "./lark-chatGPT"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "content",
			"levelKey": "type",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		fmt.Println("logger err: ", err.Error())
		panic(err)
	}
	_logger, err := cfg.Build()
	logger = _logger
	if err != nil {
		fmt.Println("logger err: ", err.Error())
		panic(err)
	}
	defer _logger.Sync()
}

func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg)
}
