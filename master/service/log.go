package service

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/crazyfrankie/cronjob/common/model/resp"
)

type LogService struct {
	logCollection *mongo.Collection
}

func NewLogService(db *mongo.Database) *LogService {
	return &LogService{
		logCollection: db.Collection("cron_job_logs"),
	}
}

func (s *LogService) GetJobLogs(ctx context.Context, name string, skip int, limit int) ([]resp.JobLog, error) {
	var res []resp.JobLog

	sortOpt := options.Find().SetSort(map[string]interface{}{"startTime": -1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	filter := bson.M{"jobName": name}
	cursor, err := s.logCollection.Find(ctx, filter, sortOpt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		item := resp.JobLog{}

		if err := cursor.Decode(&item); err != nil {
			continue
		}

		res = append(res, item)
	}

	return res, nil
}
