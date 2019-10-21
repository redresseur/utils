package connPool

import (
	"context"
	"google.golang.org/grpc"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	destTarget    = "127.0.0.1:8080"
	benchMarkPool = NewConnectionsPool(context.Background(), destTarget,
		WithUsingLimit(10),
		WithBalanceDuration(1*time.Second),
		WithConnectionFunc(func(ctx context.Context, target string) (*grpc.ClientConn, error) {
			return grpc.Dial(destTarget, grpc.WithInsecure())
		}))
)

func TestNewConnectionsPool(t *testing.T) {
	connPool := NewConnectionsPool(context.Background(), destTarget, WithConnectionFunc(func(ctx context.Context, target string) (*grpc.ClientConn, error) {
		return &grpc.ClientConn{}, nil
	}))

	connPool.PickClient()
}

func TestAlgorithmGrpcClientPool_PickClient(t *testing.T) {
	connPool := NewConnectionsPool(context.Background(), destTarget,
		WithUsingLimit(10),
		WithBalanceDuration(1*time.Second),
		WithConnectionFunc(func(ctx context.Context, target string) (*grpc.ClientConn, error) {
			return grpc.Dial(destTarget, grpc.WithInsecure())
		}))

	g := &sync.WaitGroup{}
	for i := 0; i < 1024; i++ {
		g.Add(1)
		go func() {
			defer g.Done()
			c, err := connPool.PickClient()
			if err != nil {
				panic(err)
			}

			time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)*100))
			connPool.BackClient(c)
		}()
	}

	g.Wait()
	t.Logf("%+v", connPool)

	time.Sleep(2 * time.Second)
	t.Logf("%+v", connPool)
}

// 测试超出上限时的情况
func TestAlgorithmGrpcClientPool_PickClient2(t *testing.T) {
	connPool := NewConnectionsPool(context.Background(), destTarget,
		WithBalanceDuration(1*time.Second),
		WithPoolMaxSize(5),
		WithPoolMaxFreeSize(5),
		WithUsingLimit(1),
		WithConnectionFunc(func(ctx context.Context, target string) (*grpc.ClientConn, error) {
			return grpc.Dial(destTarget, grpc.WithInsecure())
		}))

	g := &sync.WaitGroup{}
	for i := 0; i < 1024; i++ {
		g.Add(1)
		go func() {
			defer g.Done()
			c, err := connPool.PickClient()
			if err != nil {
				panic(err)
			}

			connPool.BackClient(c)
		}()
	}

	g.Wait()
	t.Logf("%+v", connPool)

	time.Sleep(2 * time.Second)
	t.Logf("%+v", connPool)
}

func BenchmarkAlgorithmGrpcClientPool_PickClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, _ := benchMarkPool.PickClient()
		benchMarkPool.BackClient(conn)
	}
}
