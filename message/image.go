package message

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"coding.net/cherrysd/wxserver/util"
)

type responseImageMessage struct {
	PublicMessage
	MediaID string `xml:"Image>MediaId"`
}

// Image 回复图片逻辑消息体(供外部使用)
type Image struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MediaID      string
	MsgID        int64
}

func (rtmsg *Image) formatLogicMsg() (string, error) {
	destMsg := responseImageMessage{}
	destMsg.MediaID = rtmsg.MediaID
	destMsg.ToUserName = rtmsg.ToUserName
	destMsg.FromUserName = rtmsg.FromUserName
	if rtmsg.CreateTime == 0 {
		destMsg.CreateTime = util.GetCurrTimeStamp()
	} else {
		destMsg.CreateTime = rtmsg.CreateTime
	}
	destMsg.MsgType = ImageMsg

	responseRawXMLMsg, err := xml.Marshal(destMsg)
	if err != nil {
		log.Println("Build Response Text Message Error")
		return "", err
	}
	result := string(responseRawXMLMsg)
	return result, err
}

// Send 向服务器发送图片消息
func (rtmsg *Image) Send(w http.ResponseWriter) error {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = xmlContentType
	}
	w.WriteHeader(200)
	strResponseMsg, err := rtmsg.formatLogicMsg()
	fmt.Fprint(w, strResponseMsg)
	return err
}
