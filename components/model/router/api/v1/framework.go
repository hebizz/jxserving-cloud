package v1

import (
  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
)

func GetFramework(c *gin.Context) {
  Framework := []string{"tensorflow", "tensorflow-lite", "tensorflow-serving", "rknn", "pytorch"}
  response.Success(c, "query framework success", map[string]interface{}{"result": Framework})
}
