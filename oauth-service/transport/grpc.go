package transport

import (
	"context"
	"github.com/go-kit/kit/transport/grpc"
	oauthServiceEndpoint "miaosha/oauth-service/endpoint"
	"miaosha/pb"
)

type grpcServer struct {
	 checkTokenServer grpc.Handler
}

func (s *grpcServer)CheckToken(ctx context.Context,r *pb.CheckTokenRequest) (*pb.CheckTokenResponse,error) {
	_, resp, err := s.checkTokenServer.ServeGRPC(ctx,r)
	if err != nil {
		return nil,err
	}
	return resp.(*pb.CheckTokenResponse),nil
}
func NewGRPCServer(ctx context.Context,endpoints oauthServiceEndpoint.OAuth2Endpoints,serverTracer grpc.ServerOption) pb.OAuthServiceServer {
	return &grpcServer{
		checkTokenServer:grpc.NewServer(
			endpoints.GRPCCheckTokenEndpoint,
			DecodeGRPCCHeckTokenRequest,
			EncodeGRPCCheckTokenResponse,
			serverTracer,
		),
	}
}