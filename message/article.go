package message

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"coding.net/cherrysd/wxserver/util"
)

type responseArticleMessage struct {
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
	PicURL      string `xml:"PicUrl"`
	URL         string `xml:"Url"`
}

type responseArticlesMessage struct {
	PublicMessage
	ArticleCount int                      `xml:"ArticleCount"`
	Content      []responseArticleMessage `xml:"Articles>item"`
}

// Article 回复图文逻辑子项消息体
type Article struct {
	Title       string
	Description string
	PicURL      string
	URL         string
}

// Articles 回复图文逻辑消息体(供外部使用)
type Articles struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	Content      []Article
}

func (rtmsg *Articles) formatLogicMsg() (string, error) {
	destMsg := responseArticlesMessage{}
	destMsg.ToUserName = rtmsg.ToUserName
	destMsg.FromUserName = rtmsg.FromUserName
	for index := 0; index < len(rtmsg.Content); index++ {
		// 超过8条就截断
		if index >= 8 {
			break
		}

		destData := responseArticleMessage{}
		destData.Title = rtmsg.Content[index].Title
		destData.PicURL = rtmsg.Content[index].PicURL
		destData.URL = rtmsg.Content[index].URL
		destData.Description = rtmsg.Content[index].Description
		destMsg.Content = append(destMsg.Content, destData)
		destMsg.ArticleCount = index + 1
	}

	if rtmsg.CreateTime == 0 {
		destMsg.CreateTime = util.GetCurrTimeStamp()
	} else {
		destMsg.CreateTime = rtmsg.CreateTime
	}
	destMsg.MsgType = ArticleMsg

	responseRawXMLMsg, err := xml.Marshal(destMsg)
	if err != nil {
		log.Println("Build Response Text Message Error")
		return "", err
	}
	result := string(responseRawXMLMsg)
	return result, err
}

// Send 向服务器发送文字消息
func (rtmsg *Articles) Send(w http.ResponseWriter) error {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = xmlContentType
	}
	w.WriteHeader(200)
	strResponseMsg, err := rtmsg.formatLogicMsg()
	fmt.Fprint(w, strResponseMsg)
	return err
}
