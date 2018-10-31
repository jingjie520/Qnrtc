package qnrtc

import (
	"core/utils/config"
	"core/utils/logutil"
	"time"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/rtc"
)

var (
	manager *rtc.Manager
	appID   string
)

func init() {
	//加载配置文件
	Conf := config.Conf
	if Conf != nil {
		//初始化七牛云
		appID = Conf.AppID

		mac := qbox.NewMac(Conf.AccessKey, Conf.SecretKey)
		manager = rtc.NewManager(mac)
	}
}

//GetToken 获取Token
func GetToken(roomName, userID string) string {
	token, err := manager.GetRoomToken(rtc.RoomAccess{AppID: appID, RoomName: roomName, UserID: userID, ExpireAt: time.Now().Unix() + 36000, Permission: "user"})

	if err == nil {
		return token
	}

	logutil.Error("获取TOKEN失败。", err)
	return "获取TOKEN失败。"
}
