package micro

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/client"
)

type ClientWrapper struct {
	client.Client
}

func (c *ClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		//正常执行逻辑
		fmt.Println(req.Service() + "." + req.Endpoint())
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(err error) error {
		fmt.Println(err)
		return err
	})
}

// NewClientHystrixWrapper 熔断器
func NewClientHystrixWrapper() client.Wrapper {
	return func(i client.Client) client.Client {
		return &ClientWrapper{i}
	}
}
