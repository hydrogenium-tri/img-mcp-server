package main

func main() {
	//解析配置文件
	config = loadConfig(Config{})
	startServer()
}