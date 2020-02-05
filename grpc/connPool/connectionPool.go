package connPool

import (
	"context"
	"github.com/redresseur/flogging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"sync"
	"time"
)

var (
	logger *flogging.FabricLogger = flogging.MustGetLogger("grpc.connPool")
)

const (
	// 一条链接同时最多允许1000个请求同时占用
	MaxUsingLimit = 1000
	ConnTimeOut   = 10 * time.Second
	MaxFreeConn   = 10
	MaxConn       = 256
	Balance       = 10 * time.Second
	Unavailable   = uint32(1<<32 - 1)
)

func SetupLogger(flogger *flogging.FabricLogger) {
	logger = flogger
}

type AlgorithmGrpcClientPool interface {
	PickClient() (*grpc.ClientConn, error)
	BackClient(conn *grpc.ClientConn)
	Release()
}

// 一条链接最多占用数量为1000
type algorithmGrpcClientPool struct {
	serverAddr  string
	ctx         context.Context
	connTimeOut time.Duration
	target      string

	pool map[*grpc.ClientConn]uint32
	l    sync.Mutex

	maxSize     uint32
	maxFreeSize uint32
	update      chan bool

	maxUsingLimit uint32

	connFuc         func(ctx context.Context, target string) (*grpc.ClientConn, error)
	balanceDuration time.Duration
}

func (ap *algorithmGrpcClientPool) Release() {
	ap.ctx.Value(ap).(context.CancelFunc)()

	// 关闭所有剩余链接
	ap.l.Lock()
	for c := range ap.pool {
		c.Close()
	}
	ap.pool = map[*grpc.ClientConn]uint32{}
	ap.l.Unlock()
}

// 执行balance的时间间隔
func WithBalanceDuration(duration time.Duration) func(*algorithmGrpcClientPool) {
	return func(pool *algorithmGrpcClientPool) {
		pool.balanceDuration = duration
	}
}

// 在使用默认连接函数时，连接超时时间
func WithConnectionTime(duration time.Duration) func(*algorithmGrpcClientPool) {
	return func(pool *algorithmGrpcClientPool) {
		pool.connTimeOut = duration
	}
}

// 替换默认的连接函数
func WithConnectionFunc(connFuc func(ctx context.Context, target string) (*grpc.ClientConn, error)) func(*algorithmGrpcClientPool) {
	return func(pool *algorithmGrpcClientPool) {
		pool.connFuc = connFuc
	}
}

// 连接池最大的连接数
func WithPoolMaxSize(maxSize uint32) func(*algorithmGrpcClientPool) {
	return func(pool *algorithmGrpcClientPool) {
		pool.maxSize = maxSize
	}
}

// 连接池空闲时保留的最大连接数
func WithPoolMaxFreeSize(maxFreeSize uint32) func(*algorithmGrpcClientPool) {
	return func(pool *algorithmGrpcClientPool) {
		pool.maxFreeSize = maxFreeSize
	}
}

// 同一条连接的最大占用数量
func WithUsingLimit(usingLimit uint32) func(*algorithmGrpcClientPool) {
	return func(pool *algorithmGrpcClientPool) {
		if usingLimit == Unavailable {
			pool.maxUsingLimit = Unavailable - 1
		}
	}
}

func NewConnectionsPool(ctx context.Context, target string, ops ...func(*algorithmGrpcClientPool)) AlgorithmGrpcClientPool {
	res := &algorithmGrpcClientPool{
		pool:            map[*grpc.ClientConn]uint32{},
		maxUsingLimit:   MaxUsingLimit,
		connTimeOut:     ConnTimeOut,
		target:          target,
		maxFreeSize:     MaxFreeConn,
		update:          make(chan bool),
		maxSize:         MaxConn,
		balanceDuration: Balance,
	}

	cctx, cancel := context.WithCancel(ctx)
	res.ctx = context.WithValue(cctx, res, cancel)

	res.connFuc = func(ctx context.Context, target string) (*grpc.ClientConn, error) {
		cCtx, cancel := context.WithTimeout(ctx, res.connTimeOut)
		defer cancel()
		return grpc.DialContext(cCtx, target, grpc.WithBlock(), grpc.WithInsecure())
	}

	for _, op := range ops {
		op(res)
	}

	go res.balance()
	return res
}

