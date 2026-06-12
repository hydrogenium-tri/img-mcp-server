package main

//本文件中描述了图片格式分析和模型的调用

import (
	"context"
	"fmt"
	"encoding/base64"

	"github.com/sashabaranov/go-openai"
)

//图像解码函数
func detectImageFormat(base64Str string) string{
	//Base64解码
	decoded, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil{
		return "image/png"
	}

	//检查文件格式
	switch{
		case decoded[0] == 0xFF:
		return "image/jpeg"
		case decoded[0] == 0x89:
		return "image/png"
		case decoded[0] == 0x47:
		return "image/gif"
		case decoded[0] == 0x52:
		return "image/webp"
		case decoded[0] == 0x42:
		return "image/bmp"
		default:
		return "image/png"
	}
	
}

//创建OpenAI客户端
func callVisionModel(ctx context.Context, base64Image string, prompt string) (string, error){
	//创建OpenaiAPI的客户端
	clientConfig := openai.DefaultConfig(config.APIKey)
	clientConfig.BaseURL = config.BaseURL
	client := openai.NewClientWithConfig(clientConfig)

	//检测图片格式
	imageFormat := detectImageFormat(base64Image)

	//构建消息内容
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: config.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleUser,
				MultiContent: []openai.ChatMessagePart{
					{
						Type: openai.ChatMessagePartTypeText,
						Text: prompt,
					},
					{
						Type: openai.ChatMessagePartTypeImageURL,
						ImageURL: &openai.ChatMessageImageURL{
							URL: "data:"+imageFormat+";base64," + base64Image,
						},
					},
				},
			},
		},
		MaxTokens: 1000,	//限制返回长度
	})
	
 	if err != nil {
        return "", err
    }
    if len(resp.Choices) > 0 {
        return resp.Choices[0].Message.Content, nil
    }
    return "", fmt.Errorf("模型未返回有效结果")
}