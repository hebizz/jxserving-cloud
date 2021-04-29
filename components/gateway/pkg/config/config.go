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
  serverRunModeKey      = "gateway.server.RunMode"
  defaultServerRunMode  = "debug"
  serverHttpPortKey     = "gateway.server.HttpPort"
  defaultServerHttpPort = ":9000"

  datasetServerHostKey   = "gateway.route.dataset.host"
  datasetServerPortKey   = "gateway.route.dataset.port"
  assistantServerHostKey = "gateway.route.assistant.host"
  assistantServerPortKey = "gateway.route.assistant.port"
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

func DatasetPortal() (host string, port string) {
  host = utils.ReadString(datasetServerHostKey)
  port = utils.ReadString(datasetServerPortKey)
  return
}

func AssistantPortal() (host string, port string) {
  host = utils.ReadString(assistantServerHostKey)
  port = utils.ReadString(assistantServerPortKey)
  return
}
