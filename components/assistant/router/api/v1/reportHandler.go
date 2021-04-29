package v1

import (
  list2 "container/list"
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/assistant/model/database"
  "gitlab.jiangxingai.com/jxserving/components/assistant/pkg/config"
  "gitlab.jiangxingai.com/jxserving/components/assistant/pkg/network"
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
  log "k8s.io/klog"
)

type AuthBody struct {
  Host string `json:"host" binding:"required"`
}

type UpdateBody struct {
  Project  string                  `json:"project" binding:"required"`
  EventId  string                  `json:"event_id" binding:"required" `
  Ignore   bool                    `json:"ignore"`
  Override []interfaces.ReportData `json:"override"`
}

type CreateBody struct {
  Project  string                  `json:"project" binding:"required"`
  Override []interfaces.ReportData `json:"override" binding:"required"`
}

type ListBody struct {
  DataList []interfaces.ReportData
}

type RecoveryBody struct {
  Name       string `json:"name" binding:"required"`
  StartIndex int64  `json:"startIndex" binding:"required"`
  Offset     int64  `json:"offset" binding:"required"`
  AlertTime  int64  `json:"alertTime" binding:"required"`
}

var list = list2.New()

// 获取待验证信息
func PostVerifyInfoHandler(c *gin.Context) {
  var body interfaces.ReportData
  err := c.ShouldBindJSON(&body)
  log.Info("report origin data: ", body)
  if err != nil {
    log.Error("post info error:", err)
  }
  //存储数据到列表
  list.PushBack(body)
  response.Success(c, "save reportData success", "")
  err = network.ReportWarnInfoToWechat()
  if err != nil {
    log.Warning("Report warning msg to wechat failed: ", err)
  }
}

/**
  前端请求
*/
func GetVerifyInfoHandler(c *gin.Context) {
  if list.Len() == 0 {
    response.Success(c, "no new data", "")
    return
  }
  length := list.Len()
  var body []interface{}
  for i := 0; i < length; i++ {
    de := list.Front()
    body = append(body, de.Value)
    list.Remove(de)
  }
  log.Info("Data waiting to be processed: ", body)
  response.Success(c, fmt.Sprintf("There are %d more pieces of data", length), body)
}

/**
  忽略 & 更新AI识别信息
*/
func UpdateVerifyInfoHandler(c *gin.Context) {
  var body UpdateBody
  err := c.ShouldBindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, nil)
    return
  }
  // 查询reMark状态
  isReMark, err := database.QueryReportDataReMarkStatus(body.Project, body.EventId)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, nil)
    return
  }
  if isReMark {
    response.ServerError(c, http.StatusForbidden, "this data has been processed", nil, nil)
    return
  }
  overrides := make([]interfaces.ReportData, 0)
  for i := 0; i < len(body.Override); i++ {
    positions := make([]interfaces.AlertPosition, 0)
    info := body.Override[i]
    for j := 0; j < len(info.AlertPosition); j++ {
      if info.AlertPosition[j].RightY != "0" {
        positions = append(positions, info.AlertPosition[j])
      }
    }
    info.AlertPosition = positions
    overrides = append(overrides, info)
  }
  body.Override = overrides
  info, err, _ := network.PostIgnoreInfo(body.Project, body)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, nil)
    return
  } else {
    log.Info("report remote data success:", info)
    response.Success(c, "update info success", info)
  }
  if len(body.Override) > 0 {
    err = network.ReportUpdateDataToDataSet(body.Override[0])
    if err != nil {
      log.Error("update dataSets error: ", err)
    } else {
      log.Info("update dataSets success")
    }
  }
  // 更新reMark状态
  err = database.UpdateReportDataReMarkStatus(body.Project, body.EventId)
  if err != nil {
    log.Info("update reMark status error: ", err)
  } else {
    log.Info("update reMark status success")
  }
}

/**
  新建AI识别信息
*/
func CreateInfoHandler(c *gin.Context) {
  var body CreateBody
  _ = c.BindJSON(&body)
  result, err, _ := network.CreateInfo(body.Project, body)
  if err != nil {
    log.Info("createInfo error", err.Error())
    response.ServerError(c, http.StatusInternalServerError, "createInfo error", nil, nil)
  } else {
    response.Success(c, "createInfo success", result)
  }
}

/**
  上传数据到dataSet
*/
func UploadDataToDataSetHandler(c *gin.Context) {
  var body interfaces.ReportData
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, nil)
    return
  }
  err = network.ReportUpdateDataToDataSet(body)
  if err != nil {
    log.Info("report UpdateData To DataSet error", err.Error())
    response.ServerError(c, http.StatusInternalServerError, "", err, nil)
    return
  }
  response.Success(c, "update data success", "")
}

/*
   ping
*/
func Ping(c *gin.Context) {
  response.Success(c, "alive", config.Version())
}

/**
  查询projectList
*/
func QueryProjectListHandler(c *gin.Context) {
  lists, err := database.QueryProjectList()
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, nil)
    return
  }
  var projectList []string
  for _, project := range lists {
    projectList = append(projectList, project.Name)
  }
  response.Success(c, "query project list success", projectList)
}

/**
  一键恢复未处理的数据
*/
func RecoveryAssistantDataHandler(c *gin.Context) {
  var body RecoveryBody
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, nil)
    return
  }
  reportData, err := database.QueryNoMarkedDataByProject(body.Name, body.StartIndex, body.Offset, body.AlertTime, "alerttime")
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, nil)
  }
  response.Success(c, "recovery reportData list success", reportData)
}
