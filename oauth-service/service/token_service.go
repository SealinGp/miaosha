package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/text/date"
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

//token加密的所有类型分发使用
type TokenGranter interface {
	Grant(ctx context.Context,grantType string,client *ClientDetails,reader *http.Request) (*OAuth2Token,error)
}
type ComposeTokenGranter struct {
	TokenGrantDict map[string]TokenGranter
}
func NewTokenGranter(tokenGrantDict map[string]TokenGranter) TokenGranter {
	return &ComposeTokenGranter{TokenGrantDict:tokenGrantDict}
}
func (tokenGranter *ComposeTokenGranter)Grant(ctx context.Context,grantType string,client *ClientDetails,reader *http.Request) (*OAuth2Token,error)  {
	dispatchGranter := tokenGranter.TokenGrantDict[grantType]
	if dispatchGranter == nil {
		return nil,ErrNotSupportGranType
	}
	return dispatchGranter.Grant(ctx,grantType,client,reader)
}

type UsernamePasswordTokenGranter struct {
	supportGranterType string
	userDetailsService UserDetailsService
	tokenService TokenService
}
func NewUsernamePasswordTokenGranter(grantType string,service UserDetailsService,tokenService TokenService) TokenGranter {
	return &UsernamePasswordTokenGranter{
		supportGranterType: grantType,
		userDetailsService: service,
		tokenService:       tokenService,
	}
}
func (tokenGranter *UsernamePasswordTokenGranter)Grant(ctx context.Context,grantType string,client *ClientDetails,reader *http.Request) (*OAuth2Token,error)  {
	if grantType != tokenGranter.supportGranterType {
		return nil,ErrNotSupportGranType
	}
	username := reader.FormValue("username")
	password := reader.FormValue("password")
	if username == "" || password == "" {
		return nil,ErrInvalidUsernameAndPasswordRequest
	}

	userDetails,err := tokenGranter.userDetailsService.GetUserDetailByUsername(ctx,username,password)
	if err != nil {
		return nil,err
	}
	return tokenGranter.tokenService.CreateAccessToken(&OAuth2Details{
		Client: client,
		User:   userDetails,
	})
}

type RefreshTokenGranter struct {
	supportGrantType string
	tokenService TokenService
}

func NewRefreshGranter(grantType string,service TokenService) TokenGranter {
	return &RefreshTokenGranter{
		supportGrantType: grantType,
		tokenService:     service,
	}
}
func (tokenGranter *RefreshTokenGranter)Grant(_ context.Context,grantType string,_ *ClientDetails,reader *http.Request) (*OAuth2Token,error)  {
	if grantType != tokenGranter.supportGrantType {
		return nil,ErrNotSupportGranType
	}
	refreshTokenValue := reader.URL.Query().Get("refresh_token")
	if refreshTokenValue == "" {
		return nil,ErrInvalidTokenRequest
	}
	return  tokenGranter.tokenService.RefreshAccessToken(refreshTokenValue)
}

//token服务 = token存储+token加密/解密
type TokenService interface {
	// 根据访问令牌获取对应的用户信息和客户端信息
	GetOAuth2DetailsByAccessToken(tokenValue string) (*OAuth2Details, error)
	// 根据用户信息和客户端信息生成访问令牌
	CreateAccessToken(oauth2Details *OAuth2Details) (*OAuth2Token, error)
	// 根据刷新令牌获取访问令牌
	RefreshAccessToken(refreshTokenValue string) (*OAuth2Token, error)
	// 根据用户信息和客户端信息获取已生成访问令牌
	GetAccessToken(details *OAuth2Details) (*OAuth2Token, error)
	// 根据访问令牌值获取访问令牌结构体
	ReadAccessToken(tokenValue string) (*OAuth2Token, error)
}
type DefaultTokenService struct {
	tokenStore TokenStore
	tokenEnhancer TokenEnhancer
}

func NewTokenService(store TokenStore,enhancer TokenEnhancer) TokenService {
	return &DefaultTokenService{
		tokenStore:    store,
		tokenEnhancer: enhancer,
	}
}

