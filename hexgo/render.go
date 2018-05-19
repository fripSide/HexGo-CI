package hexgo

import (
	"net/http"
	T "html/template"
	"fmt"
	"bytes"
	"strings"
	"os"
	"io/ioutil"
	//"net/url"
	"path"
)

/*
TODO:
1. 更换 template，写出一套好看的前端。
http://lewis.suclub.cn/about/
 */

type object = interface{}

func Render(w http.ResponseWriter, tplName string, context map[string]interface{}) {
	tpl, err := T.ParseFiles(tplName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.Execute(w, context)
}

type BlogTheme struct {
	confDir string
	cacheDir   string
	tplMap     map[string]*T.Template
	cache      map[string]*PageCache
	commonLinks map[string]object
	version    string // cache 版本
	RequestMap map[string][]byte
}

type PageCache struct {
	Url string // 路由地址
	Content string // html
	Path string // 保存路径
}

// cache conf 文件
type StaticCache struct {
	url string
	data []byte
}

func NewBlogTheme(confDir, cacheDir string) *BlogTheme {
	theme := new(BlogTheme)
	theme.confDir = confDir
	theme.cacheDir = cacheDir
	tplDir := path.Join(confDir, "template/")
	theme.tplMap = make(map[string]*T.Template)
	pages := []string {"domain", "main", "page", "book", "chapter"}
	layout := path.Join(tplDir, "layout.html")
	for _, name := range pages {
		theme.tplMap[name] = T.Must(T.New(name).ParseFiles(layout, fmt.Sprintf("%s/%s.html", tplDir, name)))
	}
	theme.buildBlog()
	return theme
}

func (theme *BlogTheme) buildBlog() {
	blogDir := path.Join(theme.confDir, "blog/")
	blog := NewBlogList(blogDir)
	fmt.Println(blog)
	theme.cache = make(map[string]*PageCache)
	theme.RequestMap = make(map[string][]byte)
	theme.commonLinks = make(map[string]object)
	// commonLinks
	pageLinks := make(map[string]map[string]string)
	for i := range blog.Pages {
		page := &blog.Pages[i]
		pageLinks[page.Title] = page.ToParamsMap()
	}

	theme.commonLinks["Pages"] = &pageLinks
	domainLinks := make([]object, 0)
	for i := range blog.Domains {
		domain := &blog.Domains[i]
		domainLinks = append(domainLinks, domain.Link.ToParamsMap())
	}
	theme.commonLinks["Domains"] = &domainLinks
	// build html
	theme.buildPages(blog)
	for i := range blog.Domains {
		theme.buildDomain(&blog.Domains[i])
	}
	//theme.RequestMap["/"] = theme.RequestMap["/index"]
}

func (theme *BlogTheme) buildPages(blog *BlogList) {
	params := make(map[string]object)
	params["Links"] = &theme.commonLinks
	for i := range blog.Pages {
		page := &blog.Pages[i]
		params["Title"] = page.Title
		params["Content"] = T.HTML(page.Md)
		pageHtml := theme.executeTpl("page", params)
		theme.genPageCache(page, pageHtml)
		theme.copyImages(page)
	}
	theme.copyStaticFiles()
}

func (theme *BlogTheme) buildDomain(domain *DomainConf) {
	// domain cover
	params := make(map[string]object)
	params["Links"] = &theme.commonLinks
	params["Title"] = domain.Link.Title
	bookLinks := make([]object, 0)
	for i := range domain.Books {
		book := &domain.Books[i]
		info := book.Link.ToParamsMap()
		info["Desc"] = book.Desc
		info["Cover"] = book.Cover
		bookLinks = append(bookLinks, info)
	}
	params["Books"] = bookLinks
	pageHtml := theme.executeTpl("domain", params)
	theme.genPageCache(&domain.Link, pageHtml)
	for i := range domain.Books {
		theme.buildBook(&domain.Books[i])
	}
}

func (theme *BlogTheme) buildBook(book *BookConf) {
	// book cover
	params := make(map[string]object)
	params["Links"] = &theme.commonLinks
	params["Cover"] = book.Cover
	params["Desc"] = book.Desc
	chapterLinks := make([]object, 0)
	for i := range book.Chapters {
		chapter := &book.Chapters[i]
		info := chapter.Link.ToParamsMap()
		info["Date"] = chapter.Date
		chapterLinks = append(chapterLinks, info)
	}
	params["Chapters"] = chapterLinks
	coverHtml := theme.executeTpl("book", params)
	theme.genPageCache(&book.Link, coverHtml)
	theme.copyImages(&book.Link)

	// book page
	for i := range book.Chapters {
		chapter := &book.Chapters[i]
		params["Content"] = T.HTML(chapter.Link.Md)
		pageHtml := theme.executeTpl("chapter", params)
		theme.genPageCache(&chapter.Link, pageHtml)
	}
}

func (theme *BlogTheme) genPageCache(page *PathLink, html string) {
	cachePath := path.Join(theme.cacheDir, page.Url, "index.html")
	cache := PageCache{page.Url, html, cachePath}
	theme.cache[cache.Url] = &cache
	cache.saveToFile()
	theme.RequestMap[cache.Url] = []byte(cache.Content)
}

func (theme *BlogTheme) executeTpl(name string, data map[string]object) string {
	var w bytes.Buffer
	err := theme.tplMap[name].ExecuteTemplate(&w, "layout", data)
	if err != nil {
		fmt.Println(err.Error())
	}
	return w.String()
}

func (theme *BlogTheme) copyStaticFiles() {
	dest := path.Join(theme.cacheDir, "static/")
	src := path.Join(theme.confDir, "template/static/")
	Copy(src, dest)
}

func (theme *BlogTheme) copyImages(page *PathLink)  {
	if page.LocalImage != "" {
		imagePath := path.Join(theme.cacheDir,"/image/", page.Url)
		Copy(page.LocalImage, imagePath)
		//fmt.Println("copy image", page.LocalImage, imagePath)
	}
}

func (c *PageCache) saveToFile() {
	i := strings.LastIndex(c.Path, "/")
	if i < 0 {
		return
	}
	d := c.Path[0:i]
	e1 := os.MkdirAll(d, 0777)
	if e1 != nil {
		panic(e1)
	}
	//fmt.Printf("Success to create %s\n", d)
	// 文件只有读写权限
	err := ioutil.WriteFile(c.Path, []byte(c.Content), 0666)
	if err != nil {
		panic(err)
	}
	//fmt.Println("Save: ", c.Path)
}
