package message

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"coding.net/cherrysd/wxserver/util"
)

type responseVoiceMessage struct {
	PublicMessage
	MediaID string `xml:"Voice>MediaId"`
}

// Voice 回复语音逻辑消息体(供外部使用)
type Voice struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MediaID      string
	MsgID        int64
}

func (rtmsg *Voice) formatLogicMsg() (string, error) {
	destMsg := responseVoiceMessage{}
	destMsg.MediaID = rtmsg.MediaID
	destMsg.ToUserName = rtmsg.ToUserName
	destMsg.FromUserName = rtmsg.FromUserName
	if rtmsg.CreateTime == 0 {
		destMsg.CreateTime = util.GetCurrTimeStamp()
	} else {
		destMsg.CreateTime = rtmsg.CreateTime
	}
	destMsg.MsgType = VoiceMsg

	responseRawXMLMsg, err := xml.Marshal(destMsg)
	if err != nil {
		log.Println("Build Response Text Message Error")
		return "", err
	}
	result := string(responseRawXMLMsg)
	return result, err
}

// Send 向服务器发送文字消息
func (rtmsg *Voice) Send(w http.ResponseWriter) error {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = xmlContentType
	}
	w.WriteHeader(200)
	strResponseMsg, err := rtmsg.formatLogicMsg()
	fmt.Fprint(w, strResponseMsg)
	return err
}
