package logic

import (
	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"net"
	"net/http"
	"time"
)

type FileDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDownloadLogic {
	return &FileDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var ToServerDone = make(chan error)
var GetPort = make(chan string)

func (l *FileDownloadLogic) FileDownload(rp *models.RepositoryPool) (resp *types.FileDownloadReply, err error) {
	resp = new(types.FileDownloadReply)
	resp.Ext = rp.Ext
	resp.Name = rp.Name
	resp.Size = rp.Size
	resp.Hash = rp.Hash
	//var wg sync.WaitGroup
	//index := strings.Trim(rp.Name, " ")

	go func(rp *models.RepositoryPool) {
		_, err = helper.FileDownloadFromCOSToServer(rp.Path, define.ServerDownloadPath, rp.Name, rp.Ext)
		ToServerDone <- err
		return
	}(rp)

	go func(rp *models.RepositoryPool) {
		err = <-ToServerDone
		if err != nil {
			return
		}
		listener, err := net.Listen("tcp", ":0") //系统自动分配一个端口号
		if err != nil {
			panic(err)
		}
		port := listener.Addr().String()
		port = port[len("[::]"):]
		GetPort <- port
		go func(port string, rp *models.RepositoryPool) {
			server := &http.Server{
				Addr:         "127.0.0.1" + port,
				ReadTimeout:  2 * time.Second,
				WriteTimeout: 2 * time.Second,
			}
			mux := http.NewServeMux()
			mux.Handle("/", http.FileServer(http.Dir(define.ServerDownloadPath+"\\"+rp.Name[:len(rp.Name)-len(rp.Ext)])))
			server.Handler = mux
			err = server.ListenAndServe()
			if err != nil {
				return
			}
			<-GetPort
		}(port, rp)
	}(rp)
	resp.Port = <-GetPort
	return
}
