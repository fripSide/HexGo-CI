package hexgo

import (
	"testing"
	"fmt"
	//"os"
	"regexp"
	"path"
)

const (
	root = "/Users/fripside/Go/src/HexGo/"
	blogPath = "conf/blog/"
	cachePath = "conf/cache/"
	tplPath = "conf/template/"
)

func TestGenBooks(t *testing.T)  {
	NewBlogList(root + "conf/blog/")
}

func TestGenBlog(t *testing.T) {
	NewBlogTheme(root + "conf", root + "cache")
}

func TestTrimUrl(t *testing.T)  {
	fmt.Println(absolutePath("/asdad///sadasd"))
	fmt.Println(absolutePath("/asdad/sadasd//"))
	fmt.Println(absolutePath("abc"))
	fmt.Println(absolutePath("abc/"))
	fmt.Println(absolutePath("//"))
}

func TestRegex(t *testing.T)  {
	re := regexp.MustCompile(`\(.*\)`)
	fmt.Println(re.FindAllStringSubmatch("(123) \n aa (456) aa", -1))
	fmt.Println(re.FindAllStringSubmatch("(123) aa (456) aa", -1))
}


func TestMarkdownLinkTransform(t *testing.T) {
	tmp := make(map[string]string)
	tmp["01.md"] = "/go/01.md"
	tmp["013.md"] = "/go/013.md"
	MarkdownPageLinkReplace("[xxx](< 01.md >) asdasd [xxx](<013.md>) ", tmp)
	fmt.Println(MarkdownImageLinkReplace("![](images/1.2.png)", "images/", "/image/go/"))
}

func TestYaml(t *testing.T) {
	blogConf := BlogList{}
	readYamlConf(path.Join(root, "config"), &blogConf)
	fmt.Println(blogConf.Pages)
	fmt.Println(blogConf.Domains)
}


func TestCodeMarkDown() {

}