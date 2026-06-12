package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

//所有负责处理upload的程序代码放在此处
func uploadHandler(w http.ResponseWriter, r *http.Request){
	//检测请求
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "不合理的请求")
		return
	}

	//处理上传的数据
	err := r.ParseMultipartForm(16<<20)		//最大16MB
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "上传数据错误：%v", err)
		return
	}

	//获取上传的文件
	file, _, err := r.FormFile("file")
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "文件错误")
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "文件读取失败：%v", err)
		return
	}
	hash := sha256.Sum256(fileData)		//计算了图片的HASH值来区分图片
	cacheID := fmt.Sprintf("%x", hash)
	cacheSet(cacheID, fileData)		//向缓存中存储图片

	//返回Json响应
	w.Header().Set("Content-Type", "application/json")

	//创建返回响应
	json.NewEncoder(w).Encode(map[string]string{
    	"cache_id": cacheID,
	})
}
