package global

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserToken struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type JWTClaims struct {
	// StandardClaims结构体实现了Claims接口(Valid()函数)
	jwt.StandardClaims
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func GenToken(user UserToken) (signedToken string, err error) {
	TokenConf := Config.Section("token")
	Secret := TokenConf.Key("secretKey").String()
	ExpireTime := TokenConf.Key("expTime").MustInt()

	claims := JWTClaims{
		Id:   user.Id,
		Name: user.Name,
	}
	claims.ExpiresAt = (time.Now().Add(time.Hour * time.Duration(ExpireTime)).Unix())
	claims.IssuedAt = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(Secret))
	if err != nil {
		log.Println(err)
		err = errors.New("Token生成失败")
		return
	}

	//将token存入redis
	err = Rds.Set(strconv.Itoa(int(user.Id)), signedToken, time.Hour*time.Duration(ExpireTime)).Err()
	if err != nil {
		log.Println(err)
		err = errors.New("Token保存Redis失败")
		return
	}

	return
}

func ParseToken(r *http.Request) (user UserToken, err error) {
	TokenConf := Config.Section("token")
	Secret := TokenConf.Key("secretKey").String()

	authString := r.Header.Get("Authorization")
	if authString == "" {
		err = errors.New("请求未携带token，无权限访问")
		return
	}

	token, err := jwt.ParseWithClaims(authString, &JWTClaims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return []byte(Secret), nil
		})

	if err != nil || token.Valid == false {
		log.Println(err)
		err = errors.New("无效的Token")
		return
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		err = errors.New("无效的Token，请重新登录")
		return
	}

	val, err := Rds.Get(strconv.Itoa(int(claims.Id))).Result()
	if err != nil {
		log.Println(err)
		err = errors.New("Token已经过期")
		return
	}

	if authString != val {
		err = errors.New("Token已经刷新")
		return
	}

	user.Id = claims.Id
	user.Name = claims.Name
	return
}
