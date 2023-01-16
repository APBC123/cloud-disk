package define

import (
	"github.com/dgrijalva/jwt-go"
	"os"
)

type UserClaim struct {
	Id       int
	Identity string
	Name     string
	jwt.StandardClaims
}

var MailPassword = os.Getenv("MailPassword")

var JwtKey = "cloud-disk-key"

// 验证码的长度
var CodeLength = 6

// 验证码的过期时间（s）
var CodeExpire = 300

// 腾讯云对象储存
var TencentSecretID = os.Getenv("TencentSecretID")
var TencentSecretKey = os.Getenv("TencentSecretKey")
var CosBucket = "https://2290312980-1316376654.cos.ap-nanjing.myqcloud.com"

// 分页的默认参数
var PageSize = 20

var Datetime = "2006-01-01 15:05:05"

var TokenExpireTime = 3600
var RefreshTokenExpireTime = 5400
var ServerDownloadPath = "G:\\ServerDownloadPath"
