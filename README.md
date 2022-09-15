# gormv2-logrus

Easily connect Gorm V2 and logrus with customizable options

Usage example :

```go
package test

import (
	gormv2logrus "github.com/thomas-tacquet/gormv2-logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func myFunctionToConnectOnMyDB(e *logrus.Entry) {
	gormLogger := gormv2logrus.NewGormlog(gormv2logrus.WithLogrusEntry(e))
	gormLogger.LogMode(logger.Error)

	gormConfig := &gorm.Config{
		Logger:                 gormLogger,
		CreateBatchSize:        1500,
		SkipDefaultTransaction: true,
	}

	db, err := gorm.Open(
		"CONNEXION STRING",
		gormConfig,
	)
}
```

Or you can use with a logrus.Logger :

```go
gormLogger := gormv2logrus.NewGormlog(gormv2logrus.WithLogrus(e))
```

You can also add customization options :

```go
	slowThresholdDuration, _ := time.ParseDuration("300ms")

	gormLogger := gormv2logrus.NewGormlog(
		gormv2logrus.WithLogrus(e),
		gormv2logrus.WithGormOptions(
			gormv2logrus.GormOptions{
				SlowThreshold: slowThresholdDuration,
				LogLevel:      logger.Error,
				LogLatency:    true,
			},
		),
	)
```

> Colorful option has been removed, because Logrus can't really handle the color code, especially if the logs are redirected to a file, the characters are not removed automatically

## Contibuting

Just feel free to open issues, ask questions, make proposals.
