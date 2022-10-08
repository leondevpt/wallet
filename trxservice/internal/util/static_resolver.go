package util

import (
	"github.com/leondevpt/wallet/trxservice/global"

	"google.golang.org/grpc/resolver"
)

// 定义 Scheme 名称
const staticScheme = global.StaticSchema

type staticResolverBuilder struct {
	addrsStore map[string][]string
}

func NewStaticResolverBuilder(addrsStore map[string][]string) *staticResolverBuilder {
	return &staticResolverBuilder{addrsStore: addrsStore}
}

func (e *staticResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// 初始化 resolver, 将 addrsStore 传递进去
	r := &staticResolver{
		target:     target,
		cc:         cc,
		addrsStore: e.addrsStore,
	}
	// 调用 start 初始化地址
	r.start()
	return r, nil
}
func (e *staticResolverBuilder) Scheme() string { return staticScheme }

type staticResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *staticResolver) start() {
	// 在静态路由表中查询此 Endpoint 对应 addrs
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	// addrs 列表转化为 state, 调用 cc.UpdateState 更新地址
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*staticResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*staticResolver) Close()                                  {}
