package handler

import (
	"net/http"
	"strconv"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/cronjob/common/model/req"
	"github.com/crazyfrankie/cronjob/common/response"
	"github.com/crazyfrankie/cronjob/master/service"
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
		jobGroup.POST("", h.CreateJob())
		jobGroup.DELETE("delete")
		jobGroup.GET("")
		jobGroup.DELETE("kill")
		jobGroup.POST("", h.CreateJob())
		jobGroup.GET("logs", h.GetJobLogs())
		jobGroup.GET("workers")
	}
}

func (h *JobHandler) CreateJob() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReq req.CreateJobReq
		if err := c.ShouldBind(&createReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

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

		res, err := h.logSvc.GetJobLogs(c.Request.Context(), name, skip, limit)
		if err != nil {
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
