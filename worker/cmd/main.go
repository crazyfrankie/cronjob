package main

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/crazyfrankie/cronjob/worker/conf"
	"github.com/crazyfrankie/cronjob/worker/service"
)

func main() {
	mg := initMongo()
	client := initETCD()
	log := service.NewLogManager(mg)
	lock := service.NewLockManager(client)
	scheduler := service.InitScheduler(log, lock)
	err := service.InitJobManager(client, scheduler)
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
	return
}

func initETCD() *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.GetConf().ETCD.Endpoint,
		DialTimeout: time.Duration(conf.GetConf().ETCD.ConnectTimeOut),
	})
	if err != nil {
		panic(err)
	}

	return client
}

func initMongo() *mongo.Database {
	monitor := &event.CommandMonitor{}

	opt := options.Client().ApplyURI(conf.GetConf().Mongo.URI).
		SetMonitor(monitor).SetConnectTimeout(time.Duration(conf.GetConf().Mongo.ConnectTimeOut) * time.Second)
	client, err := mongo.Connect(opt)
	if err != nil {
		panic(err)
	}

	db := client.Database("cronjob")

	return db
}
