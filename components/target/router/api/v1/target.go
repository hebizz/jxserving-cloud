package v1

import (
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
  "gitlab.jiangxingai.com/jxserving/components/target/models/database"
  log "k8s.io/klog"
)

type TargetBody struct {
  Id string `json:"id" binding:"required"`
}

/**
  获取标签
*/

func QueryTargetHandler(c *gin.Context) {
  // 查询所有target
  targets, err := database.QueryTargetData()
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  response.Success(c, "query target list success", targets)
}

/**
  上传target
*/
func UploadTargetHandler(c *gin.Context) {
  var target interfaces.Target
  err := c.BindJSON(&target)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  targets, err := database.QueryTargetData()
  for _, value := range targets {
    if value.Name == target.Name {
      response.ServerError(c, http.StatusBadRequest, fmt.Sprintf("%s is exists", target.Name), nil, "")
      return
    }
  }
  id, err := database.InsertTargetData(target)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  response.Success(c, fmt.Sprintf("upload target %s success", target.Name), &TargetBody{Id: id})
}

/**
  更新 target
*/
func UpdateTargetHandler(c *gin.Context) {
  var target interfaces.Target
  err := c.BindJSON(&target)
  log.Info("update target request:", target)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "request params error", nil, "")
    return
  }
  // 通过id 检查该document是否存在
  err, _ = database.QueryTarget(target.LabelId)
  if err != nil {
    response.ServerError(c, http.StatusNotFound, "this target id not exist", nil, "")
    return
  }
  // 通过id 查询target数据并更新
  err = database.UpdateTargetData(target)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  response.Success(c, fmt.Sprintf("update target %s success", target.Name), "")
}

/**
  删除target
*/
func DeleteTargetHandler(c *gin.Context) {
  var body TargetBody
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "request params error", nil, "")
    return
  }
  // 通过id 检查该document是否存在
  err, targetData := database.QueryTarget(body.Id)
  if err != nil {
    response.ServerError(c, http.StatusNotFound, "this target id not exist", nil, "")
    return
  }
  // 通过name 检查该document是否被model关联
  err = database.QueryTargetIsRelateByModel(targetData.Name)
  if err == nil {
    response.ServerError(c, http.StatusForbidden, "this label has been associated to model and not be deleted", nil, "")
    return
  } else {
    log.Info("query label is related by model error:", err)
  }
  // 通过id 删除target
  err = database.DeleteTargetData(body.Id)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "delete target by id error", nil, "")
    return
  }
  response.Success(c, fmt.Sprintf("delete target (%s) success", body.Id), "")
}
