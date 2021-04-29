package v1

import (
  "context"
  "io/ioutil"
  "net"
  "net/http"
  "strings"

  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  pt "gitlab.jiangxingai.com/jxserving/components/interfaces/protos"
  "gitlab.jiangxingai.com/jxserving/components/model/models/database"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
  "google.golang.org/grpc"
  log "k8s.io/klog"
)

/**
  check backend is alive
*/
func QueryActiveNodeHandler(c *gin.Context) {
  var body QueryListModelBody
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  conn, err := grpc.Dial(net.JoinHostPort(body.TargetIp, "50051"), grpc.WithInsecure())
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  defer conn.Close()
  client := pt.NewConnectivityClient(conn)
  reply, err := client.Ping(context.Background(), &pt.PingRequest{Client: "test.client"})
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  response.Success(c, "this ip is alive", reply.Version)
}

/**
  check backend is alive
*/
func QueryTargetNodeMachineInfoHandler(c *gin.Context) {
  var body QueryListModelBody
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  conn, err := grpc.Dial(net.JoinHostPort(body.TargetIp, "50051"), grpc.WithInsecure())
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  defer conn.Close()
  client := pt.NewConnectivityClient(conn)
  reply, err := client.ListNodeResources(context.Background(), &pt.PingRequest{Client: "test.client"})
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  response.Success(c, "query machine info success", reply)
}

/**
  query backend running model
*/
func QueryTargetNodeListModelHandler(c *gin.Context) {
  var body QueryListModelBody
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  conn, err := grpc.Dial(net.JoinHostPort(body.TargetIp, "50051"), grpc.WithInsecure())
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  defer conn.Close()
  client := pt.NewModelClient(conn)
  reply, err := client.ListStoredModel(context.Background(), &pt.PingRequest{Client: "test.client"})
  if err != nil {
    log.Error("get list storage model error:", err)
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  models := reply.List
  for _, model := range models {
    model := model.GetModel()
    log.Info("model info:", model)
  }
  response.Success(c, "query active node list success", models)
}

/**
  evaluate model
*/
func EvaluateModelHandler(c *gin.Context) {
  var body QueryListModelBody
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  conn, err := grpc.Dial(net.JoinHostPort(body.TargetIp, "50051"), grpc.WithInsecure())
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  defer conn.Close()
  //client := pt.NewModelClient(conn)
  //reply, err := client.EvaluateModel(context.Background(), &pt.EvaluateRequest{Model: "",Lable:"",Dataset:""})
  if err != nil {
    log.Error("get list storage model error:", err)
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  //leak:=reply.Leak
  //maps := reply.Map
  //miss := reply.Miss
  //acc := reply.Acc
  response.Success(c, "query model success", "")
}

/**
  在线分发model
*/
func DownloadModelOnLineHandler(c *gin.Context) {
  var body ModelAttributeBody
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  // 建立连接到gRPC服务
  conn, err := grpc.Dial(net.JoinHostPort(body.TargetIp, "50051"), grpc.WithInsecure())
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  // 函数结束时关闭连接
  defer conn.Close()
  client := pt.NewTransferFileClient(conn)
  file, err := ioutil.ReadFile(body.ModelPath)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  reply, err := client.Upload(context.Background(), &pt.Request{Buffer: file, Name: body.TargetIp})
  if err != nil {
    log.Error("grpc upload file error:", err)
    response.ServerError(c, http.StatusInternalServerError, "", err, reply)
    return
  }
  response.Success(c, "distribute model success", reply.Status)
}

/**
  下发模型
*/
func DistributeHandler(c *gin.Context) {
  var body interfaces.Distribute
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  // 查询数据库,获取下发服务器的ip
  err, data := database.QueryNodeInfo(body.Name)
  if data.Cidr == "" || err != nil {
    response.ServerError(c, http.StatusInternalServerError, "cidr is nil", nil, "")
    return
  }
  ip := strings.Split(data.Cidr, "/")[0]
  log.Info("wireguard ip: ", ip)
  // 建立连接到gRPC服务
  conn, err := grpc.Dial("10.50.1.106:50051", grpc.WithInsecure())
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  // 函数结束时关闭连接
  defer conn.Close()
  client := pt.NewTransferFileClient(conn)
  file, err := ioutil.ReadFile(body.FilePath)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  reply, err := client.Upload(context.Background(), &pt.Request{Buffer: file, Name: body.Name})
  if err != nil {
    log.Error("grpc upload file error:", err)
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  response.Success(c, "distribute model success", reply.Status)
}

/**
  create model
*/
func CreateAndLoadModelHandler(c *gin.Context) {
  var body CreateModelBody
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  conn, err := grpc.Dial(net.JoinHostPort(body.TargetIp, "50051"), grpc.WithInsecure())
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  defer conn.Close()
  client := pt.NewInferenceClient(conn)
  reply, err := client.CreateAndLoadModel(context.Background(), &pt.LoadRequest{Bid: body.Bid, Btype: body.Btype,
    Model: body.Model, Version: body.Version, Mode: body.Mode, Extra: body.Extra})
  if err != nil {
    log.Error("grpc upload file error:", err)
    response.ServerError(c, http.StatusInternalServerError, "", err, reply)
    return
  }
  response.Success(c, "create model success", reply)
}
