package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/crazyfrankie/cronjob/pkg/model"
	"github.com/crazyfrankie/cronjob/worker/conf"
)

type LogManager struct {
	logCollection  *mongo.Collection
	logChan        chan *model.JobLog
	autoCommitChan chan *model.JobLogBatch
}

func NewLogManager(db *mongo.Database) *LogManager {
	lm := &LogManager{
		logCollection:  db.Collection("cron_job_logs"),
		logChan:        make(chan *model.JobLog, 1000),
		autoCommitChan: make(chan *model.JobLogBatch, 1000),
	}

	go lm.writeLoop()

	return lm
}

func (l *LogManager) AppendLog(log *model.JobLog) {
	select {
	case l.logChan <- log:
	default:
		// queue is full, throw it
	}
}

func (l *LogManager) saveLogs(batch *model.JobLogBatch) {
	_, err := l.logCollection.InsertMany(context.Background(), batch.Logs)
	if err != nil {
		fmt.Println(err)
	}
}

func (l *LogManager) writeLoop() {
	var logBatch *model.JobLogBatch

	for {
		select {
		case logRes := <-l.logChan:
			var timer *time.Timer
			if logBatch == nil {
				logBatch = &model.JobLogBatch{}
				timer = time.AfterFunc(
					time.Duration(conf.GetConf().Mongo.LogCommitTimeout)*time.Second,
					func() func() {
						return func() {
							l.autoCommitChan <- logBatch
						}
					}())
			}

			logBatch.Logs = append(logBatch.Logs, logRes)

			if len(logBatch.Logs) >= conf.GetConf().Mongo.LogBatchSize {
				l.saveLogs(logBatch)
				logBatch = nil
				timer.Stop()
			}
		case timeoutBatch := <-l.autoCommitChan:
			if timeoutBatch != logBatch {
				continue
			}
			l.saveLogs(timeoutBatch)
			logBatch = nil
		}
	}
}
