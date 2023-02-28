package token

import (
	"account/config/redis"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 2

var MySecret = []byte("Account")

// GenToken 生成JWT
func GenToken(UserName string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		UserName, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "Account",                                  // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}

func GetToken(userName string) (interface{}, error) {
	token, err := redis.Rdb.Do("Get", userName+"token").Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

func SetToken(userName string, token string) {
	redis.Rdb.Set(userName+"token", token, time.Minute*30)
}

func DelToken(userName string) {
	redis.Rdb.Del(userName + "token")
}

func CheckToken(userName string, token string) (bool, error) {
	redisToken, err := GetToken(userName + "token")
	var mc *MyClaims
	mc, err = ParseToken(fmt.Sprint(redisToken))
	if err != nil {
		return false, err
	}
	if mc.UserName == userName {
		return true, nil
	} else {
		return false, errors.New("登录信息错误")
	}
}
