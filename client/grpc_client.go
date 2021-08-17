package client

import (
	"context"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
	"sync"

	"github.com/irisnet/core-sdk-go/types"
)

var clientConnSingleton *grpc.ClientConn
var once sync.Once

type grpcClient struct {
}

// Token token
type Token struct {
	projectId        string
	projectKey       string
	chainAccountAddr string
}

const (
	projectIdHeader           = "projectId"
	projectKeyHeader          = "projectKey"
	chainAccountAddressHeader = "chainAccountAddress"
)

// GetRequestMetadata 获取当前请求认证所需的元数据
func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{projectIdHeader: t.projectId, projectKeyHeader: t.projectKey, chainAccountAddressHeader: t.chainAccountAddr}, nil
}

// RequireTransportSecurity 是否需要基于 TLS 认证进行安全传输
func (t *Token) RequireTransportSecurity() bool {
	return true
}

func NewGRPCClient(url string, info types.BSNProjectInfo) types.GRPCClient {
	once.Do(func() {

		token := Token{
			projectId:        info.ProjectId,
			projectKey:       info.ProjectKey,
			chainAccountAddr: info.ChainAccountAddress,
		}

		dialOpts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithPerRPCCredentials(&token),
		}
		clientConn, err := grpc.Dial(url, dialOpts...)
		if err != nil {
			log.Error(err.Error())
			panic(err)
		}
		clientConnSingleton = clientConn
	})
	return &grpcClient{}
}

func (g grpcClient) GenConn() (*grpc.ClientConn, error) {
	return clientConnSingleton, nil
}
