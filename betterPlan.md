# 优化计划

## 全局错误处理函数
所有方法的错误都要换 以后再说
```go
package globalfail

func fail(msg string, err error) {
	if err != nil {
		log.Printf("%s: %v", msg, err)
	}
	Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
}

//使用示例
if err != nil{
fail("具体错误位置", err)
return
}
```
_____________________________________
## 把自动更新封面封装成函数
预计马上完成
## 全面使用redis
接口有点多 待定