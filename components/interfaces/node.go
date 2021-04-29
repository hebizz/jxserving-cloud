package interfaces


type NodeInfo struct {
  Id        string `bson:"_id" json:"id"`
  NodeName  string `json:"name" binding:"required" `
  Cidr      string `json:"cidr" binding:"required" `
  PublicKey string `json:"publicKey" binding:"required" `
  Status    string `json:"status"`
  NetSpeed  string `json:"netSpeed"`
}

type NodeList struct {
  Master   NodeInfo   `json:"master"`
  NodeList []NodeInfo `json:"nodeList"`
}

type NodeDelete struct {
  Id string `json:"id" binding:"required"`
}
