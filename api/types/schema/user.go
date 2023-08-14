package schema

type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"userName"`
	Realname string `json:"realName"`
	Rolename string `json:"roleName"`
}

// 登录用户结构
type LoginUser struct {
	Username string `form:"Username" json:"userName" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
}
