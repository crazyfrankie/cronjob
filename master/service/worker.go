package service

import (
	"context"
	"strings"

	"go.etcd.io/etcd/client/v3"

	"github.com/crazyfrankie/cronjob/common/consts"
)

type WorkerService struct {
	client *clientv3.Client
}

func NewWorkerService(client *clientv3.Client) *WorkerService {
	return &WorkerService{
		client: client,
	}
}

func (s *WorkerService) ListWorkers(ctx context.Context) ([]string, error) {
	res, err := s.client.Get(ctx, consts.JobWorkerDir, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	addrs := make([]string, 0, len(res.Kvs))
	for _, kv := range res.Kvs {
		addrs = append(addrs, strings.TrimPrefix(string(kv.Key), consts.JobWorkerDir))
	}

	return addrs, nil
}
