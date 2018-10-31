package main

import (
	"core/service/qnrtc"
	"core/utils/config"
	"core/utils/logutil"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/qiniu/api.v7/rtc"
)

var (
	manager *rtc.Manager
	appID   string
)

func init() {

	//加载配置文件
	Conf := config.Conf
	if Conf == nil {
		fmt.Println("加载配置文件conf.ini失败。请检查当前目录下是否存在该文件。")
		os.Exit(-1)
	}
}

func main() {
	//注册处理函数，用户连接，自动调用指定的处理函数

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)

	router.HandleFunc("/api/v1.qnrtc/getToken/roomName/{roomName}/userId/{userId}", GetToken)

	//监听绑定
	http.ListenAndServe(":80", router)
}

//Index 默认首页
func Index(writer http.ResponseWriter, request *http.Request) {
	logutil.RequestInfo(request)
	fmt.Fprintln(writer, "Welcome Go! ")
}

//GetToken 获取Token
func GetToken(writer http.ResponseWriter, request *http.Request) {
	logutil.RequestInfo(request)

	vars := mux.Vars(request)
	roomName := vars["roomName"]
	userID := vars["userId"]
	writer.Write([]byte(qnrtc.GetToken(roomName, userID)))
}
