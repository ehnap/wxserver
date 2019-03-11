package message

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"coding.net/cherrysd/wxserver/util"
)

type responseTextMessage struct {
	PublicMessage
	Content string `xml:"Content"`
}

// Text 回复文本逻辑消息体(供外部使用)
type Text struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	Content      string
	MsgID        int64
}

func (rtmsg *Text) formatLogicMsg() (string, error) {
	destMsg := responseTextMessage{}
	destMsg.Content = rtmsg.Content
	destMsg.ToUserName = rtmsg.ToUserName
	destMsg.FromUserName = rtmsg.FromUserName
	if rtmsg.CreateTime == 0 {
		destMsg.CreateTime = util.GetCurrTimeStamp()
	} else {
		destMsg.CreateTime = rtmsg.CreateTime
	}
	destMsg.MsgType = TextMsg

	responseRawXMLMsg, err := xml.Marshal(destMsg)
	if err != nil {
		log.Println("Build Response Text Message Error")
		return "", err
	}
	result := string(responseRawXMLMsg)
	return result, err
}

// Send 向服务器发送文字消息
func (rtmsg *Text) Send(w http.ResponseWriter) error {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = xmlContentType
	}
	w.WriteHeader(200)
	strResponseMsg, err := rtmsg.formatLogicMsg()
	fmt.Fprint(w, strResponseMsg)
	return err
}
