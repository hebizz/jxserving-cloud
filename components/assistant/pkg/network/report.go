package network

import (
  "bytes"
  "crypto/tls"
  "encoding/json"
  "errors"
  "fmt"
  "io/ioutil"
  "net"
  "net/http"

  "gitlab.jiangxingai.com/jxserving/components/assistant/model/database"
  "gitlab.jiangxingai.com/jxserving/components/assistant/pkg/config"
  "k8s.io/klog"
)

/**
  更新AI数据
*/
func PostIgnoreInfo(projectName string, object interface{}) (interface{}, error, int) {
  var contentType = "application/json"
  err, project := database.QueryProjectInfoByName(projectName)
  if err != nil {
    return "request error: ", err, http.StatusInternalServerError
  }
  url := "http://" + net.JoinHostPort(project.Host, project.Port) + project.Api
  klog.Info("update info to assistant backend url:", url)
  b, err := json.Marshal(object)
  if err != nil {
    return "marshal failed", err, http.StatusInternalServerError
  }
  var data = bytes.NewBuffer([]byte(b))
  resp, err := http.Post(url, contentType, data)
  if err != nil {
    return "request error: ", err, http.StatusInternalServerError
  }
  defer resp.Body.Close()
  statusCode := resp.StatusCode
  result, err := ioutil.ReadAll(resp.Body)
  return string(result), err, statusCode
}

/**
  新建AI数据
*/
func CreateInfo(projectName string, object interface{}) (interface{}, error, int) {
  var contentType = "application/json"
  err, project := database.QueryProjectInfoByName(projectName)
  if err != nil {
    return "request error: ", err, http.StatusInternalServerError
  }
  url := "http://" + net.JoinHostPort(project.Host, project.Port) + "/api/v1/manual/create"
  b, err := json.Marshal(object)
  if err != nil {
    return "umarshal failed", err, http.StatusInternalServerError
  }
  var data = bytes.NewBuffer([]byte(b))
  resp, err := http.Post(url, contentType, data)
  if err != nil {
    return "request error: ", err, http.StatusInternalServerError
  }
  defer resp.Body.Close()
  statusCode := resp.StatusCode
  result, err := ioutil.ReadAll(resp.Body)
  return string(result), nil, statusCode
}

/**
  人工干预后更新dataSet数据库
*/
func ReportUpdateDataToDataSet(obj interface{}) error {
  var contentType = "application/json"
  host, port := config.DataSetBackendPortal()
  url := "http://" + net.JoinHostPort(host, port) + "/api/v1/reportData/upload"
  b, err := json.Marshal(obj)
  if err != nil {
    return err
  }
  var data = bytes.NewBuffer(b)
  resp, err := http.Post(url, contentType, data)
  if err != nil {
    klog.Info("post error")
    return err
  }
  defer resp.Body.Close()
  statusCode := resp.StatusCode
  if statusCode != http.StatusOK {
    result, _ := ioutil.ReadAll(resp.Body)
    return errors.New(string(result))
  }
  return nil
}

/**
  上报告警信息到wechat
*/
func ReportWarnInfoToWechat() error {
  var contentType = "application/json"
  url := config.ReportPoliceDataToWechat()
  fmt.Println("wechat machine url :", url)
  msg := `{"msgtype":"text","text":{"content":"干预平台收到一条新告警，请打开srv.cloud.jiangxingai.com:3000网站进行操作"}}`
  var data = bytes.NewBuffer([]byte(msg))

  // 调用第三方api,会默认调用根证书rootca.pem进行验证，会与该api的官方证书产生冲突，
  // 常见错误：x509: certificate signed by unknown authority
  tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
  client := &http.Client{Transport: tr}

  resp, err := client.Post(url, contentType, data)
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  _, err = ioutil.ReadAll(resp.Body)
  return err
}
