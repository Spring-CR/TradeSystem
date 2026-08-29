package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type WeChatMsgColor string
const (
	ColorWarning WeChatMsgColor = "warning"
	ColorInfo WeChatMsgColor = "info"
)

// 企业微信消息结构体
type WeChatMessage struct {
	Msgtype  string `json:"msgtype"`
	Text struct {
		Content string `json:"content"`
		MentionedList []string `json:"mentioned_list"`
	} `json:"text"`
}

// 发送消息到企业微信机器人
func SendToWeChat(url, content string) error {
	msg := WeChatMessage{
		Msgtype: "text",
		Text: struct {
			Content string `json:"content"`
			MentionedList []string `json:"mentioned_list"`
		}{
			Content: content,
			MentionedList: []string{"@all"},
		},
	}

	// 将消息结构体转为JSON
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("构建消息JSON失败: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(msgJSON))
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("发送失败，响应状态: %s, 读取响应内容失败: %v", resp.Status, err)
		}
		return fmt.Errorf("发送失败，响应状态: %s, 响应内容: %s", resp.Status, string(body))
	}

	return nil
}

func GenWeChatMessage(passed bool, headTpl string, info string, err error) string {
	if err != nil {
		return fmt.Sprintf(headTpl, ColorWarning, time.Now().Format("2006-01-02")) + fmt.Sprintf("任务执行异常, error: %v\n", err)
	}
	if passed {
		return fmt.Sprintf(headTpl, ColorInfo, time.Now().Format("2006-01-02")) + info
	}
	return fmt.Sprintf(headTpl, ColorWarning, time.Now().Format("2006-01-02")) + info
}