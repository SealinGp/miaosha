package service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	. "miaosha/oauth-service/model"
	"net/http"
	"time"
)

var (
	ErrNotSupportGranType  = errors.New("grant type is not supported")
	ErrNotSupportOperation = errors.New("no support operation")
	ErrInvalidUsernameAndPasswordRequest = errors.New("invalid username,password")
	ErrInvalidTokenRequest = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token is expired")
)

type TokenGranter interface {
	Grant(ctx context.Context,grantType string,client *ClientDetails,reader *http.Request) (*OAuth2Token,error)
}

type ComposeTokenGranter struct {
	TokenGrantDict map[string] TokenGranter
}

//token存储
type TokenStore interface {
	//存储
	StoreAccessToken(oauth2Token *OAuth2Token, oauth2Details *OAuth2Details)
	//根据令牌值访问令牌结构体
	ReadAccessToken(tokenValue string) (*OAuth2Token,error)
	//根据令牌值获取客户端和用户信息
	ReadAuth2Details(tokenValue string) (*OAuth2Details,error)
	//根据客户端信息和用户信息获取访问令牌
	GetAccessToken(oauth2Details *OAuth2Details) (*OAuth2Token,error)
	//移除访问令牌
	RemoveAccessToken(tokenValue string)
	// 存储刷新令牌
	StoreRefreshToken(oauth2Token *OAuth2Token, oauth2Details *OAuth2Details)
	// 移除存储的刷新令牌
	RemoveRefreshToken(oauth2Token string)
	//存储刷新访问令牌
	ReadRefreshToken(tokenValue string) (*OAuth2Token,error)
	//根据令牌值获取刷新令牌对应的客户端和用户信息
	ReadOAuth2DetailsForReFreshToken(tokenValue string) (*OAuth2Details, error)
}

//token 加密/解密
type TokenEnhancer interface {
	Enhance(oauth2Token *OAuth2Token,oAuth2Details *OAuth2Details) (*OAuth2Token,error)
	Extract(tokenValue string) (*OAuth2Token,*OAuth2Details,error)
}

type OAuth2TokenCustomClaims struct {
	UserDetails
	ClientDetails
	RefreshToken OAuth2Token
	jwt.StandardClaims
}

type JwtTokenEnhancer struct {
	secretKey []byte
}
func NewJwtTokenEnhancer(secretKey string) TokenEnhancer {
	return &JwtTokenEnhancer{
		secretKey:[]byte(secretKey),
	}
}
func (enhancer *JwtTokenEnhancer)Enhance(oauth2Token *OAuth2Token,oAuth2Details *OAuth2Details) (*OAuth2Token,error)  {
	return enhancer.sign(oauth2Token,oAuth2Details)
}
func (enhancer *JwtTokenEnhancer)Extract(tokenValue string) (*OAuth2Token,*OAuth2Details,error)  {
	token, err := jwt.ParseWithClaims(tokenValue,&OAuth2TokenCustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return enhancer.secretKey,nil
	})
	if err != nil {
		return nil,nil,err
	}
	claims := token.Claims.(*OAuth2TokenCustomClaims)
	exopiresTime := time.Unix(claims.ExpiresAt,0)
	return &OAuth2Token{
		RefreshToken:&claims.RefreshToken,
		TokenValue:tokenValue,
		ExpiresTime:&exopiresTime,
	},&OAuth2Details{
		User:&claims.UserDetails,
		Client:&claims.ClientDetails,
	},nil
}
func (enhancer *JwtTokenEnhancer)sign(oauth2Token *OAuth2Token,oauth2Details *OAuth2Details) (*OAuth2Token,error) {
	expireTime    := oauth2Token.ExpiresTime
	clientDetails := *oauth2Details.Client
	userDetails   := *oauth2Details.User
	clientDetails.ClientSecret = ""
	userDetails.Password = ""

	claims := OAuth2TokenCustomClaims{
		UserDetails:userDetails,
		ClientDetails:clientDetails,
		StandardClaims:jwt.StandardClaims{
			ExpiresAt:expireTime.Unix(),
			Issuer:"System",
		},
	}

	if oauth2Token.RefreshToken != nil {
		claims.RefreshToken = *oauth2Token.RefreshToken
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenValue, err := token.SignedString(enhancer.secretKey)
	if err != nil {
		return nil,err
	}
	oauth2Token.TokenValue = tokenValue
	oauth2Token.TokenType  = "jwt"
	return oauth2Token,nil
}