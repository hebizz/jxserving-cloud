package response

import (
  "fmt"
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"
  log "k8s.io/klog"
)

type Response struct {
  Code    string      `json:"c"`
  Message string      `json:"msg"`
  Data    interface{} `json:"data"`
}

func ServerError(c *gin.Context, code int, msg string, err error, d interface{}) {
  response := &Response{
    Code:    strconv.Itoa(code),
    Message: msg,
    Data:    d,
  }
  if err != nil {
    response.Message = fmt.Sprintf("%s", err)
  }
  log.Info("response error：", err)
  if d != "" {
    log.Info("response error data:", d)
  }
  c.JSON(code, response)
}

func Success(c *gin.Context, msg string, d interface{}) {
  response := &Response{
    Code:    strconv.Itoa(http.StatusOK),
    Message: msg,
    Data:    d,
  }
  log.Info("response success:", d)
  c.JSON(http.StatusOK, response)
}
