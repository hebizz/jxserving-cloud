package interfaces

type UserInfo struct {
  Username    string `json:"username"`
  Pwd         string `json:"pwd"`
  Phone       string `json:"phone"`
  Permission  string `json:"permission"`
}
