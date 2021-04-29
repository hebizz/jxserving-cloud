package config

import (
  "fmt"

  "gitlab.jiangxingai.com/jxserving/components/utils"
  "k8s.io/klog"
)

var (
  buildVersion string
  buildTime    string
  buildHash    string
  goVersion    string
)

const (
  serverRunModeKey      = "model.server.RunMode"
  defaultServerRunMode  = "debug"
  serverHttpPortKey     = "model.server.HttpPort"
  defaultServerHttpPort = ":9004"
  modelZipSavePath      = "model.app.ModelSavePath"
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
    klog.Info("runMode:::", ret)

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

func ModelSavePath() string {
  return utils.ReadString(modelZipSavePath)
}
