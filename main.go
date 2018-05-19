package main

import (

	//"fmt"
	//"strconv"
	"hexgo"
	"os"
	"fmt"
	"io/ioutil"
	"path"
	"flag"
)

/*
TODO:
1. 完善CSS。
2. 支持Domain，目录，或者大纲（markdown ToC功能）。
3. 支持Github Pages，重构代码。
4. 集成CI并部署。

速度上线，用约定俗成来代替配置。
 */
const (
	workDir = "/Users/fripside/Go/src/HexGo/"
	DEBUG = true
	ARG_GEN = "gen"
	ARG_DEV = "dev"
)

// 从github pull新内容
func Update()  {

}

func GenBookCache() *hexgo.BlogTheme {
	d, err := os.Getwd()
	dir := workDir
	if err == nil {
		dir = d + "/"
	}
	//fmt.Println("Book Content: ", bookList.Books["go/"].Content)
	confDir := path.Join(dir, "conf")
	cacheDir := path.Join(dir, "cache")
	theme := hexgo.NewBlogTheme(confDir, cacheDir)
	return theme
}

func Dev() {
	h := hexgo.CreateApp(":8888")
	theme := GenBookCache()
	h.SetupFunc = func() {
		fmt.Println("Restart Func")
		theme := GenBookCache()
		h.SetRegisterMap(theme.RequestMap)
	}
	fmt.Println("exec restart func")
	h.RegisterPages(theme)
	h.RegisterStaticDir("/static/", "cache/static")
	h.RegisterStaticDir("/image/","cache/image")
	h.Get("/css", func(context *hexgo.Context) {
		fp, _ := ioutil.ReadFile(workDir + "conf/index.html")
		context.Writer.Write(fp)
	})

	if DEBUG {
		hexgo.LiveReload(h)
	} else {
		h.Run()
	}
}

func main() {
	action := flag.String("action", "gen", "-action [dev] or [gen]")
	flag.Parse()
	if *action == "gen" {
		GenBookCache()
	} else if *action == "dev" {
		Dev()
	}
}
