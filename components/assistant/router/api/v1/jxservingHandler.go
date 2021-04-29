package v1

import (
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/assistant/pkg/network"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
)

type Model struct {
  Id string `json:"id" binding:"required"`
}

type ImageInfo struct {
  Image string `json:"image" binding:"required"`
}

/**
  switch model
*/
func JxServingSwitchModelHandler(c *gin.Context) {
  var body Model
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  err = network.SwitchModelHandler(body)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  response.Success(c, fmt.Sprintf("switch model (%s) success", body.Id), "")
}

/**
  detect image
*/
func JxServingDetectImageHandler(c *gin.Context) {
  var body ImageInfo
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  result, err := network.DetectImage(body)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  response.Success(c, "detect image success", result)
}
