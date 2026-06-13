package main

import(
	"encoding/json"
	"fmt"
	"log"
	"os"
)
//所有配置相关代码在这个文件中
// 定义配置文件结构体
type Config struct {
	Provider string `json:"provider"`
	BaseURL  string `json:"base_url"`
	APIKey   string `json:"api_key"`
	Model    string `json:"model"`
	Port     int    `json:"port"`
	AuthToken string `json:"auth_token"`
}

// 定义读取配置文件函数
func loadConfig(config Config) Config {
	//打开配置文件
	jsonFile, err := os.Open("config.json")
	//打开失败的错误处理
	if err != nil {
		fmt.Printf("无法打开Json文件：%v\n", err)
	}
	fmt.Println("成功打开了Json文件")

	//尝试解析Json
	if err := json.NewDecoder(jsonFile).Decode(&config); err != nil {
		log.Fatalf("解析JSON配置失败：%v", err)
	}
	//最后需要关闭文件
	defer jsonFile.Close()

	//函数返回
	return config
}

var config Config
