package logic

import (
	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/models"
	"context"
	"errors"
	"time"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MailCodeSendRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMailCodeSendRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MailCodeSendRegisterLogic {
	return &MailCodeSendRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MailCodeSendRegisterLogic) MailCodeSendRegister(req *types.MailCodeSendRequest) (resp *types.MailCodeSendReply, err error) {
	// todo: add your logic here and delete this line
	//该邮箱在数据库中不存在
	cnt, err := l.svcCtx.Engine.Where("email=?", req.Email).Count(new(models.UserBasic))
	if err != nil {
		return
	}
	if cnt > 0 {
		err = errors.New("该邮箱被注册")
		return
	}
	// 校验验证码是否在有效期内
	codeTTL, _ := l.svcCtx.RDB.TTL(l.ctx, req.Email).Result()
	if codeTTL.Seconds() > 0 || codeTTL.Seconds() == -1 {
		return nil, errors.New("the verify code has not expired")
	}
	//获取验证码
	code := helper.RandCode()
	//存储验证码至Redis中
	l.svcCtx.RDB.Set(l.ctx, req.Email, code, time.Second*time.Duration(define.CodeExpire))
	//发送验证码
	err = helper.MailSendCode(req.Email, code)

	return
}
