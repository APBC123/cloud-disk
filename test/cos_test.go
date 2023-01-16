package test

import (
	"bytes"
	"cloud-disk/core/define"
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestFileUpLoadByFilepath(t *testing.T) {
	u, _ := url.Parse("https://2290312980-1316376654.cos.ap-nanjing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := "cloud-disk/1671180003499.jpg"

	_, _, err := client.Object.Upload(
		context.Background(), key, "./img/1671180003499.jpg", nil,
	)
	if err != nil {
		panic(err)
	}

}

func TestFileUpLoadByReader(t *testing.T) {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := "cloud-disk/1671180003499.jpg"

	f, err := os.ReadFile("./img/1671180003499.jpg")

	_, err = client.Object.Put(
		context.Background(), key, bytes.NewReader(f), nil,
	)
	if err != nil {
		panic(err)
	}

}
