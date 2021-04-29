package config

import (
  "fmt"

  "gitlab.jiangxingai.com/jxserving/components/utils"
)

var (
  buildVersion string
  buildTime    string
  buildHash    string
  goVersion    string
)

const (
  serverRunModeKey      = "node.server.RunMode"
  defaultServerRunMode  = "debug"
  serverHttpPortKey     = "node.server.HttpPort"
  defaultServerHttpPort = ":8000"

  jxsOnboard            = "node.app.Onboard"
  defaultJxsOnboard     = ":8080/api/v1alpha/switch"
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

func JxsOnboardApi() string {
  if ret := utils.ReadString(jxsOnboard); ret != "" {
    return ret
  } else {
    return defaultJxsOnboard
  }
}