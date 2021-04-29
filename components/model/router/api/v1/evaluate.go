package v1

import (
  "net/http"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/spf13/cast"
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "gitlab.jiangxingai.com/jxserving/components/model/models/database"
  "gitlab.jiangxingai.com/jxserving/components/model/router/api/v1Branch"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
  "gitlab.jiangxingai.com/jxserving/components/utils"
  log "k8s.io/klog"
)

const AppKey = "app"

func EvaluateModel(c *gin.Context) {
  var body interfaces.Modeldate
  err := c.BindJSON(&body)
  body.DownType = 1
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  if ret := utils.ReadMap(AppKey); ret != nil {
    value := cast.ToStringMapString(ret)
    err, model := database.QueryModelInfo(body.Modelmd5)
    modelPath := model.ModelPath
    rawName := model.RawName
    //unzip model
    err = v1Branch.Zipmodel(modelPath, value, rawName)
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    path, err := v1Branch.Getdateset(body, value)
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    err = v1Branch.UnzipDateset(path, value)
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    date, test, err := v1Branch.Handler(value)
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    data := interfaces.ModelPerformance{
      ModelMd5:         body.Modelmd5,
      EvalDataset:      body.Name,
      Datesettype:      body.Type,
      ErrorRate:        test.ErrorRate,
      LeakRate:         test.LeakRate,
      MeanAp:           test.Map,
      CreatedTimestamp: time.Now().Unix(),
    }
    err = database.InsertPerformanceData(data)
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    log.Info("evaluate result:", test)
    response.Success(c, "evaluate success", map[string]interface{}{
      "result": date,
    })
  }

}

/**
  历史评价
*/
func HistoryEvaluate(c *gin.Context) {
  var body AuthId
  err := c.ShouldBindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  } else {
    performances, err := database.QueryPerformanceData(body.ModelMd5)
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    response.Success(c, "getmsg succ", map[string]interface{}{
      "result": performances,
    })
  }
}
