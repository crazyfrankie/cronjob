package consts

const (
	JobSaveDir   = "cron/jobs/"
	JobWorkerDir = "cron/workers/"
	JobKillDir   = "cron/killer/"
	JobLockDir   = "cron/lock/"
)

const (
	JobSave = iota
	JobDel
	JobKill
)
