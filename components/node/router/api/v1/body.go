package v1

import (
  pt "gitlab.jiangxingai.com/jxserving/components/interfaces/protos"
)

type QueryListModelBody struct {
  Cidr string `json:"cidr" binding:"required"`
}

type NodeInfoBody struct {
  NodeStatus *pt.ResourcesReply `json:"status"`
  Version    string             `json:"version"`
  Models     []*pt.ModelInfo    `json:"models"`
}

type SubNodeInfoBody struct {
  Version string          `json:"version"`
  Models  []*pt.ModelInfo `json:"models"`
  Name    string          `json:"name"`
  Ip      string          `json:"ip"`
}
