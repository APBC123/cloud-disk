package logic

import (
	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterRequest) (resp *types.UserRegisterReply, err error) {
	// todo: add your logic here and delete this line
	//判断Code是否一致
	code, err := l.svcCtx.RDB.Get(l.ctx, req.Email).Result()
	if err != nil {
		return nil, errors.New("未获取该邮箱的验证码")
	}
	if code != req.Code {
		err = errors.New("验证码错误")
		return
	}
	//判断用户名是否已存在
	cnt, err := l.svcCtx.Engine.Where("name=?", req.Name).Count(new(models.UserBasic))
	if cnt > 0 {
		err = errors.New("用户名已存在")
		return
	}
	//数据入库
	user := &models.UserBasic{
		Name:        req.Name,
		Password:    helper.Md5(req.Password),
		Email:       req.Email,
		Identity:    helper.UUID(),
		NowVolume:   define.NewVolume,
		TotalVolume: define.NewVolume,
	}
	n, err := l.svcCtx.Engine.Insert(user)
	if err != nil {
		return nil, err
	}
	log.Println("insert user row:", n)
	return
}
