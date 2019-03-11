package menu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"coding.net/cherrysd/wxserver/server"
)

// Type 菜单类型
type Type string

const jsonContentType = "application/json;charset=utf-8"
const createMenuURL = "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s"

// 菜单类型枚举
const (
	ClickType             = "click"
	ViewType              = "view"
	ScancodePushType      = "scancode_push"
	ScancodeWatingMsgType = "scancode_waitmsg"
	PicSysPhotoType       = "pic_sysphoto"
	PicPhotoAlbumType     = "pic_photo_or_album"
	PicWXType             = "pic_weixin"
	LocationSelectType    = "location_select"
	MediaIDType           = "media_id"
	ViewLimitType         = "view_limited"
	MiniProgramType       = "miniprogram"
)

type jsonButtons struct {
	Button []jsonButton `json:"button"`
}

type jsonButton struct {
	Name       string          `json:"name,omitempty"`
	ButtonType Type            `json:"type,omitempty"`
	Key        string          `json:"key,omitempty"`
	URL        string          `json:"url,omitempty"`
	MediaID    string          `json:"media_id,omitempty"`
	AppID      string          `json:"appid,omitempty"`
	PagePath   string          `json:"pagepath,omitempty"`
	SubButton  []jsonSubButton `json:"sub_button,omitempty"`
}

type jsonSubButton struct {
	Name       string `json:"name,omitempty"`
	ButtonType Type   `json:"type,omitempty"`
	Key        string `json:"key,omitempty"`
	URL        string `json:"url,omitempty"`
	MediaID    string `json:"media_id,omitempty"`
	AppID      string `json:"appid,omitempty"`
	PagePath   string `json:"pagepath,omitempty"`
}

// LevelButton 一级按钮(弹出二级菜单)
type LevelButton struct {
	Name      string
	SubButton []FunctionButton
}

// FunctionButton 功能按钮
type FunctionButton struct {
	Name       string
	ButtonType Type
	Key        string
	URL        string
	MediaID    string
	AppID      string
	PagePath   string
}

type innerButton struct {
	Name       string
	ButtonType Type
	Key        string
	URL        string
	MediaID    string
	AppID      string
	PagePath   string
	SubButton  []innerSubButton
}

type innerSubButton struct {
	Name       string
	ButtonType Type
	Key        string
	URL        string
	MediaID    string
	AppID      string
	PagePath   string
}

// MainMenu 菜单实例
type MainMenu struct {
	buttons  []innerButton
	dbServer *server.Server
}

// Error 菜单接口调用结果JSON
type Error struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// CreateMenu 创建菜单
func CreateMenu(menu *MainMenu) bool {
	jsonContent, err := menu.formatJSON()
	if err != nil {
		return false
	}

	body := bytes.NewBuffer(jsonContent)
	menuURL := fmt.Sprintf(createMenuURL, menu.dbServer.GetAccessToken())
	response, postErr := http.Post(menuURL, jsonContentType, body)
	if postErr != nil {
		return false
	}

	result, responseError := ioutil.ReadAll(response.Body)
	if responseError != nil {
		return false
	}

	resultJSON := Error{}
	jsonError := json.Unmarshal(result, &resultJSON)
	if jsonError != nil {
		return false
	}

	if resultJSON.ErrCode != 0 {
		fmt.Println(resultJSON.ErrMsg)
		return false
	}

	return true
}

// CreateMenuWithToken 传入token创建菜单
func CreateMenuWithToken(menu *MainMenu, accessToken string) bool {
	jsonContent, err := menu.formatJSON()
	if err != nil {
		return false
	}

	body := bytes.NewBuffer(jsonContent)
	menuURL := fmt.Sprintf(createMenuURL, accessToken)
	response, postErr := http.Post(menuURL, jsonContentType, body)
	if postErr != nil {
		return false
	}

	result, responseError := ioutil.ReadAll(response.Body)
	if responseError != nil {
		return false
	}

	resultJSON := Error{}
	jsonError := json.Unmarshal(result, &resultJSON)
	if jsonError != nil {
		return false
	}

	if resultJSON.ErrCode != 0 {
		fmt.Println(resultJSON.ErrMsg)
		return false
	}

	return true
}

