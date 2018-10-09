package library

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

/**
获取源码
*/
func FetchSource(url string) (doc *goquery.Document) {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	dec := mahonia.NewDecoder("gb18030")
	rd := dec.NewReader(res.Body)

	doc, err = goquery.NewDocumentFromReader(rd)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

/**
下载文件
*/
func FetchFile(url, imagePath, referer string) (err error, fullpath string, size int64) {
	exist, err := PathExists(imagePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !exist {
		err = os.MkdirAll(imagePath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	i := strings.LastIndex(url, ".")
	j := strings.LastIndex(url, "-")
	suffix := url[i:j]

	filename := url
	filenameBype := []byte(filename)
	md5Filename := md5.Sum(filenameBype)
	filename = fmt.Sprintf("%x%s", md5Filename, suffix)
	//fmt.Println(filename, suffix)
	filepath := imagePath + "/" + filename

	exist, err = PathExists(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if exist {
		err = errors.New("文件已存在")
		return
	}

	imageFile, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("[writeImage create file]: fileName: %s\n href: %s\nerror: %s\n", filepath, url, err.Error())
		return
	}

	client := &http.Client{}

	//提交请求
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//增加header选项
	reqest.Header.Add("NT", "1")
	reqest.Header.Add("If-Modified-Since", "Thu, 06 Sep 2018 03:54:19 GMT")
	reqest.Header.Add("If-None-Match", "BDE9E8B0317BF99A37BE8FE52763AF1E")
	reqest.Header.Add("Referer", referer)

	//处理返回结果
	res, _ := client.Do(reqest)

	//fmt.Println(res.StatusCode)
	if res.StatusCode != 200 {
		fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status)
		os.Remove(filepath)
		return
	}

	size, err = io.Copy(imageFile, res.Body)
	if err != nil {
		fmt.Printf("io.Copy: error: %s  href: %s\n", err.Error(), url)
		os.Remove(filepath)
		return
	}
	fmt.Printf("Get From %s: %d bytes\n", url, size)

	fullpath = filepath
	return
}
