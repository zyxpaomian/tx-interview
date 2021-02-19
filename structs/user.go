package structs

type User struct {
	Id       int64  `json:"id"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

type Token struct {
	AgentVersion string `json:"version"`
	UpdateTime   string `json:"updatetime"`
}
