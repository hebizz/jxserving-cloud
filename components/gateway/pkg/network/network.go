package network

import (
  "bytes"
  "encoding/json"
  "errors"
  "io/ioutil"
  "net"
  "net/http"

  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/gateway/pkg/config"
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  log "k8s.io/klog"
)

type ResponseBody struct {
  C    string `json:"c"`
  Msg  string `json:"msg"`
  Data string `json:"data"`
}

/**
  上传告警数据到数据集管理
*/
func DatasetStoreRawDataHandler(c *gin.Context) (string, error) {
  host, port := config.DatasetPortal()
  url := "http://" + net.JoinHostPort(host, port) + "/api/v1/reportData/upload"
  log.Info("dataset post url:", url)
  raw, err := c.GetRawData()
  if err != nil {
    log.Error("get raw data  error: ", err)
    return "", err
  }
  data := bytes.NewBuffer(raw)
  resp, err := http.Post(url, "application/json", data)
  if err != nil {
    log.Error("post dataSets error: ", err)
    return "", err
  }
  defer resp.Body.Close()
  result, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Error("read dataSets response body error: ", err)
    return "", err
  }
  if resp.StatusCode != http.StatusOK {
    log.Error("dataset return error: ", string(result))
    return "", errors.New("dataset didn't response 200OK")
  }
  var response ResponseBody
  err = json.Unmarshal([]byte(string(result)), &response)
  if err != nil {
    log.Error("json unMarshal result error: ", err)
    return "", err
  }
  c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(raw))
  c.Next()
  return response.Data, nil
}

/**
  assistant is alive
*/
func IsAliveHandler() error {
  host, port := config.AssistantPortal()
  url := "http://" + net.JoinHostPort(host, port) + "/api/v1/assistant/ping"
  log.Info("ping assistant url :", url)
  resp, err := http.Get(url)
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  if resp.StatusCode != http.StatusOK {
    return errors.New("assistant server didn't working")
  }
  return nil
}

/**
  上传告警数据到assistant
*/
func PublishRawDataToAssistantHandler(c *gin.Context, rawDataId string) error {
  host, port := config.AssistantPortal()
  url := "http://" + net.JoinHostPort(host, port) + "/api/v1/assistant/verifyInfo"
  var body interfaces.ReportData
  _ = c.ShouldBindJSON(&body)
  if rawDataId != "" {
    body.DataId = rawDataId
  }
  b, err := json.Marshal(body)
  if err != nil {
    return err
  }
  data := bytes.NewBuffer(b)
  resp, err := http.Post(url, "application/json", data)
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  if resp.StatusCode != http.StatusOK {
    return errors.New("assistant verifyInfo didn't response 200OK")
  }
  return nil
}
