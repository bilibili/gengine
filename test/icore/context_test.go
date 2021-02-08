package icore

import (
	"context"
	"sync"
	"testing"
	"time"
)

func Test_conext(t *testing.T) {

	var lock sync.Mutex
	lock.Lock()
	ctx, cancal := context.WithCancel(context.Background())
	cancal()
	select {
	default:
		println("111111")
	case <-ctx.Done():
		lock.Unlock()
		println("xxxxx")
	}

	println("yyyy")
	time.Sleep(5 * time.Second)
	//ctx.Done()
	println("done")
	time.Sleep(5 * time.Second)
}
