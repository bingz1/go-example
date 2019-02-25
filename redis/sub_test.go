package redis

import (
	"fmt"
	"testing"
)

func TestSub(t *testing.T) {
	ch := make(chan int, 1)
	go Sub(ch)
	<- ch
	fmt.Println("over")
}
