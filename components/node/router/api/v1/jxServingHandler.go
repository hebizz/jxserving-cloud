package v1

import (
  "context"
  "net/http"
  "strings"
  "time"

  "github.com/gin-gonic/gin"
  pt "gitlab.jiangxingai.com/jxserving/components/interfaces/protos"
  "gitlab.jiangxingai.com/jxserving/components/node/database"
  "gitlab.jiangxingai.com/jxserving/components/node/router/api/v1Branch"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
  "gitlab.jiangxingai.com/jxserving/components/utils"
)

/**
  获取backend信息
*/
func QueryTargetNodeInfoHandler(c *gin.Context) {
  var body QueryListModelBody
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  if !utils.RegexpExecute(body.Cidr) {
    response.ServerError(c, http.StatusBadRequest, "cidr is wrong", nil, "")
    return
  }
  targetIp := strings.Split(body.Cidr, "/")[0]
  conn, err := v1Branch.ConnectGRPC(targetIp)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  defer conn.Close()
  replyConn, replyNode, replyModel, err := v1Branch.JxServingGRPC(conn)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  response.Success(c, "query backend info success", NodeInfoBody{Version: replyConn.Version, NodeStatus: replyNode, Models: replyModel.List})
}

/**
  查询子节点信息
*/
func QuerySubNodeModelInfo(c *gin.Context) {
  // 查询处于连接状态的子节点信息
  infos, err := database.QueryNodeInfoList()
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  var subNodeList [] SubNodeInfoBody
  for _, info := range infos {
    targetIp := strings.Split(info.Cidr, "/")[0]
    conn, err := v1Branch.ConnectGRPC(targetIp)
    if err != nil {
      continue
    }
    ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*600)
    client := pt.NewConnectivityClient(conn)
    replyConn, err := client.Ping(ctx, &pt.PingRequest{Client: "test.client"})
    cancel()
    if err != nil {
      conn.Close()
      continue
    }
    ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*3)
    clientModel := pt.NewModelClient(conn)
    replyModels, err := clientModel.ListStoredModel(ctx2, &pt.PingRequest{Client: "test.client"})
    cancel2()
    if err != nil {
      conn.Close()
      continue
    }
    subNode := SubNodeInfoBody{Version: replyConn.Version, Name: info.NodeName, Ip: targetIp, Models: replyModels.List}
    subNodeList = append(subNodeList, subNode)
    conn.Close()
  }
  response.Success(c, "query model list info success", subNodeList)
}
