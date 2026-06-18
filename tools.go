package main

//本文件描述了MCP中存在的所有的工具

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

//处理图像分析函数
func handleAnalyzeImage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 从请求参数中提取 base64 图片和 prompt 以及cache_id
	
	// 对base64处理转换Arguments为map
	arguments, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultText("错误，无法解析参数"), nil
	}

	prompt, ok := arguments["prompt"].(string)
	if !ok {
	    return mcp.NewToolResultText("错误：缺少提示词参数"), nil
	}

	// 首先尝试寻找cache_id
	// 优先检查 cache_id
	var base64Image string
	if cacheID, ok := arguments["cache_id"].(string); ok && cacheID != "" {
	    // 从缓存获取图片
	    imageData, found := cacheGet(cacheID)
	    if !found {
	        return mcp.NewToolResultText("缓存ID不存在或已过期"), nil
	    }
	    base64Image = base64.StdEncoding.EncodeToString(imageData)
	} else if b64, ok := arguments["base64"].(string); ok && b64 != "" {
	    base64Image = b64
	} else {
	    return mcp.NewToolResultText("请提供 cache_id 或 base64 参数"), nil
	}
	
	result, err := callVisionModel(ctx, base64Image, prompt)
	//处理返回结果
	if err != nil{
		return mcp.NewToolResultText(fmt.Sprintf("调用模型失败：%v", err)), nil
	}

	return mcp.NewToolResultText(result), nil
}

//处理图像比较函数
func handleCompareImage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error){
	// 从请求参数中提取 base64 图片和 prompt 以及cache_id
	
	// 对base64处理转换Arguments为map
	arguments, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultText("错误，无法解析参数"), nil
	}

	prompt, ok := arguments["prompt"].(string)
	if !ok {
	    return mcp.NewToolResultText("错误：缺少提示词参数"), nil
	}

	// 首先尝试寻找cache_id
	// 优先检查 cache_id
	var base64ImageOne string
	if cacheID, ok := arguments["cache_id_1"].(string); ok && cacheID != "" {
	    // 从缓存获取图片
	    imageData, found := cacheGet(cacheID)
	    if !found {
	        return mcp.NewToolResultText("缓存ID不存在或已过期"), nil
	    }
	    base64ImageOne = base64.StdEncoding.EncodeToString(imageData)
	} else if b64, ok := arguments["base64_1"].(string); ok && b64 != "" {
	    base64ImageOne = b64
	} else {
	    return mcp.NewToolResultText("请提供 cache_id_1 或 base64_1 参数"), nil
	}

	var base64ImageTwo string
	if cacheID, ok := arguments["cache_id_2"].(string); ok && cacheID != "" {
	    // 从缓存获取图片
	    imageData, found := cacheGet(cacheID)
	    if !found {
	        return mcp.NewToolResultText("缓存ID不存在或已过期"), nil
	    }
	    base64ImageTwo = base64.StdEncoding.EncodeToString(imageData)
	} else if b64, ok := arguments["base64_2"].(string); ok && b64 != "" {
	    base64ImageTwo = b64
	} else {
	    return mcp.NewToolResultText("请提供 cache_id_2 或 base64_2 参数"), nil
	}
	
	result, err := callVisionModelDuo(ctx, base64ImageOne, base64ImageTwo, prompt)
	//处理返回结果
	if err != nil{
		return mcp.NewToolResultText(fmt.Sprintf("调用模型失败：%v", err)), nil
	}

	return mcp.NewToolResultText(result), nil
}