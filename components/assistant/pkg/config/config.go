package config

import (
  "fmt"
  "os"

  "gitlab.jiangxingai.com/jxserving/components/utils"
  log "k8s.io/klog"
)

var (
  buildVersion string
  buildTime    string
  buildHash    string
  goVersion    string
)

const (
  serverRunModeKey      = "assistant.server.RunMode"
  defaultServerRunMode  = "debug"
  serverHttpPortKey     = "assistant.server.HttpPort"
  defaultServerHttpPort = ":10000"

  reportServerHostKey = "assistant.report-backend.host"
  reportServerPortKey = "assistant.report-backend.port"

  dataSetServerHostKey = "assistant.dataset-backend.host"
  dataSetServerPortKey = "assistant.dataset-backend.port"

  jxservingHostKey = "assistant.jxserving-backend.host"
  jxservingPortKey = "assistant.jxserving-backend.port"
)

func Version() string {
  ver := fmt.Sprintf(
    "Version: %s\nBuild Time: %s\nBuild Hash: %s\nGo Version: %s\n",
    buildVersion,
    buildTime,
    buildHash,
    goVersion)
  return ver
}

func RunMode() string {
  if ret := utils.ReadString(serverRunModeKey); ret != "" {
    return ret
  } else {
    return defaultServerRunMode
  }
}

func HttpPort() string {
  if ret := utils.ReadString(serverHttpPortKey); ret != "" {
    return ret
  } else {
    return defaultServerHttpPort
  }
}

func DataSetBackendPortal() (string, string) {
  port := os.Getenv("DATASET_PORT")
  host := os.Getenv("DATASET_HOST")
  log.Info("dataset env host & port:", host, port)
  if port == "" {
    host = utils.ReadString(dataSetServerHostKey)
    port = utils.ReadString(dataSetServerPortKey)
  }
  return host, port
}

func JxServingPortal() (string, string) {
  host := utils.ReadString(jxservingHostKey)
  port := utils.ReadString(jxservingPortKey)
  return host, port
}

func ReportPoliceDataToWechat() string {
  return os.Getenv("REPORT_URL")
}
