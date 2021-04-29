package middleware

import (
  "errors"
  "net/http"

  "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
  log "k8s.io/klog"
)

type JWT struct {
  SigningKey []byte
}

type CustomClaims struct {
  Username string `json:"username"`
  Pwd      string `json:"pwd"`
  Permission string `json:"permission"`
  jwt.StandardClaims
}

var (
  TokenExpired     error  = errors.New("Token is expired")
  TokenNotValidYet error  = errors.New("Token not active yet")
  TokenMalformed   error  = errors.New("That's not even a token")
  TokenInvalid     error  = errors.New("Couldn't handle this token:")
  SignKey          string = "jxs"
)

func NewJWT() *JWT {
  return &JWT{
    []byte(GetSignKey()),
  }
}

func GetSignKey() string {
  return SignKey
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
  token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
    return j.SigningKey, nil
  })
  if err != nil {
   if ve, ok := err.(*jwt.ValidationError); ok {
     if ve.Errors&jwt.ValidationErrorMalformed != 0 {
       return nil, TokenMalformed
     } else if ve.Errors&jwt.ValidationErrorExpired != 0 {
       return nil, TokenExpired
     } else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
       return nil, TokenNotValidYet
     } else {
       return nil, TokenInvalid
     }
   }
  }
  if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
    log.Info(claims)
    return claims, nil
  } else {
    return nil, TokenInvalid
  }
}

func Jwt() gin.HandlerFunc {
  return func(c *gin.Context) {
    token := c.Request.Header.Get("token")
    if token == "" {
      response.ServerError(c, http.StatusInternalServerError, "no token", nil, "")
      c.Abort()
      return
    }
    log.Info("get token: ", token)
    j := NewJWT()
    claims, err := j.ParseToken(token)
    if err != nil {
      if err == TokenExpired {
        response.ServerError(c, http.StatusInternalServerError, "token expired", nil, "")
        c.Abort()
        return
      }
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      c.Abort()
      return
    }
    c.Set("claims", claims)
  }
}
