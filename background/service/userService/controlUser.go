package userService

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"takeout/db"
	"takeout/models"
	"takeout/models/vo"
	"takeout/utils"
	"time"
)

type UserLoginDTO struct {
	Code string `json:"code"`
}

type WeChatLoginVO struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	ErrMsg     string `json:"errmsg"`
	ErrCode    int32  `json:"errcode"`
	UnionId    string `json:"unionid"`
}

func Login(e echo.Context) error {
	var userLogin UserLoginDTO
	if err := e.Bind(&userLogin); err != nil {
		result := new(models.Result[interface{}])
		result.Error(err.Error())
		fmt.Println("没有code")
		return e.JSON(http.StatusBadRequest, result)
	}
	//基础url
	getUrl := "https://api.weixin.qq.com/sns/jscode2session"

	// 创建URL值并添加查询参数
	params := url.Values{}
	params.Add("appid", utils.WeChatAppId)
	params.Add("secret", utils.WeChatAppSecret)
	params.Add("js_code", userLogin.Code)
	params.Add("grant_type", "authorization_code")

	// 构造完整的URL
	fullURL := fmt.Sprintf("%s?%s", getUrl, params.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		result := new(models.Result[interface{}])
		result.Error(err.Error())
		fmt.Println("请求失败")
		return e.JSON(http.StatusBadRequest, result)
	}

	var weChatLogin WeChatLoginVO
	err = json.NewDecoder(resp.Body).Decode(&weChatLogin)
	if err != nil {
		result := new(models.Result[interface{}])
		result.Error(err.Error())
		fmt.Println("解析失败")
		return e.JSON(http.StatusBadRequest, result)
	}

	if weChatLogin.ErrCode != 0 {
		result := new(models.Result[interface{}])
		result.Error(weChatLogin.ErrMsg)
		fmt.Printf("错误码: %d\n", weChatLogin.ErrCode)
		fmt.Println("登录失败" + weChatLogin.ErrMsg)
		return e.JSON(http.StatusBadRequest, result)
	}

	var user models.User
	db.DB.Where("openid = ?", weChatLogin.OpenId).First(&user)
	if user == (models.User{}) {
		user.CreateTime = utils.CustomTime{Time: time.Now()}
		user.OpenId = weChatLogin.OpenId
		user.Avatar = "https://heathen-project.oss-cn-beijing.aliyuncs.com/672f0b29b3fc49b6aca8b99ac2f500e7.jpg"
		user.Name = "用户" + weChatLogin.OpenId
		err = db.DB.Create(&user).Error
		if err != nil {
			result := new(models.Result[interface{}])
			result.Error(err.Error())
			fmt.Println("用户创建失败")
			return e.JSON(http.StatusBadRequest, result)
		}
	}

	if user.Id == 0 {
		db.DB.Where("open_id = ?", weChatLogin.OpenId).Find(&user)
	}
	// 创建jwt
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.Id
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	tokenAnswer, err := token.SignedString([]byte("itheima"))

	if err != nil {
		result := new(models.Result[interface{}])
		result.Error(err.Error())
		fmt.Println("token生成失败")
		return e.JSON(http.StatusBadRequest, result)
	}
	result := new(models.Result[vo.UserVO])
	result.SuccessWithObject(vo.UserVO{
		Id:     user.Id,
		Token:  tokenAnswer,
		OpenId: user.OpenId,
	})
	return e.JSON(http.StatusOK, result)
}
