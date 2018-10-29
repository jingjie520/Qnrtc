package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/rtc"
)

var (
	manager *rtc.Manager
	appID   string
)

func init() {
	accessKey := "J3HMHHvYeUaYYBoYdegtP9ZKMC_k1H--RQp2fWuG"
	secretKey := ""
	appID = "drnocxlxj"
	mac := qbox.NewMac(accessKey, secretKey)
	manager = rtc.NewManager(mac)
}

func main() {
	//注册处理函数，用户连接，自动调用指定的处理函数

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)

	router.HandleFunc("/api/v1.qnrtc/getToken/roomName/{roomName}/userId/{userId}", GetToken)

	//监听绑定
	http.ListenAndServe(":80", router)
}

func Index(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Welcome Go! ")
}

func GetToken(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	roomName := vars["roomName"]
	userID := vars["userId"]

	token, err := getRoomToken(appID, roomName, userID)

	if err == nil {
		writer.Write([]byte(token))
	} else {
		fmt.Fprintln(writer, "获取TOKEN失败。", err)
	}

}

func getRoomToken(appId, roomName, userID string) (token string, err error) {
	token, err = manager.GetRoomToken(rtc.RoomAccess{AppID: appId, RoomName: roomName, UserID: userID, ExpireAt: time.Now().Unix() + 36000, Permission: "user"})
	return
}
