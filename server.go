package main

//所有关于MCP服务器的启动、工具注册等相关代码在这个文件中

import (
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"fmt"
)

func startServer() {
	//创建服务器实例
	s := server.NewMCPServer(
		"image-recognition",
		"0.0.1",
	)

	//新建一个工具
	analyzeImage := mcp.NewTool("analyze_image",
		mcp.WithDescription("传入图像字符串，使用一个多模态模型解析它"),
		mcp.WithString("base64",
      		mcp.Description("字符串，base64编码的图像数据（与cache_id二选一）"),
		),
		mcp.WithString("prompt",
   			mcp.Required(),
      		mcp.Description("字符串，对图像模型的提示词"),
		),
		mcp.WithString("cache_id",
			mcp.Description("上传图片后返回的缓存ID（与base64二选一）"),
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
	mux := http.NewServeMux()
	mux.Handle("/mcp", sseServer)
	mux.HandleFunc("/upload", uploadHandler)

	http.ListenAndServe(addr, mux)
}