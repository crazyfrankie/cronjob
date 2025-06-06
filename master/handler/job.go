package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/cronjob/master/service"
	"github.com/crazyfrankie/cronjob/pkg/model"
	"github.com/crazyfrankie/cronjob/pkg/response"
)

type JobHandler struct {
	jobSvc    *service.JobService
	logSvc    *service.LogService
	workerSvc *service.WorkerService
}

func NewJobHandler(job *service.JobService, log *service.LogService, worker *service.WorkerService) *JobHandler {
	return &JobHandler{
		jobSvc:    job,
		logSvc:    log,
		workerSvc: worker,
	}
}

func (h *JobHandler) RegisterRoute(r *gin.Engine) {
	jobGroup := r.Group("api/job")
	{
		jobGroup.POST("", h.SaveJob())
		jobGroup.DELETE("delete", h.DeleteJob())
		jobGroup.GET("", h.GetJobList())
		jobGroup.DELETE("kill", h.KillJob())
		jobGroup.GET("logs", h.GetJobLogs())
		jobGroup.GET("workers", h.GetWorkerList())
	}
}

func (h *JobHandler) SaveJob() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReq model.Job
		if err := c.ShouldBind(&createReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := h.jobSvc.SaveJob(c.Request.Context(), &createReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

func (h *JobHandler) DeleteJob() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")

		err := h.jobSvc.DeleteJob(c.Request.Context(), name)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

func (h *JobHandler) GetJobList() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := h.jobSvc.GetJobList(c.Request.Context())
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.SuccessWithData(c, res)
	}
}

func (h *JobHandler) KillJob() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")

		err := h.jobSvc.KillJob(c.Request.Context(), name)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

func (h *JobHandler) GetJobLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		skipParam := c.Query("skip")
		limitParam := c.Query("limit")

		var err error
		var skip, limit int
		if skip, err = strconv.Atoi(skipParam); err != nil {
			skip = 0
		}
		if limit, err = strconv.Atoi(limitParam); err != nil {
			limit = 20
		}

		// 创建具有较短超时的上下文，防止请求挂起太长时间
		timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), 6*time.Second)
		defer cancel()

		// 使用带超时的上下文调用服务
		res, err := h.logSvc.GetJobLogs(timeoutCtx, name, skip, limit)
		if err != nil {
			// 判断是否是超时错误
			if err == context.DeadlineExceeded {
				response.Error(c, http.StatusGatewayTimeout, gerrors.NewBizError(30001, "查询日志超时，请稍后重试或缩小查询范围"))
				return
			}
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.SuccessWithData(c, res)
	}
}

func (h *JobHandler) GetWorkerList() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := h.workerSvc.ListWorkers(c.Request.Context())
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.SuccessWithData(c, res)
	}
}
