package logic

import (
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"net"
	"net/url"
	"sync"
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

var mtx sync.Mutex

func (l *FileDownloadLogic) FileDownload(req *types.FileDownloadRequest, userIdentity string) (resp *types.FileDownloadReply, err error) {
	ur := new(models.UserRepository)
	has, err := l.svcCtx.Engine.Where("identity = ? AND user_identity = ?", req.User_repository_identity, userIdentity).Get(ur)
	if err != nil {
		return nil, err
	}
	if !has {
		err = errors.New("用户库中无对应文件")
		return nil, err
	}
	//
	rp := new(models.RepositoryPool)
	has, err = l.svcCtx.Engine.Where("identity = ?", ur.RepositoryIdentity).Get(rp)
	if err != nil {
		return nil, err
	}
	if !has {
		err = errors.New("文件不存在")
		return nil, err
	}
	resp = new(types.FileDownloadReply)
	resp.Ext = rp.Ext
	resp.FileURL = url.QueryEscape(rp.Name) //对filename进行URL编码
	resp.Size = rp.Size
	resp.Hash = rp.Hash

	mtx.Lock()
	listener, err := net.Listen("tcp", ":0") //系统自动分配一个端口号
	if err != nil {
		return
	}
	port := listener.Addr().String()
	port = port[len("[::]"):]
	resp.Port = port
	err = listener.Close()
	if err != nil {
		return
	}
	go helper.Download(rp, port)
	mtx.Unlock()
	//var wg sync.WaitGroup
	//index := strings.Trim(rp.Name, " ")
	/*
		go func(rp *models.RepositoryPool) {
			listener, err := net.Listen("tcp", ":0") //系统自动分配一个端口号
			if err != nil {
				return
			}
			port := listener.Addr().String()
			port = port[len("[::]"):]
			GetPort <- port
			_, err = helper.FileDownloadFromCOSToServer(rp.Path, define.ServerDownloadPath, rp.Ext)
			//ToServerDone <- port
		}(rp)
	*/
	//resp.Port = <-GetPort

	/*
		go func(rp *models.RepositoryPool) {
			select {
			case port := <-ToServerDone:
				{
					server := &http.Server{
						Addr:         "127.0.0.1" + port,
						ReadTimeout:  4800 * time.Second,
						WriteTimeout: 4800 * time.Second,
					}
					/*
						mux := http.NewServeMux()
						mux.Handle("/", http.FileServer(http.Dir(define.ServerDownloadPath+"\\"+rp.Name[:len(rp.Name)-len(rp.Ext)])))
						server.Handler = mux

					http.HandleFunc("/", helper.FileDownloadFromServerToClient)
					log.Fatal(server.ListenAndServe())

				}
			default:
				return
			}
		}(rp)*/
	return
}
