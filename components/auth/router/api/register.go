package api

import (
  "net/http"

  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "go.mongodb.org/mongo-driver/bson"

  "gitlab.jiangxingai.com/jxserving/components/portal/database"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
  log "k8s.io/klog"
)

type RegisterInfo struct {
  Username  string `json:"username"`
  Phone     string `json:"phone"`
  Pwd       string `json:"pwd"`
  SecondPwd string `json:"secondpwd"`
  Type      int64  `json:"type"`
}

func GetIdentity(body RegisterInfo) (identity string) {
  switch body.Type {
  case 0:
    identity = "user"
  case 1:
    identity = "admin"
  case 2:
    identity = "superadmin"
  }
  return
}

func Register(c *gin.Context) {
  var body RegisterInfo
  err := c.BindJSON(&body)
  if err != nil {
    log.Error(err)
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  } else {
    if body.Pwd != body.SecondPwd {
      response.ServerError(c, http.StatusInternalServerError, "pwd is different,please re-enter", nil, "")
      return
    }
    Database, _ := database.NewDatabase(database.TyMongo)
    err := Database.NewConnection()
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    Database.UpdateTable("auth")
    data, err := Database.Query(bson.M{"username": body.Username})
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    log.Info(len(data.([]bson.M)))
    if len(data.([]bson.M)) != 0 {
      log.Info("user exists")
      response.ServerError(c, http.StatusInternalServerError, "user exists", err, "")
      return
    }
    err = Database.Insert(&interfaces.UserInfo{
      Username:   body.Username,
      Pwd:        body.Pwd,
      Phone:      body.Phone,
      Permission: GetIdentity(body),
    })
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    response.Success(c, "register succ", map[string]interface{}{})
  }
}