// 维持链接在均衡水准
func (ap *algorithmGrpcClientPool) balance() {
	for run := true; run; {
		timer := time.After(ap.balanceDuration)
		select {
		case <-ap.ctx.Done():
			run = false
		case <-ap.update:
			{
			}
		case <-timer:
			{
				// 检查错误链接
				ap.l.Lock()
				for c, cu := range ap.pool {
					if cu == Unavailable {
						switch c.GetState() {
						case connectivity.Shutdown:
							{
								delete(ap.pool, c)
								c.Close()
								c = nil
							}
						case connectivity.Ready:
							{
								ap.pool[c] = ap.maxUsingLimit
							}
						default:
							// TODO: 处理其他状态
						}
					}
				}
				ap.l.Unlock()

				// 去关闭未使用的空闲链接
				ap.l.Lock()
				poolSize := len(ap.pool)
				if uint32(poolSize) < ap.maxFreeSize {
					ap.l.Unlock()
					continue
				}

				// 统计空闲的链接
				freeConnSets := []*grpc.ClientConn{}
				for c, cu := range ap.pool {
					if cu >= ap.maxUsingLimit {
						freeConnSets = append(freeConnSets, c)
					}
				}

				// 释放空闲连接，最多保留 maxFreeSize 个
				freeSize := uint32(len(freeConnSets))
				if freeSize > ap.maxFreeSize {
					freeSize = freeSize - ap.maxFreeSize
					for i := uint32(0); i < freeSize; i++ {
						delete(ap.pool, freeConnSets[i])
						freeConnSets[i].Close()
					}
				}
				ap.l.Unlock()
			}
		}
	}

	return
}

func (ap *algorithmGrpcClientPool) BackClient(conn *grpc.ClientConn) {
	ap.l.Lock()
	defer ap.l.Unlock()

	switch conn.GetState() {
	case connectivity.TransientFailure:
		{
			// 当链接处于错误状态时，允许其自动重连
			// 暂时设置为不可用状态
			ap.pool[conn] = Unavailable
		}
	case connectivity.Shutdown:
		{
			delete(ap.pool, conn)
			conn.Close()
			conn = nil
		}
	default:
		{
			if ap.pool[conn] < ap.maxUsingLimit {
				ap.pool[conn] += 1
			} else {
				ap.pool[conn] = ap.maxUsingLimit
			}
		}
	}

	return
}

func (ap *algorithmGrpcClientPool) poolSize() int {
	ap.l.Lock()
	defer ap.l.Unlock()
	return len(ap.pool)
}

// 检查是否达到连接数上限
func (ap *algorithmGrpcClientPool) checkMaxLimit() bool {
	ap.l.Lock()
	defer ap.l.Unlock()
	return len(ap.pool) >= int(ap.maxSize)
}

func (ap *algorithmGrpcClientPool) updateSignal() {
	go func() {
		ap.update <- true
	}()
	return
}

func (ap *algorithmGrpcClientPool) PickClient() (conn *grpc.ClientConn, err error) {
	ap.l.Lock()
	cus := uint32(0)
	for cc, cu := range ap.pool {
		// 如果不可用则跳过
		if cu == Unavailable {
			continue
		}

		if cu > cus {
			cus = cu
			conn = cc
		}
	}
	if conn != nil {
		ap.pool[conn] -= 1
	}
	ap.l.Unlock()

	if conn != nil {
		return
	}

	// 判斷當前連接數是否達到上限
	if !ap.checkMaxLimit() {
		// 建立一条新的连接
		if conn, err = ap.connFuc(ap.ctx, ap.target); err != nil {
			logger.Errorf("DialWithHttpTarget err: %v", err)
			return
		} else {
			ap.l.Lock()
			ap.pool[conn] = ap.maxUsingLimit - 1
			ap.l.Unlock()

			// 更新狀態
			ap.updateSignal()
		}

	} else {
		time.Sleep(1 * time.Second)
		return ap.PickClient()
	}

	return
}
