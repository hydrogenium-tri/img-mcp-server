package main

import (
	"sync"
	"time"
)

//这个文件描述了图片的缓存机制

var (
    cache     = make(map[string]imageCache)
    cacheMu   sync.RWMutex
    cacheTTL  = 10 * time.Minute   // 10 分钟过期
)

//定义了存储图片以及清理图片的缓存结构体
type imageCache struct{
	Data []byte		//图片的原始数据
	ExpriresAt time.Time		//过期时间
}

//向缓存中写入图片数据
func cacheSet(id string, data []byte){
	cacheMu.Lock()		//写时锁死缓存
	cache[id] = imageCache{Data: data, ExpriresAt: time.Now().Add(cacheTTL)}
	cacheMu.Unlock()	//写完放开缓存
}

//从缓存中提取出图片数据
func cacheGet(id string)([]byte, bool){
	cacheMu.RLock()		//使用读取锁防止在读取时被修改，同时允许多线程读取
	entry, ok := cache[id]
	cacheMu.RUnlock()	//读取完解锁
	if !ok || time.Now().After(entry.ExpriresAt){
		return nil, false
	}
	return entry.Data, true
}

//清理函数，每一会运行一次清除过期缓存
func startCacheCleanup(){
	go func(){
		ticker := time.NewTicker(time.Minute)
		for range ticker.C{
			cacheMu.Lock()
			for id, entry := range cache{
				if time.Now().After(entry.ExpriresAt){
					delete(cache, id)
				}
			}
			cacheMu.Unlock()
		}
	}()
}