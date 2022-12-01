package global

import (
	"context"
	"fmt"
	"go-klikdokter/helper/config"
	"go-klikdokter/helper/message"

	"github.com/go-kit/kit/auth/jwt"
	jwtgo "github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type JWTObj struct {
	UserIdLegacy interface{} `json:"user_id_legacy"`
	Fullname     interface{} `json:"fullname"`
	Avatar       interface{} `json:"avatar"`
}

func SetJWTInfoFromContext(ctx context.Context) (JWTObj, message.Message) {
	jwtObj := JWTObj{}
	var avatar string
	defaultAvatar := config.GetConfigString(viper.GetString("image.default-avatar"))
	token, _, err := new(jwtgo.Parser).ParseUnverified(fmt.Sprint(ctx.Value(jwt.JWTContextKey)), jwtgo.MapClaims{})
	if err != nil {
		return jwtObj, message.ErrNoAuth
	}

	if claims, ok := token.Claims.(jwtgo.MapClaims); ok {
		// Get claim value
		var userIdLegacy = claims["id"]
		//if token.Method == jwtgo.SigningMethodRS256 {
		//	userIdLegacy = claims["user_id_legacy"]
		//} else {
		//	userIdLegacy = claims["sub"]
		//}

		fullname := claims["full_name"]
		rawAvatar, ok := claims["avatar"]
		if !ok || rawAvatar == nil {
			avatar = defaultAvatar
		} else {
			avatar = rawAvatar.(string)
			if avatar == "" {
				avatar = defaultAvatar
			}
		}

		// Set value to JWTObj
		if _, ok := userIdLegacy.(float64); ok {
			jwtObj.UserIdLegacy = fmt.Sprintf("%.0f", userIdLegacy)
		} else {
			jwtObj.UserIdLegacy = fmt.Sprintf("%s", userIdLegacy)
		}
		jwtObj.Fullname = fmt.Sprintf("%s", fullname)
		jwtObj.Avatar = avatar

		return jwtObj, message.SuccessMsg
	} else {
		return jwtObj, message.ErrNoAuth
	}
}
