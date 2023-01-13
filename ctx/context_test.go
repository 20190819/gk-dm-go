package ctx

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContextTimeout(t *testing.T) {

	ctx1, cancel1 := context.WithTimeout(context.Background(), 1*time.Second)
	ctx2, cancel2 := context.WithTimeout(ctx1, 3*time.Second)
	defer cancel2()
	defer cancel1()

	go func() {
		<-ctx1.Done()
		fmt.Println("timeout 1")
	}()
	go func() {
		<-ctx2.Done()
		fmt.Println("timeout 2")
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("end")
}

func TestTimeoutExample(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var signal chan struct{}
	signal = make(chan struct{})

	go func() {
		doBiz()
		signal <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("progress timeout")
	case <-signal:
		fmt.Println("biz ok")
	}
}

func doBiz() {
	time.Sleep(2 * time.Second)
}

func TestTimeoutTimeAfter(t *testing.T) {
	signal := make(chan struct{})
	go func() {
		doBiz()
		signal <- struct{}{}
	}()

	timer := time.AfterFunc(time.Second, func() {
		fmt.Println("timeout")
	})
	<-signal
	fmt.Println("biz ok")
	timer.Stop()
}
