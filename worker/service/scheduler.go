package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/crazyfrankie/cronjob/pkg/consts"
	"github.com/crazyfrankie/cronjob/pkg/model"
)

type Scheduler struct {
	cron            *cron.Cron
	mu              sync.Mutex
	jobEventChan    chan *model.JobEvent
	jobResultChan   chan *model.JobExecuteResult
	jobPlanTable    map[string]*model.JobScheduleInfo
	jobExecuteTable map[string]*model.JobExecuteInfo
	log             *LogManager
	lock            *LockManager
}

func InitScheduler(log *LogManager, lock *LockManager) *Scheduler {
	cronOptions := []cron.Option{
		cron.WithSeconds(), // 秒级别调度
		cron.WithChain( // 任务执行链，可以添加一些中间件功能
			cron.Recover(cron.DefaultLogger), // 捕获并记录panic
		),
	}
	s := &Scheduler{
		cron:            cron.New(cronOptions...),
		jobEventChan:    make(chan *model.JobEvent),
		jobResultChan:   make(chan *model.JobExecuteResult),
		jobPlanTable:    make(map[string]*model.JobScheduleInfo),
		jobExecuteTable: make(map[string]*model.JobExecuteInfo),
		log:             log,
		lock:            lock,
	}

	go s.watchEvent()

	return s
}

func (s *Scheduler) watchEvent() {
	s.cron.Start()

	for {
		select {
		case evt := <-s.jobEventChan:
			s.handleJobEvent(evt)
		case evt := <-s.jobResultChan:
			s.handleJobResult(evt)
		}
	}
}

func (s *Scheduler) handleJobEvent(evt *model.JobEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch evt.EventType {
	case consts.JobSave:
		jobSch := &model.JobScheduleInfo{
			Job: evt.Job,
		}
		entryID, err := s.cron.AddFunc(jobSch.Job.CronExpr, func() {
			s.TryStartJob(jobSch)
		})
		if err != nil {
			return
		}
		jobSch.EntryID = entryID
		s.jobPlanTable[evt.Job.Name] = jobSch
	case consts.JobDel:
		if exists, ok := s.jobPlanTable[evt.Job.Name]; ok {
			s.cron.Remove(exists.EntryID)
			delete(s.jobPlanTable, evt.Job.Name)
		}
	case consts.JobKill:
		if exists, ok := s.jobExecuteTable[evt.Job.Name]; ok {
			exists.CancelFunc()
		}
	}
}

func (s *Scheduler) handleJobResult(evt *model.JobExecuteResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.jobExecuteTable, evt.ExecuteInfo.Job.Name)

	if !errors.Is(evt.Err, ErrFailedToPreemptLock) {
		jobLog := &model.JobLog{
			JobName:      evt.ExecuteInfo.Job.Name,
			Command:      evt.ExecuteInfo.Job.Command,
			Output:       string(evt.Output),
			ScheduleTime: evt.ExecuteInfo.RealTime.Unix(),
			StartTime:    evt.StartTime.Unix(),
			EndTime:      evt.EndTime.Unix(),
		}
		if evt.Err != nil {
			jobLog.Err = evt.Err.Error()
		} else {
			jobLog.Err = ""
		}
		s.log.AppendLog(jobLog)
	}
}

func (s *Scheduler) TryStartJob(sch *model.JobScheduleInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.jobExecuteTable[sch.Job.Name]
	if ok {
		return
	}
	executeInfo := &model.JobExecuteInfo{
		Job:      sch.Job,
		RealTime: time.Now(),
	}
	executeInfo.CancelCtx, executeInfo.CancelFunc = context.WithCancel(context.Background())

	s.jobExecuteTable[executeInfo.Job.Name] = executeInfo

	fmt.Println("执行任务:", executeInfo.Job.Name, executeInfo.RealTime)
	// execute job
	go func() {
		res := &model.JobExecuteResult{
			ExecuteInfo: executeInfo,
			StartTime:   time.Now(),
			Output:      make([]byte, 0),
		}
		lock := s.lock.CreateJobLock(executeInfo.Job.Name)

		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		err := lock.TryLock()
		if err != nil {
			lock.Unlock()
			res.Err = err
			res.EndTime = time.Now()
		} else {
			res.StartTime = time.Now()

			cmd := exec.CommandContext(executeInfo.CancelCtx, "/bin/bash", "-c", executeInfo.Job.Command)
			output, err := cmd.CombinedOutput()
			res.Err = err
			res.Output = output
			res.EndTime = time.Now()
			lock.Unlock()
		}
		s.PushJobResult(res)
	}()
}

func (s *Scheduler) PushJobEvent(jobEvt *model.JobEvent) {
	s.jobEventChan <- jobEvt
}

func (s *Scheduler) PushJobResult(result *model.JobExecuteResult) {
	s.jobResultChan <- result
}
