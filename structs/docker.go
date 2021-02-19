package structs

type ContainerInfo struct {
	ContainerID string `json:"containerid"`
	UserName    string `json:"username"`
	PassWord    string `json:"password"`
	Port        string `json:"port"`
}
