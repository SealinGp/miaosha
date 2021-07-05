package model

import (
	"encoding/json"
	"miaosha/pkg/mysql"
)

type ClientDetails struct {
	ClientId string
	//client 秘钥
	ClientSecret string
	//访问令牌有效时间,秒
	AccessTokenValiditySeconds int
	//刷新令牌有效时间,秒
	RefreshTokenValiditySeconds int
	//重定向地址,授权码类型中使用
	RegisteredRedirectUri string
	//可以使用的授权类型
	AuthorizedGrantTypes []string
}

func (clientDetails *ClientDetails)IsMatch(clientId, clientSecret string) bool {
	return clientId == clientDetails.ClientId && clientSecret == clientDetails.ClientSecret
}
type ClientDetailsModel struct {

}

func NewClientDetailsModel() *ClientDetailsModel {
	return &ClientDetailsModel{}
}

func (p *ClientDetailsModel)getTableName() string {
	return "client_details"
}

func (p *ClientDetailsModel)GetClientDetailsByClientId(clientId string) (*ClientDetails,error) {
	conn := mysql.DB()
	result, err := conn.Table(p.getTableName()).Where(map[string]interface{}{"client_id":clientId}).First()
	if err != nil {
		return nil,err
	}
	var grantTypes []string
	_ = json.Unmarshal([]byte(result["authorized_grant_types"].(string)),&grantTypes)
	return &ClientDetails{
		ClientId:clientId,
		ClientSecret:result["client_secret"].(string),
		AccessTokenValiditySeconds:int(result["access_token_validity_seconds"].(int64)),
		RefreshTokenValiditySeconds:int(result["refresh_token_validity_seconds"].(int64)),
		RegisteredRedirectUri:result["registered_redirect_uri"].(string),
		AuthorizedGrantTypes:grantTypes,
	},nil
}

func (p *ClientDetailsModel)CreateClientDetails(clientDetails *ClientDetails) error {
	conn               := mysql.DB()
	grantTypeString, _ := json.Marshal(clientDetails.AuthorizedGrantTypes)
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"client_id":clientDetails.ClientId,
		"client_secret":clientDetails.ClientSecret,
		"access_token_validity_seconds":clientDetails.AccessTokenValiditySeconds,
		"refresh_token_validity_seconds":clientDetails.RefreshTokenValiditySeconds,
		"registered_redirect_uri":clientDetails.RegisteredRedirectUri,
		"authorized_grant_types":grantTypeString,
	}).Insert()
	if err != nil {
		return err
	}
	return nil
}