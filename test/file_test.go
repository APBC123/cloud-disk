package test

import (
	"bufio"
	"cloud-disk/core/define"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
)

// 分片大小
const chunkSize = 100 * 1024 * 1024 //100MB
// 文件分片
func TestGenerateChunkFile(t *testing.T) {
	fileInfo, err := os.Stat("./img/12.jpg")
	if err != nil {
		t.Fatal(err)
	}
	//分片的个数
	chunkNum := math.Ceil(float64(fileInfo.Size()) / float64(chunkSize)) //Ceil向上取整
	myFile, err := os.OpenFile("./img/12.jpg", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	b := make([]byte, chunkSize)
	r := bufio.NewReader(myFile)
	for i := 0; i < int(chunkNum); i++ {
		//指定读取文件的起始位置
		_, err = myFile.Seek(int64(i*chunkSize), 0)
		if err != nil {
			t.Fatal(err)
		}
		if chunkSize > fileInfo.Size()-int64(i*chunkSize) {
			b = make([]byte, fileInfo.Size()-int64(i*chunkSize))
		}
		_, err = r.Read(b)
		if err != nil {
			t.Fatal(err)
		}
		f, err := os.OpenFile("./"+strconv.Itoa(i)+".chunk", os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			t.Fatal(err)
		}
		w := bufio.NewWriter(f)
		w.Write(b) //f.Write(b)
		w.Flush()
		f.Close()

	}
	myFile.Close()
}

// 分片文件的合并
func TestMergeChunk(t *testing.T) {
	myFile, err := os.OpenFile("test2.mkv", os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		t.Fatal(err)
	}
	fileInfo, err := os.Stat("test.mkv")
	if err != nil {
		t.Fatal(err)
	}
	//分片的个数
	chunkNum := math.Ceil(float64(fileInfo.Size()) / float64(chunkSize)) //Ceil向上取整
	for i := 0; i < int(chunkNum); i++ {
		f, err := os.OpenFile("./"+strconv.Itoa(i)+".chunk", os.O_WRONLY, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		b, err := os.ReadFile("./" + strconv.Itoa(i) + ".chunk")
		if err != nil {
			t.Fatal(err)
		}
		w := bufio.NewWriter(myFile)
		w.Write(b) //myFile.Write(b)
		w.Flush()
		f.Close()
	}
	myFile.Close()
}

// 文件一致性校验
func TestCheck(t *testing.T) {
	//获取原始文件的信息
	f1, err := os.OpenFile("test.mkv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	b1, err := os.ReadFile("test.mkv")
	if err != nil {
		t.Fatal(err)
	}
	//获取merge后文件的信息
	f2, err := os.OpenFile("test2.mkv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := os.ReadFile("test2.mkv")
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	f1.Close()
	f2.Close()
	s1 := fmt.Sprintf("%x", md5.Sum(b1))
	s2 := fmt.Sprintf("%x", md5.Sum(b2))
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s1 == s2)
}

// 文件由腾讯云下载到后端服务器

func TestFileDownload(t *testing.T) {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := "cloud-disk/[POPGO][Ghost_in_the_Shell][S.A.C._2nd_GIG][01][AVC_FLACx2+AC3][BDrip][1080p][228A2730].mkv"
	file := "G:\\localfile.mkv"
	optDownload := &cos.MultiDownloadOptions{
		ThreadPoolSize: 8,
		CheckPoint:     true,
	}
	_, err := client.Object.Download(
		context.Background(), key, file, optDownload,
	)
	if err != nil {
		panic(err)
	}
	return
}
