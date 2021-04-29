package network

import (
  "bytes"
  "encoding/json"
  "errors"
  "io/ioutil"
  "net"
  "net/http"

  "gitlab.jiangxingai.com/jxserving/components/assistant/pkg/config"
)

type DetectResultBody struct {
  Result [][]string `json:"result"`
}

/**
  切换模型
*/

func SwitchModelHandler(object interface{}) error {
  var contentType = "application/json"
  host, port := config.JxServingPortal()
  url := "http://" + net.JoinHostPort(host, port) + "/api/v1/switch"
  b, err := json.Marshal(object)
  if err != nil {
    return err
  }
  var data = bytes.NewBuffer([]byte(b))
  resp, err := http.Post(url, contentType, data)
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  if resp.StatusCode != http.StatusOK {
    return errors.New("jxServing switch model error")
  }
  return nil
}

/**
 * 图片识别
 */
func DetectImage(object interface{}) (interface{}, error) {
  //base64Str, err := utils.ImageToBase64Str(imagePath)
  var contentType = "application/json"
  host, port := config.JxServingPortal()
  url := "http://" + net.JoinHostPort(host, port) + "/api/v1/detect"
  b, err := json.Marshal(object)
  if err != nil {
    return nil, err
  }
  var data = bytes.NewBuffer([]byte(b))
  resp, err := http.Post(url, contentType, data)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  if resp.StatusCode != http.StatusOK {
    result, _ := ioutil.ReadAll(resp.Body)
    return nil, errors.New(string(result))
  }
  res, err := ioutil.ReadAll(resp.Body)
  var body DetectResultBody
  err = json.Unmarshal(res, &body)
  if err != nil {
    return nil, err
  }
  return body, nil
}
