package e2eTest

import (
  "bytes"
  "encoding/json"
  "io"
  "io/ioutil"
  "mime/multipart"
  "net/http"
  "net/http/httptest"
  "os"
  "testing"

  "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
)

var ResponseBody struct {
  Code    string      `json:"c"`
  Message string      `json:"msg"`
  Data    interface{} `json:"data"`
}

/**
  start test server and perform real request
*/
func PerformRequest(r http.Handler, method, path string, request io.Reader) *httptest.ResponseRecorder {
  req, _ := http.NewRequest(method, path, request)
  w := httptest.NewRecorder()
  w.Header().Set("Content-Type", "application/json")
  r.ServeHTTP(w, req)
  return w
}

/**
  generate request body and router for post
*/
func GeneratePostRequestBody(body interface{}) (*gin.Engine, io.Reader, error) {
  b, err := json.Marshal(body)
  reader := bytes.NewReader(b)
  engine := gin.Default()
  return engine, reader, err
}

/**
  generate request body and router for get
*/
func GenerateGetRequestBody() *gin.Engine {
  engine := gin.Default()
  return engine
}

/**
  response assert ok
*/
func AssertResponseOk(t *testing.T, w *httptest.ResponseRecorder, err error) (interface{}, string) {
  assert.Equal(t, http.StatusOK, w.Code)
  err = json.Unmarshal([]byte(w.Body.String()), &ResponseBody)
  msg := ResponseBody.Message
  data := ResponseBody.Data
  assert.Nil(t, err)
  assert.NotNil(t, data)
  return data, msg
}

/**
  response assert error
*/
func AssertResponseError(t *testing.T, w *httptest.ResponseRecorder, dstMsg string, err error) {
  err = json.Unmarshal([]byte(w.Body.String()), &ResponseBody)
  msg := ResponseBody.Message
  assert.Nil(t, err)
  assert.Equal(t, dstMsg, msg)
}

/**
  generate multipart request
*/
func NewFileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  fileContents, err := ioutil.ReadAll(file)
  if err != nil {
    return nil, err
  }
  fi, err := file.Stat()
  if err != nil {
    return nil, err
  }
  defer file.Close()

  body := new(bytes.Buffer)
  writer := multipart.NewWriter(body)
  part, err := writer.CreateFormFile(paramName, fi.Name())
  if err != nil {
    return nil, err
  }
  _, err = part.Write(fileContents)
  if err != nil {
    return nil, err
  }

  for key, val := range params {
    _ = writer.WriteField(key, val)
  }
  err = writer.Close()
  if err != nil {
    return nil, err
  }
  request, err := http.NewRequest("POST", uri, body)
  if err != nil {
    return nil, err
  }
  request.Header.Add("Content-Type", writer.FormDataContentType())
  return request, err
}
