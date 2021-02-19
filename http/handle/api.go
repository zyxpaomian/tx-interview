package handle

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"tx-interview/common"
	log "tx-interview/common/formatlog"
	"tx-interview/controller"
)

// 用户认证
func apiUserAuth(res http.ResponseWriter, req *http.Request) {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.Errorf("[http] 请求报文解析失败")
		common.ReqBodyInvalid(res)
		return
	}

	request := &Request{}
	if err := common.ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("[http] 解析模板JSON失败")
		common.ResMsg(res, 400, err.Error())
		return
	}

	type Response struct {
		Token string `json:"token"`
	}

	token, err := controller.UserController.GetToken(request.Username, request.Password)
	if err != nil {
		log.Errorln("[http] 用户认证失败")
		common.ResMsg(res, 400, err.Error())
		return
	}

	response := &Response{Token: token}
	result, err := json.Marshal(response)
	if err != nil {
		log.Errorf("[http] apiUserAuth JSON生成失败, %v", err.Error())
		common.ResMsg(res, 500, err.Error())
		return
	}
	common.ResMsg(res, 200, string(result))
}

// 获取所有用户名
func apiGetAllUser(res http.ResponseWriter, req *http.Request) {
	users, err := controller.UserController.GetAllUsers()
	if err != nil {
		log.Errorf("[http] apiGetAllUser 数据处理失败, %v", err.Error())
		common.ResMsg(res, 500, err.Error())
		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		log.Errorf("[http] apiGetAllAgents JSON生成失败, %v", err.Error())
		common.ResMsg(res, 400, err.Error())
		return
	}
	common.ResMsg(res, 200, string(response))
}

// 获取所有的docker image
func apiGetAllImage(res http.ResponseWriter, req *http.Request) {
	images, err := controller.DockerController.ListImage()
	if err != nil {
		log.Errorf("[http] apiGetAllImage 数据处理失败, %v", err.Error())
		common.ResMsg(res, 500, err.Error())
		return
	}

	type Response struct {
		ImageList []string `json:"imagelist"`
	}

	response := &Response{ImageList: images}
	result, err := json.Marshal(response)
	if err != nil {
		log.Errorf("[http] apiGetAllImage JSON生成失败, %v", err.Error())
		common.ResMsg(res, 500, err.Error())
		return
	}
	common.ResMsg(res, 200, string(result))
}

// 启动容器
func apiRunContainer(res http.ResponseWriter, req *http.Request) {
	type Request struct {
		ContainerType string `json:"containertype"`
	}

	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.Errorf("[http] 请求报文解析失败")
		common.ReqBodyInvalid(res)
		return
	}

	request := &Request{}
	if err := common.ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("[http] 解析模板JSON失败")
		common.ResMsg(res, 400, err.Error())
		return
	}

	containerResult, err := controller.DockerController.CreateContainer(request.ContainerType)
	if err != nil {
		log.Errorf("[http] apiRunContainer 数据处理失败, %v", err.Error())
		common.ResMsg(res, 500, err.Error())
		return
	}

	response, err := json.Marshal(containerResult)
	if err != nil {
		log.Errorf("[http] apiRunContainer JSON生成失败, %v", err.Error())
		common.ResMsg(res, 400, err.Error())
		return
	}
	common.ResMsg(res, 200, string(response))
}
