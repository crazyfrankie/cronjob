package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/crazyfrankie/cronjob/pkg/model"
)

type LogService struct {
	logCollection *mongo.Collection
}

func NewLogService(db *mongo.Database) *LogService {
	logSvc := &LogService{
		logCollection: db.Collection("cron_job_logs"),
	}

	go logSvc.ensureIndexes()

	return logSvc
}

func (s *LogService) ensureIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobNameIndex := mongo.IndexModel{
		Keys: map[string]interface{}{
			"jobName": 1,
		},
	}

	startTimeIndex := mongo.IndexModel{
		Keys: map[string]interface{}{
			"startTime": -1,
		},
	}

	compoundIndex := mongo.IndexModel{
		Keys: map[string]interface{}{
			"jobName":   1,
			"startTime": -1,
		},
	}

	_, _ = s.logCollection.Indexes().CreateOne(ctx, jobNameIndex)
	_, _ = s.logCollection.Indexes().CreateOne(ctx, startTimeIndex)
	_, _ = s.logCollection.Indexes().CreateOne(ctx, compoundIndex)
}

func (s *LogService) GetJobLogs(ctx context.Context, name string, skip int, limit int) ([]model.JobLog, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var res []model.JobLog

	sortOpt := options.Find().
		SetSort(map[string]interface{}{"startTime": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	filter := &model.JobLogFilter{JobName: name}
	if name == "" {
		filter = &model.JobLogFilter{}
	}

	cursor, err := s.logCollection.Find(timeoutCtx, filter, sortOpt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(timeoutCtx)

	resultChan := make(chan []model.JobLog, 1)
	errorChan := make(chan error, 1)

	go func() {
		var localRes []model.JobLog

		if limit > 0 {
			localRes = make([]model.JobLog, 0, limit)
		} else {
			localRes = make([]model.JobLog, 0, 20)
		}

		for cursor.Next(timeoutCtx) {
			item := model.JobLog{}

			if err := cursor.Decode(&item); err != nil {
				continue
			}

			localRes = append(localRes, item)

			if limit > 0 && len(localRes) >= limit {
				break
			}
		}

		if err := cursor.Err(); err != nil {
			errorChan <- err
			return
		}

		resultChan <- localRes
	}()

	select {
	case <-timeoutCtx.Done():
		return nil, timeoutCtx.Err()
	case err := <-errorChan:
		return nil, err
	case res = <-resultChan:
		return res, nil
	}
}
