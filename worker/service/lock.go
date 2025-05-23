package service

import (
	"context"
	"fmt"
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/crazyfrankie/cronjob/pkg/consts"
)

type LockManager struct {
	client *clientv3.Client
}

type Lock struct {
	client     *clientv3.Client
	jobName    string
	cancelFunc context.CancelFunc // 用于终止自动续租
	leaseId    clientv3.LeaseID   // 租约ID
	isLocked   bool
}

func NewLockManager(client *clientv3.Client) *LockManager {
	return &LockManager{client: client}
}

func (l *LockManager) CreateJobLock(name string) *Lock {
	return &Lock{
		client:  l.client,
		jobName: name,
	}
}

func (l *Lock) TryLock() error {
	leaseResp, err := l.client.Grant(context.Background(), 5)
	if err != nil {
		return err
	}
	l.leaseId = leaseResp.ID

	ctx, cancel := context.WithCancel(context.Background())
	l.cancelFunc = cancel

	go func() {
		ch, err := l.client.KeepAlive(ctx, l.leaseId)
		if err != nil {
			log.Printf("keep alive failed lease id:%d", l.leaseId)
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

	txn := l.client.Txn(context.Background())

	key := consts.JobLockDir + l.jobName
	txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "", clientv3.WithLease(l.leaseId))).
		Else(clientv3.OpGet(key))

	// commit txn
	txnRes, err := txn.Commit()
	if err != nil {
		cancel()                                   // 取消自动续租
		l.client.Revoke(context.TODO(), l.leaseId) //  释放租约
		return err
	}

	if !txnRes.Succeeded {
		return ErrFailedToPreemptLock
	}

	l.isLocked = true
	return nil
}

func (l *Lock) Unlock() {
	if l.isLocked {
		l.cancelFunc()                             // 取消自动续租的协程
		l.client.Revoke(context.TODO(), l.leaseId) // 释放租约
	}
}
