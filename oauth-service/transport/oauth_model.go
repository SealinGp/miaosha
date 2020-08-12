package transport

import (
	"context"
	oauthServiceEndpoint "miaosha/oauth-service/endpoint"
	"miaosha/oauth-service/model"
	"miaosha/pb"
)

//req->pbReq
func EncodeGRPCCheckTokenRequest(_ context.Context,r interface{}) (interface{},error) {
	req := r.(*oauthServiceEndpoint.CheckTokenRequest)
	return &pb.CheckTokenRequest{
		Token:req.Token,
	},nil
}
//pbReq->req
func DecodeGRPCCHeckTokenRequest(_ context.Context,r interface{}) (interface{},error) {
	req := r.(*pb.CheckTokenRequest)
	return &oauthServiceEndpoint.CheckTokenRequest{
		Token: req.Token,
	},nil
}

//resp->pbResp
func EncodeGRPCCheckTokenResponse(_ context.Context,r interface{}) (interface{},error) {
	resp := r.(oauthServiceEndpoint.CheckTokenResponse)
	if resp.Error != "" {
		return &pb.CheckTokenResponse{
			IsValidToken:  false,
			Err:           resp.Error,
		},nil
	}

	return &pb.CheckTokenResponse{
		UserDetails:   &pb.UserDetails{
			UserId:      resp.OAuth2Details.User.UserId,
			Username:    resp.OAuth2Details.User.Password,
			Authorities: resp.OAuth2Details.User.Authorities,
		},
		ClientDetails: &pb.ClientDetails{
			ClientId:                    resp.OAuth2Details.Client.ClientId,
			AccessTokenValiditySeconds:  int32(resp.OAuth2Details.Client.AccessTokenValiditySeconds),
			RefreshTokenValiditySeconds: int32(resp.OAuth2Details.Client.RefreshTokenValiditySeconds),
			AuthorizedGrantTypes:        resp.OAuth2Details.Client.AuthorizedGrantTypes,
		},
		IsValidToken:  false,
		Err:           "",
	},nil
}

//pbResp->resp
func DecodeGRPCCheckTokenResponse(_ context.Context,r interface{}) (interface{},error) {
	resp := r.(*pb.CheckTokenResponse)
	if resp.Err != "" {
		return oauthServiceEndpoint.CheckTokenResponse{
			OAuth2Details: nil,
			Error:         resp.Err,
		},nil
	}
	return oauthServiceEndpoint.CheckTokenResponse{
		OAuth2Details: &model.OAuth2Details{
			Client: &model.ClientDetails{
				ClientId:                    resp.GetClientDetails().ClientId,
				AccessTokenValiditySeconds:  int(resp.GetClientDetails().AccessTokenValiditySeconds),
				RefreshTokenValiditySeconds: int(resp.GetClientDetails().RefreshTokenValiditySeconds),
				AuthorizedGrantTypes:        resp.GetClientDetails().GetAuthorizedGrantTypes(),
			},
			User:   &model.UserDetails{
				UserId:      resp.UserDetails.UserId,
				Username:    resp.UserDetails.Username,
				Authorities: resp.UserDetails.Authorities,
			},
		},
		Error:         "",
	},nil
}