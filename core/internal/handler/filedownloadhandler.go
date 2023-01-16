package handler

import (
	"cloud-disk/core/models"
	"errors"
	"net/http"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileDownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileDownloadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		ur := new(models.UserRepository)
		has, err := svcCtx.Engine.Where("identity = ?", req.User_repository_identity).Get(ur)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		if !has {
			err = errors.New("用户库中无对应文件")
			httpx.Error(w, err)
			return
		}
		//
		rp := new(models.RepositoryPool)
		has, err = svcCtx.Engine.Where("identity = ?", ur.RepositoryIdentity).Get(rp)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		if !has {
			err = errors.New("文件不存在")
			httpx.Error(w, err)
			return
		}
		l := logic.NewFileDownloadLogic(r.Context(), svcCtx)
		resp, err := l.FileDownload(rp)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
