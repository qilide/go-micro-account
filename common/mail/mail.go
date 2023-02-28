package mail

import (
	"account/common/micro"
	"account/config/logger"
	"account/config/redis"
	"fmt"
	"gopkg.in/gomail.v2"
	"math/rand"
	"time"
)

func InitSendMail(mailTo string, subject string, body string, code string) (string, error) {
	// 设置邮箱主体
	mailConn := map[string]string{
		"user":   micro.ConsulInfo.Email.User,   //发送人邮箱（邮箱以自己的为准）
		"pass":   micro.ConsulInfo.Email.Pass,   //发送人邮箱的密码，现在可能会需要邮箱 开启授权密码后在pass填写授权码
		"host":   micro.ConsulInfo.Email.Host,   //邮箱服务器（此时用的是qq邮箱）
		"rename": micro.ConsulInfo.Email.Rename, //邮箱别名
	}
	m := gomail.NewMessage(
		//发送文本时设置编码，防止乱码。 如果txt文本设置了之后还是乱码，那可以将原txt文本在保存时
		//就选择utf-8格式保存
		gomail.SetEncoding(gomail.Base64),
	)
	m.SetHeader("From", m.FormatAddress(mailConn["user"], mailConn["rename"])) // 添加别名
	m.SetHeader("To", mailTo)                                                  // 发送给用户(可以多个)
	m.SetHeader("Subject", subject)                                            // 设置邮件主题
	m.SetBody("text/html", body)                                               // 设置邮件正文
	/*
	   创建SMTP客户端，连接到远程的邮件服务器，需要指定服务器地址、端口号、用户名、密码，如果端口号为465的话，
	   自动开启SSL，这个时候需要指定TLSConfig
	*/
	d := gomail.NewDialer(mailConn["host"], int(micro.ConsulInfo.Email.Port), mailConn["user"], mailConn["pass"]) // 设置邮件正文
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(m)
	if err != nil {
		return "", err
	}
	return code, nil
}

// SendRegisterMail 发送注册邮件
func SendRegisterMail(mailTo string) (string, error) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	subject := "欢迎注册"
	body := "欢迎注册，您的邮箱验证码为:" + code + "请复制到注册窗口中完成注册，感谢您的支持。验证码在十分钟内有效"
	return InitSendMail(mailTo, subject, body, code)
}

// SendResetPwdMail 发送重置密码邮件
func SendResetPwdMail(mailTo string) (string, error) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	subject := "重置密码"
	body := "欢迎重置密码，您的邮箱验证码为:" + code + "请复制到重置密码窗口中完成重置密码，感谢您的支持。验证码在十分钟内有效"
	return InitSendMail(mailTo, subject, body, code)
}

func GetMail(email string) (interface{}, error) {
	mail, err := redis.Rdb.Get(email).Result()
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return mail, nil
}

func SetMail(email string, code string) {
	redis.Rdb.Set(email+"mail", code, time.Minute*10)
}

func DelMail(email string) {
	redis.Rdb.Del(email + "mail")
}

func CheckMail(email string, code string) (bool, error) {
	redisCode, err := GetMail(email + "mail")
	if fmt.Sprint(redisCode) == code {
		return true, nil
	} else {
		logger.Error(err)
		return false, err
	}
}
