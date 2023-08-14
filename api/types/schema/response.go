package schema

// 返回token信息
type Token struct {
	Token   string `json:"token"`
	Expired string `json:"expired"`
}
