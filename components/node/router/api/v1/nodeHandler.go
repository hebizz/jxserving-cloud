package v1

import (
  "net/http"
  "strings"

  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "gitlab.jiangxingai.com/jxserving/components/node/database"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
  "gitlab.jiangxingai.com/jxserving/components/utils"
  "k8s.io/klog"
)

/**
  添加wireGuard节点
*/
func AddNodeHandler(c *gin.Context) {
  var body interfaces.NodeInfo
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "query node list params error", err, "")
    return
  }
  // 校验cidr输入格式
  if isMatch := utils.RegexpExecute(body.Cidr); !isMatch {
    response.ServerError(c, http.StatusBadRequest, "cidr is not correct, example 9.101.0.1/32", nil, "")
    return
  }
  num, err := database.QueryNodeByCidr(body.Cidr)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "insert node info error", err, "")
    return
  }
  if num > 0 {
    response.ServerError(c, http.StatusForbidden, "this cidr is exist", nil, "")
    return
  }
  // 插入数据库
  err = database.InsertNodeInfo(body)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "insert node info error", err, "")
    return
  }
  response.Success(c, "add node list success", "")
}

/**
  查询node list信息
*/
func QueryNodeListHandler(c *gin.Context) {
  // 查询数据库
  infoList, err := database.QueryNodeInfoList()
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "query node list error", err, "")
    return
  }
  // 获取本机的wg状态
  contents, err := utils.ExecCommandAndPrintLines("wg", "show")
  if err != nil {
    klog.Error("wg show error: ", err)
    response.ServerError(c, http.StatusInternalServerError, "wg show error", err, "")
    return
  }
  var master interfaces.NodeInfo
  if contents.Len() == 0 {
    master.Status = "已断开"
    master.PublicKey = "--"
  } else {
    master.Status = "已连接"
  }
  // 获取公钥
  for item := contents.Front(); nil != item; item = item.Next() {
    res := strings.Split(item.Value.(string), ":")
    r := strings.TrimSpace(res[0])
    if r == "public key" {
      publicKey := strings.TrimSpace(res[1])
      master.PublicKey = publicKey
    }
  }
  // cidr 直接给定值，wg show 无法获取wg ip,只能获取所属网段
  master.Cidr = "9.101.0.100/32"
  // 对每个节点进行网络测速
  var lists = make([]interfaces.NodeInfo, 0)
  for _, node := range infoList {
    // ping
    cidr := strings.Split(node.Cidr, "/")[0]
    contents, _ := utils.ExecCommandAndPrintLines("ping", cidr, "-w", "1")
    for item := contents.Front(); nil != item; item = item.Next() {
      if strings.Contains(item.Value.(string), "time=") {
        node.Status = "已连接"
        res := strings.Split(item.Value.(string), "time=")
        if res != nil {
          node.NetSpeed = strings.TrimSpace(res[1])
        }
        break
      } else {
        node.Status = "已断开"
        node.NetSpeed = "--"
      }
    }
    lists = append(lists, node)
  }
  responseBody := &interfaces.NodeList{Master: master, NodeList: lists}
  response.Success(c, "query node list success", responseBody)
}

/**
  delete node
*/
func DeleteNodeHandler(c *gin.Context) {
  var body interfaces.NodeDelete
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "request for deleteNode params error", err, "")
    return
  }
  err = database.DeleteNode(body.Id)
  if err != nil {
    response.ServerError(c, http.StatusNotFound, "delete node error", err, "")
    return
  }
  response.Success(c, "delete node success", "")
}

/**
  update node info
*/
func UpdateNodeInfoHandler(c *gin.Context) {
  var body interfaces.NodeInfo
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "query node list params error", err, "")
    return
  }
  // 校验cidr输入格式
  if isMatch := utils.RegexpExecute(body.Cidr); !isMatch {
    response.ServerError(c, http.StatusBadRequest, "cidr is not correct, example 9.101.0.1/32", nil, "")
    return
  }
  err = database.UpdateNodeInfo(body)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "update node error", err, "")
    return
  }
  response.Success(c, "update node success", "")
}
