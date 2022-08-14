package global

import (
	"context"
	"fmt"
	"go-klikdokter/helper/message"

	"github.com/go-kit/kit/auth/jwt"
	jwtgo "github.com/golang-jwt/jwt"
)

type JWTObj struct {
	UserIdLegacy string `json:"user_id_legacy"`
	Avatar       string `json:"avatar"`
}

func SetJWTInfoFromContext(ctx context.Context) (JWTObj, message.Message) {
	jwtObj := JWTObj{}
	token, _, err := new(jwtgo.Parser).ParseUnverified(fmt.Sprint(ctx.Value(jwt.JWTContextKey)), jwtgo.MapClaims{})
	if err != nil {
		return jwtObj, message.ErrNoAuth
	}

	if claims, ok := token.Claims.(jwtgo.MapClaims); ok {
		userIdLegacy := claims["sub"].(float64)
		var avatar string
		defaultAvatar := "https://asset-cdn.medkomtek.com/assets/images/profile/user-default-original.jpg"
		rawAvatar, ok := claims["avatar"]

		if !ok || rawAvatar == nil {
			avatar = defaultAvatar
		} else {
			avatar = rawAvatar.(string)

			if avatar == "" {
				avatar = defaultAvatar
			}
		}

		jwtObj.UserIdLegacy = fmt.Sprintf("%.0f", userIdLegacy)
		jwtObj.Avatar = fmt.Sprintf("%s", avatar)
		return jwtObj, message.SuccessMsg
	} else {
		return jwtObj, message.ErrNoAuth
	}
}
