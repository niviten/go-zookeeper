package main

import (
	"fmt"
	"zk_engine/internal/zookeeper"
)

func main() {
	zk := zookeeper.New("127.0.0.1:2181")
	defer zk.Close()

	fmt.Printf("zk connected: %t\n", zk.IsConnected())

	fmt.Println("___done")
}
