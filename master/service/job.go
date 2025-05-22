package service

import (
	"context"

	"github.com/bytedance/sonic"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/crazyfrankie/cronjob/pkg/consts"
	"github.com/crazyfrankie/cronjob/pkg/model"
)

type JobService struct {
	client *clientv3.Client
}

func NewJobService(client *clientv3.Client) *JobService {
	return &JobService{client: client}
}

func (s *JobService) SaveJob(ctx context.Context, job *model.Job) error {
	key := consts.JobSaveDir + job.Name

	data, err := sonic.Marshal(job)
	if err != nil {
		return err
	}

	if _, err := s.client.Put(ctx, key, string(data), clientv3.WithPrevKV()); err != nil {
		return err
	}

	return nil
}

func (s *JobService) DeleteJob(ctx context.Context, name string) error {
	key := consts.JobSaveDir + name
	_, err := s.client.Delete(ctx, key, clientv3.WithPrevKV())

	return err
}

func (s *JobService) GetJobList(ctx context.Context) ([]model.Job, error) {
	res, err := s.client.Get(ctx, consts.JobSaveDir, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	jobs := make([]model.Job, 0, len(res.Kvs))
	for _, kv := range res.Kvs {
		var job model.Job
		if err := sonic.Unmarshal(kv.Value, &job); err != nil {
			continue
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (s *JobService) KillJob(ctx context.Context, name string) error {
	key := consts.JobKillDir + name

	// Let the worker listen for a put operation,
	// and create a lease that expires automatically later.
	leaseResp, err := s.client.Grant(ctx, 1)
	if err != nil {
		return err
	}

	_, err = s.client.Put(ctx, key, "", clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	return nil
}
