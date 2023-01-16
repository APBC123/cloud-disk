package handler

import (
	"cloud-disk/core/helper"
	"cloud-disk/core/models"
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"path"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 获取上传的文件
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			err = errors.New("获取上传文件失败")
			httpx.Error(w, err)
			return
		}

		// 判断文件在数据库中是否已经存在
		b := make([]byte, fileHeader.Size)
		_, err = file.Read(b)
		if err != nil {
			err = errors.New("文件未成功读取")
			httpx.Error(w, err)
			return
		}
		hash := fmt.Sprintf("%x", md5.Sum(b))
		rp := new(models.RepositoryPool)
		has, err := svcCtx.Engine.Where("hash = ?", hash).Get(rp)
		if err != nil {
			err = errors.New("hash匹配错误")
			httpx.Error(w, err)
			return
		}
		if has {
			// 文件已经存在，直接返回信息
			err = errors.New("文件已存在")
			httpx.Error(w, err)
			return
		}

		// 往 COS 中存储文件
		cosPath, err := helper.CosUpload(r)
		if err != nil {
			err = errors.New("向COS中储存文件时出错")
			httpx.Error(w, err)
			return
		}

		// 往 logic 传递 request
		req.Name = fileHeader.Filename
		req.Ext = path.Ext(fileHeader.Filename)
		req.Size = fileHeader.Size
		req.Hash = hash
		req.Path = cosPath

		l := logic.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
