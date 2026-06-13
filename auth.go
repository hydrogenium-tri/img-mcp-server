package main

import(
	"net/http"
)

//认证文件，判断客户端是否有权限访问

func authMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		header := r.Header.Get("Authorization")
		if config.AuthToken == "" || "Bearer " + config.AuthToken == header {
			next.ServeHTTP(w, r)
			return 
		}	else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
	})
}