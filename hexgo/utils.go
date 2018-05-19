package hexgo

import (
	//"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"strings"
	//"bytes"
	//"fmt"
	"regexp"
	"bytes"
	"os"
	"io"
	"path/filepath"
	"fmt"
)

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func readYamlConf(path string, data interface{})  {
	conf, e1 := ioutil.ReadFile(path)
	panicError(e1)
	e2 := yaml.Unmarshal(conf, data)
	panicError(e2)
}

// 将url和文件Path变成：/url
func absolutePath(url string) string {
	// 去除中间多个 /
	for strings.Contains(url, "//") {
		url = strings.Replace(url, "//", "/", -1)
	}
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	if url != "/" && strings.HasSuffix(url, "/") {
		i := strings.LastIndex(url, "/")
		url = url[:i]
	}

	return url
}

// 将url和path变成: path/
func relativePath(url string) string {
	for strings.Contains(url, "//") {
		url = strings.Replace(url, "//", "/", -1)
	}
	if strings.HasPrefix(url, "/") {
		url = strings.Replace(url, "/", "", 1)
	}
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return url
}

// go 内置版
func StringReplaceAll(md string, replaceMap map[string]string) string {
	rp := make([]string, 0)
	for k, v := range replaceMap {
		rp = append(rp, k, v)
	}
	fmt.Println(md)
	replacer := strings.NewReplacer(rp...)
	return replacer.Replace(md)
}

func StringReplacePrefix(str, pre, replace string) string {
	if strings.HasPrefix(str, pre) {
		return strings.Replace(str, pre, replace, 1)
	}
	return str
}

func MarkdownPageLinkReplace(md string, replaceMap map[string]string) string {
	buf := bytes.Buffer{}
	re := regexp.MustCompile(`\[[^\]]*\]\s*\(\s*<\s*([^>\s]*)\s*>\s*\)`)
	indexes := re.FindAllStringSubmatchIndex(md, -1)
	i1 := 0
	for _, idx := range indexes {
		s := md[idx[2]:idx[3]]
		//fmt.Println("raw", s)
		if replaceMap[s] != "" {
			buf.WriteString(md[i1:idx[2]])
			i1 = idx[3]
			buf.WriteString(replaceMap[s])
			//fmt.Println("replace", replaceMap[s])
		}
	}
	buf.WriteString(md[i1:])
	//fmt.Println(buf.String())
	return buf.String()
}

// 根据replace map将md中的链接进行替换，eq是否需要严格相等
func MarkdownImageLinkReplace(md string, raw, replace string) string {
	buf := bytes.Buffer{}
	re := regexp.MustCompile(`!\[[^\]]*\]\(\s*([^\)]*)\s*\)`)
	indexes := re.FindAllStringSubmatchIndex(md, -1)
	i1 := 0
	for _, idx := range indexes {
		s := md[idx[2]:idx[3]]
		//fmt.Println("raw", s, raw)
		if strings.Contains(s, raw) {
			buf.WriteString(md[i1:idx[2]])
			i1 = idx[3]
			ns := strings.Replace(s, raw, replace, 1)
			buf.WriteString(ns)
			//fmt.Println("replace", ns)
		}
	}
	buf.WriteString(md[i1:])
	//fmt.Println("length: ", indexes)
	return buf.String()
}

// copy dir, ignore hide files or dir (eg. .scss)
func Copy(src, dest string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	return copy(src, dest, info)
}

func copy(src, dest string, info os.FileInfo) error {
	if info.IsDir() {
		return dirCopy(src, dest, info)
	}
	return fileCopy(src, dest, info)
}

func createFileOrDir(dest string) error {
	return nil
}

func fileCopy(src, dest string, info os.FileInfo) error {
	if strings.HasPrefix(info.Name(), ".") {
		return nil
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	return err
}

func dirCopy(src, dest string, info os.FileInfo) error {
	if strings.HasPrefix(info.Name(), ".") {
		return nil
	}
	if err := os.MkdirAll(dest, info.Mode()); err != nil {
		return err
	}

	infos, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if err := copy(
			filepath.Join(src, info.Name()),
			filepath.Join(dest, info.Name()),
			info,
		); err != nil {
			return err
		}
	}

	return nil
}

func traverseFile(path string, fun func(filePath string, info os.FileInfo)) error {

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if strings.HasPrefix(info.Name(), ".") {
		//fmt.Println("Ignore: ", info.Name())
		return nil
	}

	if info.IsDir() {
		infos, _ := ioutil.ReadDir(path)
		for _, info := range infos {
			traverseFile(filepath.Join(path, info.Name()), fun)
		}
	} else {
		fun(path, info)
	}
	return nil
}

func file_is_exists(f string) bool {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}