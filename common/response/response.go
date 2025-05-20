package response

import (
	"net/http"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    20000,
		Message: "OK",
	})
}

func SuccessWithData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    20000,
		Message: "OK",
		Data:    data,
	})
}

func Error(c *gin.Context, status int, err error) {
	if bizErr, ok := gerrors.FromBizStatusError(err); ok {
		c.JSON(status, Response{
			Code:    bizErr.BizStatusCode(),
			Message: bizErr.BizMessage(),
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		Code:    50000,
		Message: err.Error(),
	})
}
