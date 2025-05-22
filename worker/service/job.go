package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"

	"github.com/crazyfrankie/cronjob/pkg/consts"
	"github.com/crazyfrankie/cronjob/pkg/model"
)

var (
	ErrNoLocalIPFound      = errors.New("no NIC IP found")
	ErrFailedToPreemptLock = errors.New("lock: Failed Lock Grab")
)

type JobManager struct {
	scheduler *Scheduler
	client    *clientv3.Client
	ip        string
}

func InitJobManager(client *clientv3.Client, scheduler *Scheduler) error {
	job := &JobManager{client: client, scheduler: scheduler}

	ip, err := getLocalIP()
	if err != nil {
		return err
	}
	job.ip = ip

	err = job.registerWorker()
	if err != nil {
		return err
	}

	go job.watchJobs()

	go job.watchKiller()

	return nil
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	// take the ip of the first non-lo NIC.
	var ipv4 string
	for _, addr := range addrs {
		if ipNet, isIpNet := addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// skip ipv6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return ipv4, nil
			}
		}
	}

	return "", ErrNoLocalIPFound
}

func (j *JobManager) registerWorker() error {
	em, err := endpoints.NewManager(j.client, "cron/workers")
	if err != nil {
		return err
	}

	key := consts.JobWorkerDir + j.ip
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	leaseResp, err := j.client.Grant(ctx, 180)
	if err != nil {
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{Addr: j.ip}, clientv3.WithLease(leaseResp.ID))

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ch, err := j.client.KeepAlive(ctx, leaseResp.ID)
		if err != nil {
			log.Printf("keep alive failed lease id:%d", leaseResp.ID)
			return
		}
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					log.Println("KeepAlive channel closed")
					return
				}
				fmt.Println("Lease renewed")
			case <-ctx.Done():
				return
			}
		}
	}()

	return err
}

func (j *JobManager) watchJobs() error {
	res, err := j.client.Get(context.Background(), consts.JobSaveDir, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range res.Kvs {
		var job model.Job
		if err := sonic.Unmarshal(kv.Value, &job); err != nil {
			continue
		}
		j.scheduler.PushJobEvent(&model.JobEvent{
			EventType: consts.JobSave,
			Job:       &job,
		})
	}

	go func() {
		newRevision := res.Header.Revision + 1
		watch := j.client.Watch(context.Background(), consts.JobSaveDir, clientv3.WithPrefix(), clientv3.WithRev(newRevision))
		// handle watch event
		for watchRes := range watch {
			for _, evt := range watchRes.Events {
				switch evt.Type {
				case mvccpb.PUT:
					job := model.Job{}
					if err := sonic.Unmarshal(evt.Kv.Value, &job); err != nil {
						continue
					}
					jobEvt := &model.JobEvent{
						EventType: consts.JobSave,
						Job:       &job,
					}
					j.scheduler.PushJobEvent(jobEvt)
				case mvccpb.DELETE:
					jobName := strings.TrimPrefix(string(evt.Kv.Key), consts.JobSaveDir)
					j.scheduler.PushJobEvent(&model.JobEvent{
						EventType: consts.JobDel,
						Job: &model.Job{
							Name: jobName,
						},
					})
				}
			}
		}
	}()

	return nil
}

func (j *JobManager) watchKiller() {
	go func() {
		watch := j.client.Watch(context.Background(), consts.JobKillDir, clientv3.WithPrefix())
		for watchRes := range watch {
			for _, evt := range watchRes.Events {
				switch evt.Type {
				case mvccpb.PUT:
					jobName := strings.TrimPrefix(string(evt.Kv.Key), consts.JobKillDir)
					j.scheduler.PushJobEvent(&model.JobEvent{
						EventType: consts.JobKill,
						Job: &model.Job{
							Name: jobName,
						},
					})
				case mvccpb.DELETE:
				}
			}
		}
	}()
}
