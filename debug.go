// +build ignore

package skipset

import (
	"fmt"
	"unsafe"
)

func nodeInfo(n *int64Node) string {
	if n == nil {
		return "nil"
	}
	return fmt.Sprintf("address->%+v value: %+v", (*uint)(unsafe.Pointer(n)), n)
}
