package v1Branch

import (
  "context"
  "net"
  "time"

  pt "gitlab.jiangxingai.com/jxserving/components/interfaces/protos"
  "google.golang.org/grpc"
  log "k8s.io/klog"
)

func ConnectGRPC(targetIp string) (*grpc.ClientConn,error)  {
  conn, err := grpc.Dial(net.JoinHostPort(targetIp, "50051"), grpc.WithInsecure())
  if err != nil {
    return nil, err
  }
  return conn,nil
}

func JxServingGRPC(conn *grpc.ClientConn) (*pt.PingReply, *pt.ResourcesReply, *pt.ListReply, error) {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()
  client := pt.NewConnectivityClient(conn)
  replyConn, err := client.Ping(ctx, &pt.PingRequest{Client: "test.client"})
  if err != nil {
    return nil, nil, nil, err
  }
  replyNode, err := client.ListNodeResources(ctx, &pt.PingRequest{Client: "test.client"})
  if err != nil {
    return nil, nil, nil, err
  }
  clientModel := pt.NewModelClient(conn)
  reply, err := clientModel.ListStoredModel(ctx, &pt.PingRequest{Client: "test.client"})
  if err != nil {
    log.Error("get list storage model error:", err)
    return nil, nil, nil, err
  }
  return replyConn, replyNode, reply, nil
}