// NewMenu 新建菜单实例
func NewMenu(svr *server.Server) *MainMenu {
	newSvr := new(MainMenu)
	newSvr.dbServer = svr
	return newSvr
}

// AddFirstLevelButton 添加一级菜单按钮
func (mm *MainMenu) AddFirstLevelButton(btn *LevelButton) bool {
	if len(btn.SubButton) > 5 {
		return false
	}
	innerBtn := innerButton{}
	innerBtn.Name = btn.Name
	for index := 0; index < len(btn.SubButton); index++ {
		subBtn := innerSubButton{}
		subBtn.Name = btn.SubButton[index].Name
		subBtn.ButtonType = btn.SubButton[index].ButtonType
		subBtn.Key = btn.SubButton[index].Key
		subBtn.URL = btn.SubButton[index].URL
		subBtn.MediaID = btn.SubButton[index].MediaID
		subBtn.AppID = btn.SubButton[index].AppID
		subBtn.PagePath = btn.SubButton[index].PagePath
		innerBtn.SubButton = append(innerBtn.SubButton, subBtn)
	}
	mm.buttons = append(mm.buttons, innerBtn)
	return true
}

// AddFunctionButton 添加功能按钮
func (mm *MainMenu) AddFunctionButton(btn *FunctionButton) bool {
	if len(mm.buttons) > 2 {
		return false
	}

	innerBtn := innerButton{}
	innerBtn.Name = btn.Name
	innerBtn.ButtonType = btn.ButtonType
	innerBtn.Key = btn.Key
	innerBtn.URL = btn.URL
	innerBtn.MediaID = btn.MediaID
	innerBtn.AppID = btn.AppID
	innerBtn.PagePath = btn.PagePath

	mm.buttons = append(mm.buttons, innerBtn)
	return true
}

// AddFunctionButton 添加功能按钮
func (menuBtn *LevelButton) AddFunctionButton(btn *FunctionButton) bool {
	if len(menuBtn.SubButton) > 4 {
		return false
	}
	menuBtn.SubButton = append(menuBtn.SubButton, *btn)
	return true
}

func (mm *MainMenu) formatJSON() ([]byte, error) {
	jsonButtons := mm.getJSONButtons()
	return json.Marshal(jsonButtons)
}

func (mm *MainMenu) getJSONButtons() jsonButtons {
	jsonBtns := jsonButtons{}
	for index := 0; index < len(mm.buttons); index++ {
		jsonBtn := jsonButton{}
		jsonBtn.Name = mm.buttons[index].Name
		if len(mm.buttons[index].SubButton) > 0 {
			//一级菜单
			for i := 0; i < len(mm.buttons[index].SubButton); i++ {
				jsonSubBtn := jsonSubButton{}
				jsonSubBtn.Name = mm.buttons[index].SubButton[i].Name
				jsonSubBtn.ButtonType = mm.buttons[index].SubButton[i].ButtonType
				jsonSubBtn.Key = mm.buttons[index].SubButton[i].Key
				jsonSubBtn.URL = mm.buttons[index].SubButton[i].URL
				jsonSubBtn.MediaID = mm.buttons[index].SubButton[i].MediaID
				jsonSubBtn.AppID = mm.buttons[index].SubButton[i].AppID
				jsonSubBtn.PagePath = mm.buttons[index].SubButton[i].PagePath
				jsonBtn.SubButton = append(jsonBtn.SubButton, jsonSubBtn)
			}
		} else {
			jsonBtn.ButtonType = mm.buttons[index].ButtonType
			jsonBtn.Key = mm.buttons[index].Key
			jsonBtn.URL = mm.buttons[index].URL
			jsonBtn.MediaID = mm.buttons[index].MediaID
			jsonBtn.AppID = mm.buttons[index].AppID
			jsonBtn.PagePath = mm.buttons[index].PagePath
		}
		jsonBtns.Button = append(jsonBtns.Button, jsonBtn)
	}

	return jsonBtns
}
