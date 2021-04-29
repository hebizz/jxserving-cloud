package middleware

import (
  "net/http"

  "github.com/casbin/casbin"
  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
  log "k8s.io/klog"
)

func Interceptor(e *casbin.Enforcer) gin.HandlerFunc {
  return func(c *gin.Context) {
    info,_ := c.Get("claims")
    log.Info(info.(*CustomClaims).Permission)
    sub := info.(*CustomClaims).Permission
    obj := c.Request.URL.RequestURI()
    act := c.Request.Method
    res, err := e.Enforce(sub, obj, act)
    if err != nil {
      log.Error(err)
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      c.Abort()
      return
    }
    if res == true {
      log.Info("auth access")
      c.Next()
      return
    } else {
      log.Error("auth fail")
      response.ServerError(c, http.StatusInternalServerError, "auth fail",nil, "")
      c.Abort()
      return
    }
  }
}

