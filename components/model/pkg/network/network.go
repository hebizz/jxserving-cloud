package network

import (
  "io/ioutil"
  "net/http"
  "strings"

  log "k8s.io/klog"
)

func DatasetHandler(bodyJson []byte, url string) (string, error) {
  res, err := http.Post(url, "appliction/json", strings.NewReader(string(bodyJson)))
  if err != nil {
    log.Error(err)
    return "", err
  }
  data, err := ioutil.ReadAll(res.Body)
  if err != nil {
    log.Error(err)
    return "", err
  }
  return string(data), nil
}
