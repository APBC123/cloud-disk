package logic

import (
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadLogic) FileUpload(req *types.FileUploadRequest) (resp *types.FileUploadReply, err error) {
	// 数据入库
	rp := &models.RepositoryPool{
		Identity: helper.UUID(),
		Hash:     req.Hash,
		Name:     req.Name,
		Ext:      req.Ext,
		Size:     req.Size,
		Path:     req.Path,
	}
	_, err = l.svcCtx.Engine.Insert(rp)
	if err != nil {
		err = errors.New("上传记录插入失败")
		return
	}

	resp = new(types.FileUploadReply)
	resp.Identity = rp.Identity
	resp.Ext = rp.Ext
	resp.Name = rp.Name

	return
}
