package message

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"coding.net/cherrysd/wxserver/util"
)

type responseVideoMessage struct {
	PublicMessage
	MediaID     string `xml:"Video>MediaId"`
	Title       string `xml:"Video>Title"`
	Description string `xml:"Video>Description"`
}

// Video 回复视频逻辑消息体(供外部使用)
type Video struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MediaID      string
	Title        string
	Description  string
	MsgID        int64
}

// ShortVideo 小视频消息体
type ShortVideo struct {
}

func (rtmsg *Video) formatLogicMsg() (string, error) {
	destMsg := responseVideoMessage{}
	destMsg.MediaID = rtmsg.MediaID
	destMsg.ToUserName = rtmsg.ToUserName
	destMsg.FromUserName = rtmsg.FromUserName
	destMsg.Title = rtmsg.Title
	destMsg.Description = rtmsg.Description
	if rtmsg.CreateTime == 0 {
		destMsg.CreateTime = util.GetCurrTimeStamp()
	} else {
		destMsg.CreateTime = rtmsg.CreateTime
	}
	destMsg.MsgType = VideoMsg

	responseRawXMLMsg, err := xml.Marshal(destMsg)
	if err != nil {
		log.Println("Build Response Text Message Error")
		return "", err
	}
	result := string(responseRawXMLMsg)
	return result, err
}

// Send 向服务器发送文字消息
func (rtmsg *Video) Send(w http.ResponseWriter) error {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = xmlContentType
	}
	w.WriteHeader(200)
	strResponseMsg, err := rtmsg.formatLogicMsg()
	fmt.Fprint(w, strResponseMsg)
	return err
}
