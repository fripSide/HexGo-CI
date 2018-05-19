package hexgo

import (
	_ "log"
	_ "regexp"
	//"reflect"
	"net/http"
	"fmt"
	"context"
)

/*
封装自己的Web框架
https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/13.3.md
https://github.com/gin-gonic/gin
https://github.com/go-martini/martini

https://github.com/julienschmidt/httprouter

TODO:
添加file cache
 */

type HandleFunc func(*Context)

type Context struct {
	Writer http.ResponseWriter
	Request *http.Request
	params map[string]string // url传的参数
}

type HexGo struct {
	Host      string
	server    *http.Server
	reqMap    map[string][]byte // 记录请求响应
	SetupFunc func()
}

func CreateApp(host string) *HexGo {
	h := new(HexGo)
	server := &http.Server{Addr: host, Handler: nil}
	h.server = server
	fmt.Println("Listen: ", host)
	return h
}

func (h *HexGo) Get(pattern string, handler HandleFunc) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		c := &Context{w, r, nil}
		handler(c)
	})
}

func (h *HexGo) Run() {
	fmt.Println("Start Server")
	err := h.server.ListenAndServe()
	if err != nil {
		fmt.Errorf(err.Error())
	}
}

func (h *HexGo) Stop() {
	h.server.Shutdown(context.Background())
	fmt.Println("Stop Server")
}

func (h *HexGo) SetRegisterMap(data map[string][]byte) {
	h.reqMap = data
}

func (h *HexGo) RegisterPages(theme *BlogTheme) {
	h.reqMap = theme.RequestMap
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := absolutePath(r.URL.Path)

		if r.Method == http.MethodGet {
			fmt.Println("GET ", r.URL)

			if h.reqMap[url] != nil {
				w.Write(h.reqMap[url])
			} else {
				w.Write([]byte(fmt.Sprint("Page not found: ", r.URL)))
			}
		}

	})
}

func (h *HexGo) RegisterStaticDir(pattern, dir string) {
	fs := http.FileServer(http.Dir(dir))
	http.Handle(pattern, http.StripPrefix(pattern, fs))
}
