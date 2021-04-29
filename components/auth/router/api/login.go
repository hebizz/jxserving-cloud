package api

import (
  "net/http"
  "time"

  "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/portal/database"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
  "go.mongodb.org/mongo-driver/bson"
  log "k8s.io/klog"
)

type LoginInfo struct {
  Username string `json:"username"`
  Pwd      string `json:"pwd"`
  Permission string `json:"permission"`
  jwt.StandardClaims
}

type JWT struct {
  SigningKey []byte
}

func (j *JWT) CreateToken(claims LoginInfo) (string, error) {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString(j.SigningKey)
}

func generateToken(user LoginInfo) (string, error) {
  j := &JWT{
    []byte("jxs"),
  }
  claims := LoginInfo{
    user.Username,
    user.Pwd,
    user.Permission,
    jwt.StandardClaims{
      NotBefore: int64(time.Now().Unix()),
      ExpiresAt: int64(time.Now().Unix() + 100000),
      Issuer:    "jxs",
    },
  }
  token, err := j.CreateToken(claims)
  if err != nil {
    log.Error(err)
    return "", err
  }
  log.Info(token)
  return token, nil

}

func Login(c *gin.Context) {
  var body LoginInfo
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  } else {
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
    if data.([]bson.M)[0]["username"] == body.Username && data.([]bson.M)[0]["pwd"] == body.Pwd {
      log.Info("pwd correct")
      mid, _ := data.([]bson.M)[0]["permission"]
      log.Info(mid)
      body.Permission = mid.(string)
      token, err := generateToken(body)
      if err != nil {
        response.ServerError(c, http.StatusInternalServerError, "", err, "")
        return
      }
      response.Success(c, "login succ", map[string]interface{}{"token": token})
    } else {
      log.Info("pwd incorrect")
      response.ServerError(c, http.StatusInternalServerError, "pwd incorrect", nil, "")
    }
    }
}
