package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sashabaranov/go-openai"
)

var config Config

func main() {
	//解析配置文件
	config = loadConfig(Config{})
	
	
	//创建服务器实例
	s := server.NewMCPServer(
		"image-recognition",
		"0.0.1",
	)

	//新建一个工具
	analyzeImage := mcp.NewTool("analyze_image",
		mcp.WithDescription("传入图像字符串，使用一个多模态模型解析它"),
		mcp.WithString("base64",
   			mcp.Required(),
      		mcp.Description("字符串，base64编码的图像数据"),
		),
		mcp.WithString("prompt",
   			mcp.Required(),
      		mcp.Description("字符串，对图像模型的提示词"),
		),
	)

	//注册这个工具
	s.AddTool(analyzeImage, handleAnalyzeImage)

	//验证配置文件加载
	fmt.Printf("配置加载成功，模型：%s\n", config.Model)
	
	//启动函数
	addr := fmt.Sprintf(":%d", config.Port)
	sseServer := server.NewStreamableHTTPServer(s,
		server.WithStateLess(true),
	)
	if err := sseServer.Start(addr); err != nil {
		log.Fatalf("服务器启动失败：%v", err)
	}
}

//处理函数
func handleAnalyzeImage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 从请求参数中提取 base64 图片和 prompt
	// 首先转换Arguments为map
	arguments, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultText("错误，无法解析参数"), nil
	}
	
	base64Image, ok := arguments["base64"].(string)
	if !ok {
    	return mcp.NewToolResultText("错误：缺少图片参数"), nil
	}

	prompt, ok := arguments["prompt"].(string)
	if !ok {
	    return mcp.NewToolResultText("错误：缺少提示词参数"), nil
	}

	//创建OpenaiAPI的客户端
	clientConfig := openai.DefaultConfig(config.APIKey)
	clientConfig.BaseURL = config.BaseURL
	client := openai.NewClientWithConfig(clientConfig)

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
							URL: "data:image/jpg;base64," + base64Image,
						},
					},
				},
			},
		},
		MaxTokens: 1000,	//限制返回长度
	})
	
	//处理返回结果
	if err != nil{
		return mcp.NewToolResultText(fmt.Sprintf("调用模型失败：%v", err)), nil
	}

	if len(resp.Choices) > 0 {
		return mcp.NewToolResultText(resp.Choices[0].Message.Content), nil
	}

	return mcp.NewToolResultText("模型未返回有效结果"), nil
	
    // 用config中的信息来创建一个OpenAI客户端
    // 调用视觉模型来分析图片
    // 拿到结果并且返回给LLM
}

//定义配置文件结构体
type Config struct{
	Provider string `json:"provider"`
	BaseURL string `json:"base_url"`
	APIKey string `json:"api_key"`
	Model string `json:"model"`
	Port int `json:"port"`
}

//定义读取配置文件函数
func loadConfig(config Config) Config {
	//打开配置文件
	jsonFile, err := os.Open("config.json")
	//打开失败的错误处理
	if err != nil{
		fmt.Printf("无法打开Json文件：%v\n", err)
	}
	fmt.Println("成功打开了Json文件")

	//尝试解析Json
	if err := json.NewDecoder(jsonFile).Decode(&config); err != nil{
		log.Fatalf("解析JSON配置失败：%v", err)
	}
	//最后需要关闭文件
	defer jsonFile.Close()

	//函数返回
	return config
}
