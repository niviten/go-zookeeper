package zookeeper

/*
#cgo LDFLAGS: -L/usr/local/lib -lzookeeper_mt
#include <zookeeper/zookeeper.h>

void get_children_cb(int rc, const struct String_vector* children, const void* data) {
    int i;

    if (rc == ZOK) {
        printf("get children success\n");
        for (i = 0; i < children->count; i++) {
            printf("c: %s\n", children->data[i]);
        }
    } else {
        printf("get children failed\n");
    }
}
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

func (zk *ZooKeeper) GetChildren(path string) []string {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	r := C.zoo_aget_children(zk.zh, cPath, 0, (*[0]byte)(C.get_children_cb), nil)

	fmt.Println("r: ", r)

	time.Sleep(1 * time.Second)

	return []string{}
}
