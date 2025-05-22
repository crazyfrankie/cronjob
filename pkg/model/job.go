package model

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
)

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

type JobEvent struct {
	EventType int //  SAVE, DELETE
	Job       *Job
}

type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo // 执行状态
	Output      []byte          // 脚本输出
	Err         error           // 脚本错误原因
	StartTime   time.Time       // 启动时间
	EndTime     time.Time       // 结束时间
}

type JobScheduleInfo struct {
	Job     *Job
	EntryID cron.EntryID // cron任务ID
}

type JobExecuteInfo struct {
	Job        *Job               // 任务信息
	RealTime   time.Time          // 实际的调度时间
	CancelCtx  context.Context    // 任务command的context
	CancelFunc context.CancelFunc //  用于取消command执行的cancel函数
}

type JobLogFilter struct {
	JobName string `bson:"jobName"`
}
