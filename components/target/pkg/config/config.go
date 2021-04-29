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
  serverRunModeKey       = "target.server.RunMode"
  defaultServerRunMode   = "debug"
  serverHttpPortKey      = "target.server.HttpPort"
  defaultServerHttpPort  = ":9003"
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