package hexgo
// go get -v 获取新包
import (
	"github.com/shurcooL/github_flavored_markdown"
	//"gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"fmt"
	//"gopkg.in/yaml.v2"
	"path"
)
/*
网站的组织架构：
domain: header栏，领域，例如：编程语言
	book: 一序列文章，可组织成书，或者独立文章，例如：Java基础教程
		chapter: 一篇文章，可设置收费，免费，等各种字段

pages: 独立的页面，例如：about
*/

type BlogList struct {
	BaseDir string `yaml:"-"`
	Version string
	Domains []DomainConf
	Pages []PathLink
}

type DomainConf struct {
	Link PathLink	`yaml:",inline"`
	Books []BookConf
}

type BookConf struct {
	Link PathLink	`yaml:",inline"`
	Desc string
	Cover string
	Chapters []ChapterConf
}

type ChapterConf struct {
	Link PathLink	`yaml:",inline"`
	Date string
}

// 指向文件夹，或者文件
type PathLink struct{
	Title string
	Path string
	Url string
	LocalImage string `yaml:"local_image"`
	Md string `yaml:"-"`
}

func NewBlogList(dir string) *BlogList {
	blogList := new(BlogList)
	blogList.BaseDir = dir
	blogList.genAllPages()
	return blogList
}

func (b *BlogList) genAllPages() {
	blogConf := b
	readYamlConf(path.Join(b.BaseDir, "config"), &blogConf)
	// pages
	for i := range blogConf.Pages {
		page := &blogConf.Pages[i]
		page.Path = path.Join(b.BaseDir, page.Path)
		page.Url = absolutePath(page.Url)
		d := genMarkDown(page.Path, func(s string) string {
			return TransformImageLink(s, page.LocalImage, page.Url)
		})
		page.Md = d
		if page.LocalImage != "" {
			page.LocalImage = path.Join(b.BaseDir, "pages", page.LocalImage)
		}
	}
	//fmt.Println("Md: ", blogConf.Pages)
	// domain
	base := PathLink{"", b.BaseDir, "", "", ""}
	for i := range blogConf.Domains {
		domain := &blogConf.Domains[i]
		domain.Link.mergePath(&base)
		b.genDomain(domain)
	}
}

func (b *BlogList) genDomain(domain *DomainConf) {
	//fmt.Println("pase domain ", domain.Link)
	readYamlConf(path.Join(domain.Link.Path, "config"), domain)
	//fmt.Println(domain)
	for i := range domain.Books {
		book := &domain.Books[i]
		url := book.Link.Url
		book.Link.mergePath(&domain.Link)
		b.genBook(book, url)
	}
}

func (b *BlogList) genBook(book *BookConf, bookUrl string) {
	readYamlConf(path.Join(book.Link.Path, "config"), book)
	// 将文章中的相对链接改成绝对路径
	file2Url := make(map[string]string)
	for i := range book.Chapters {
		chapter := &book.Chapters[i]
		p := chapter.Link.Path
		chapter.Link.mergePath(&book.Link)
		file2Url[p] = chapter.Link.Url
	}

	for i := range book.Chapters {
		chapter := &book.Chapters[i]
		chapter.Link.Md = genMarkDown(chapter.Link.Path, func(str string) string {
			str = TransformImageLink(str, book.Link.LocalImage, bookUrl)
			str = MarkdownPageLinkReplace(str, file2Url)
			return str
		})
	}

	if book.Link.LocalImage != "" {
		localImage := relativePath(book.Link.LocalImage)
		book.Link.LocalImage = path.Join(book.Link.Path, localImage)
		// 尝试将封面的路径替换成绝对路径
		imgUrl := absolutePath(fmt.Sprintf("/image/%s/", book.Link.Url)) + "/"
		book.Cover = StringReplacePrefix(book.Cover, localImage, imgUrl)
	}
}

func (link *PathLink) mergePath(parent *PathLink) {
	link.Path = path.Join(parent.Path, link.Path)
	link.Url = absolutePath(parent.Url + "/" + link.Url)
}

func (link *PathLink) ToParamsMap() map[string]string {
	info := make(map[string]string)
	info["Link"] = link.Url
	info["Title"] = link.Title
	return info
}

func genMarkDown(path string, transformFunc func(str string) string) string {
	//fmt.Println("Gen Page: ", path)
	fp, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err.Error())
		return string(err.Error())
	}
	data := transformFunc(string(fp))


	//return string(blackfriday.Run([]byte(data)))
	return string(github_flavored_markdown.Markdown([]byte(data)))
}


func TransformImageLink(md, localImg, url string) string {
	if localImg != "" {
		imgUrl := absolutePath(fmt.Sprintf("/image/%s/", url)) + "/"
		md = MarkdownImageLinkReplace(md, relativePath(localImg), imgUrl)
	}
	return md
}