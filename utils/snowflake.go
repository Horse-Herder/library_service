package utils

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	node, _ = snowflake.NewNode(1)
}

func GetSnowFlakeId() string {
	unId := node.Generate().String()
	return unId
}