func (tokenService *DefaultTokenService)GetOAuth2DetailsByAccessToken(tokenValue string) (*OAuth2Details, error)  {
	token,detail,err := tokenService.tokenEnhancer.Extract(tokenValue)
	if err != nil {
		return nil,err
	}
	if token != nil && token.IsExpired() {
		return nil,ErrExpiredToken
	}
	if detail == nil {
		return nil,ErrInvalidTokenRequest
	}
	return detail,nil
}
func (tokenService *DefaultTokenService)CreateAccessToken(oauth2Details *OAuth2Details) (*OAuth2Token, error)  {
	token, err := tokenService.tokenStore.GetAccessToken(oauth2Details)
	var refreshToken *OAuth2Token
	if token != nil && err == nil {
		//存在未失效令牌,直接返回
		if !token.IsExpired() {
			tokenService.tokenStore.StoreAccessToken(token,oauth2Details)
			return token,nil
		}

		//令牌已失效,移除
		tokenService.tokenStore.RemoveAccessToken(token.TokenValue)
		if token.RefreshToken != nil {
			refreshToken = token.RefreshToken
		//	tokenService.tokenStore.RemoveRefreshToken(refreshToken.TokenValue)
		}
	}

	//refreshToken没有或者过期了,则重新生成token
	if refreshToken == nil || refreshToken.IsExpired() {
		refreshToken,err = tokenService.createRefreshToken(oauth2Details)
		if err != nil {
			return nil,err
		}
	}

	accessToken, err := tokenService.createAccessToken(refreshToken,oauth2Details)
	if err != nil {
		return nil,err
	}
	tokenService.tokenStore.StoreAccessToken(accessToken,oauth2Details)
	//tokenService.tokenStore.StoreRefreshToken(refreshToken,oauth2Details)
	return accessToken,nil
}
func (tokenService *DefaultTokenService)createAccessToken(refreshToken *OAuth2Token,details *OAuth2Details) (*OAuth2Token,error) {
	validitySeconds := details.Client.AccessTokenValiditySeconds
	s, _       := time.ParseDuration(fmt.Sprintf("%ds",validitySeconds))
	expireTime := time.Now().Add(s)
	accessToken := &OAuth2Token{
		RefreshToken: refreshToken,
		ExpiresTime:  &expireTime,
		TokenValue:   uuid.NewV4().String(),
	}
	if tokenService.tokenEnhancer != nil {
		return tokenService.tokenEnhancer.Enhance(accessToken,details)
	}
	return accessToken,nil
}
func (tokenService *DefaultTokenService)createRefreshToken(details *OAuth2Details) (*OAuth2Token,error) {
	validitySeconds := details.Client.AccessTokenValiditySeconds
	s, _       := time.ParseDuration(fmt.Sprintf("%ds",validitySeconds))
	expireTime := time.Now().Add(s)
	refreshToken := &OAuth2Token{
		ExpiresTime:&expireTime,
		TokenValue:uuid.NewV4().String(),
	}
	if tokenService.tokenEnhancer != nil {
		return tokenService.tokenEnhancer.Enhance(refreshToken,details)
	}
	return refreshToken,nil
}
func (tokenService *DefaultTokenService)RefreshAccessToken(refreshTokenValue string) (*OAuth2Token, error)  {
	//解密获取token
	refreshToken, err := tokenService.tokenStore.ReadRefreshToken(refreshTokenValue)
	if err != nil {
		return nil,err
	}
	if refreshToken.IsExpired() {
		return nil,ErrExpiredToken
	}

	//解密获取tokenDetails
	oauth2Details, err := tokenService.tokenStore.ReadOAuth2DetailsForReFreshToken(refreshTokenValue)
	if err != nil {
		return nil,err
	}
	oauth2Token, err := tokenService.tokenStore.GetAccessToken(oauth2Details)
	if err != nil {
		return nil,err
	}
	//移除当前的tokenVal跟refreshTokenVal
	tokenService.tokenStore.RemoveAccessToken(oauth2Token.TokenValue)
	tokenService.tokenStore.RemoveRefreshToken(refreshTokenValue)

	//重新生成
	refreshToken, err = tokenService.createRefreshToken(oauth2Details)
	if err != nil {
		return nil,err
	}
	accessToken,err := tokenService.createAccessToken(refreshToken,oauth2Details)
	if err != nil {
		return nil,err
	}

	tokenService.tokenStore.StoreAccessToken(accessToken,oauth2Details)
	tokenService.tokenStore.StoreRefreshToken(refreshToken,oauth2Details)
	return accessToken,nil
}
func (tokenService *DefaultTokenService)GetAccessToken(details *OAuth2Details) (*OAuth2Token, error)  {

}
func (tokenService *DefaultTokenService)ReadAccessToken(tokenValue string) (*OAuth2Token, error)  {

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
type JwtTokenStore struct {
	jwtTokenEnhancer *JwtTokenEnhancer
}
type tokenInfo struct {
	Details  OAuth2Details `json:"details"`
	OldToken OAuth2Token   `json:"old_token"`
	NewToken OAuth2Token   `json:"new_token"`
}
func NewJwtTokenStore(enhancer *JwtTokenEnhancer) TokenStore {
	return &JwtTokenStore{jwtTokenEnhancer:enhancer}
}
func (tokenStore *JwtTokenStore)StoreAccessToken(oauth2Token *OAuth2Token, oauth2Details *OAuth2Details)  {
	token,_ := tokenStore.jwtTokenEnhancer.Enhance(oauth2Token,oauth2Details)

	//todo 保存到redis
	k  := genAccessTokenRedisKey(oauth2Details.Client.ClientId,oauth2Details.User.UserId)
	tf := tokenInfo{
		Details :*oauth2Details,
		OldToken:*oauth2Token,
		NewToken:*token,
	}
	v1,_ := json.Marshal(tf)
	v    := string(v1)
	fmt.Println(k,v)
}
func (tokenStore *JwtTokenStore)ReadAccessToken(tokenValue string) (*OAuth2Token,error)  {
	token,_,err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return token,err
}
func (tokenStore *JwtTokenStore)ReadAuth2Details(tokenValue string) (*OAuth2Details,error)  {
	_,details,err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return details,err
}
func (tokenStore *JwtTokenStore)GetAccessToken(oauth2Details *OAuth2Details) (*OAuth2Token,error)  {
	//todo 从redis获取,然后做json.Unmarsha1,拿不到说明没有授权|没有该客户端没有创建tokenVal
	k  := fmt.Sprintf("%s:%d",oauth2Details.Client.ClientId,oauth2Details.User.UserId)
	tf := tokenInfo{}
	fmt.Println(k,tf)
	return &tf.NewToken,ErrNotSupportOperation
}
func (tokenStore *JwtTokenStore)RemoveAccessToken(tokenValue string)  {
	details,err := tokenStore.ReadAuth2Details(tokenValue)
	if err != nil {
		return
	}

	//todo 从redis删除key
	k  := genAccessTokenRedisKey(details.Client.ClientId,details.User.UserId)
	fmt.Println(k)
}
func (tokenStore *JwtTokenStore)StoreRefreshToken(oauth2Token *OAuth2Token, oauth2Details *OAuth2Details)  {
	token,_ := tokenStore.jwtTokenEnhancer.Enhance(oauth2Token,oauth2Details)

	//todo 保存到redis
	k  := genRefreshTokenRedisKey(oauth2Details.Client.ClientId,oauth2Details.User.UserId)
	tf := tokenInfo{
		Details :*oauth2Details,
		OldToken:*oauth2Token,
		NewToken:*token,
	}
	v1,_ := json.Marshal(tf)
	v    := string(v1)
	fmt.Println(k,v)
}
func (tokenStore *JwtTokenStore)RemoveRefreshToken(oauth2TokenVal string)  {
	details,err := tokenStore.ReadAuth2Details(oauth2TokenVal)
	if err != nil {
		return
	}

	//todo 从redis删除key
	k  := genRefreshTokenRedisKey(details.Client.ClientId,details.User.UserId)
	fmt.Println(k)
}
func (tokenStore *JwtTokenStore)ReadRefreshToken(tokenValue string) (*OAuth2Token,error)  {
	token,err :=  tokenStore.ReadAccessToken(tokenValue)
	if err != nil {
		return nil,err
	}
	return token.RefreshToken,nil
}

func (tokenStore *JwtTokenStore)ReadOAuth2DetailsForReFreshToken(tokenValue string) (*OAuth2Details, error)  {
	return tokenStore.ReadAuth2Details(tokenValue)
}
func genAccessTokenRedisKey(clientId string,userId int64) string {
	return fmt.Sprintf("%s:%d:accessToken",clientId,userId)
}
func genRefreshTokenRedisKey(clientId string,userId int64) string {
	return fmt.Sprintf("%s:%d:refreshToken",clientId,userId)
}


//token加密/解密
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
	expiresTime := time.Unix(claims.ExpiresAt,0)
	return &OAuth2Token{
		RefreshToken:&claims.RefreshToken,
		TokenValue:tokenValue,
		ExpiresTime:&expiresTime,
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