package service

import (
	"context"
	
	clientv3 "go.etcd.io/etcd/client/v3"
)

type JobService struct {
	client *clientv3.Client
}

func NewJobService(client *clientv3.Client) *JobService {
	return &JobService{client: client}
}

func (s *JobService) CreateJob(ctx context.Context) error {
	return nil
}
