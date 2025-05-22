//go:build wireinject

package ioc

import (
	"context"
	"github.com/crazyfrankie/cronjob/master/conf"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/crazyfrankie/cronjob/master/handler"
	"github.com/crazyfrankie/cronjob/master/service"
)

func InitETCD() *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.GetConf().ETCD.Endpoint,
		DialTimeout: time.Duration(conf.GetConf().ETCD.ConnectTimeOut),
	})
	if err != nil {
		panic(err)
	}

	return client
}

func InitMongo() *mongo.Database {
	monitor := &event.CommandMonitor{}

	opt := options.Client().ApplyURI(conf.GetConf().Mongo.URI).
		SetMonitor(monitor).
		SetConnectTimeout(time.Duration(conf.GetConf().Mongo.ConnectTimeOut) * time.Second).
		SetServerSelectionTimeout(5 * time.Second)
	client, err := mongo.Connect(opt)
	if err != nil {
		panic(err)
	}

	db := client.Database("cronjob")
	err = InitMongoCollections(db)
	if err != nil {
		panic(err)
	}

	return db
}

func InitMongoCollections(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := db.CreateCollection(ctx, "cron_job_logs")
	if err != nil {
		return err
	}

	return nil
}

func InitGin(job *handler.JobHandler) *gin.Engine {
	srv := gin.Default()

	job.RegisterRoute(srv)

	return srv
}

func InitServer() *gin.Engine {
	wire.Build(
		InitMongo,
		InitETCD,

		service.NewJobService,
		service.NewLogService,
		service.NewWorkerService,

		handler.NewJobHandler,

		InitGin,
	)

	return new(gin.Engine)
}
