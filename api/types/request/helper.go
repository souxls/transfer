package request

type RequestUsers struct {
	Data []string `json:"data"`
}

type LoginInfo struct {
	Username string `form:"Username" json:"userName"`
}
