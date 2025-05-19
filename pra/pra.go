package main

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
)

func main() {
	// 创建一个节点，参数是机器ID
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 生成ID
	id := node.Generate()
	fmt.Println(id.Int64())
	// 打印ID
	fmt.Printf("ID: %d\n", id)
}
