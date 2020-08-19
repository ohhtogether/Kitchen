package global

import (
	"github.com/bwmarrin/snowflake"
)

// ID ID的节点实例
var GenId *snowflake.Node

// SetIDNode 设置ID节点
func SetIDNode() error {
	var err error
	GenId, err = snowflake.NewNode(0)
	return err
}
