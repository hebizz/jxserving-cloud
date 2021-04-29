package database

import (
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

/**
  查询所有nodeList
*/
func QueryNodeInfoList() ([]interfaces.NodeInfo, error) {
  db.UpdateTable("node")
  var dataList [] interfaces.NodeInfo
  cursor, ctx, err := db.QueryAll(bson.M{}, nil)
  if err != nil {
    return dataList, err
  }
  defer cursor.Close(ctx)
  err = cursor.All(ctx, &dataList)
  return dataList, err
}
/**
   查询 cidr 是否已经存在
 */
func QueryNodeByCidr(cidr string) (int64,error) {
  db.UpdateTable("node")
  selector:=bson.M{"cidr":cidr}
  return db.QueryCount(selector)
}


/**
  插入nodeInfo数据
*/
func InsertNodeInfo(data interfaces.NodeInfo) error {
  db.UpdateTable("node")
  data.Id = primitive.NewObjectID().Hex()
  return db.Insert(data)
}

/**
  删除node节点
*/
func DeleteNode(id string) error {
  db.UpdateTable("node")
  selector := bson.M{"_id": id}
  return db.Delete(selector)
}

/**
  更新node节点信息
*/
func UpdateNodeInfo(data interfaces.NodeInfo) error {
  db.UpdateTable("node")
  selector := bson.M{"_id": data.Id}
  updateData := bson.M{"$set": bson.M{"nodename": data.NodeName, "cidr": data.Cidr, "publickey": data.PublicKey}}
  return db.Update(selector, updateData)
}
