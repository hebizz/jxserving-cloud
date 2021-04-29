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
  serverRunModeKey       = "dataset.server.RunMode"
  defaultServerRunMode   = "debug"
  serverHttpPortKey      = "dataset.server.HttpPort"
  defaultServerHttpPort  = ":9002"
  ImageSourceSavePath    = "dataset.storage.ImageSourceSavePath"
  VocDataXmlSavePath     = "dataset.storage.VocDataXmlSavePath"
  VocDataUploadZipPath   = "dataset.storage.VocDataUploadZipPath"
  VocDataUploadUnZipPath = "dataset.storage.VocDataUploadUnZipPath"
  CocoDataUploadPath     = "dataset.storage.CocoDataUploadPath"
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

func ImageSavePath() string {
  return utils.ReadString(ImageSourceSavePath)
}

func VOCDataSavePath() string {
  return utils.ReadString(VocDataXmlSavePath)
}

func VOCDataUploadZipPath() string {
  return utils.ReadString(VocDataUploadZipPath)
}
func VOCDataUploadUnzipPath() string {
  return utils.ReadString(VocDataUploadUnZipPath)
}

func COCODataUploadPath() string {
  return utils.ReadString(CocoDataUploadPath)
}
