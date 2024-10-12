package zookeeper

/*
#cgo LDFLAGS: -L/usr/local/lib -lzookeeper_mt
#include <string.h>
#include <zookeeper/zookeeper.h>

typedef struct {
    int count;
    char** children;
    int done;
} GetChildrenResult;

GetChildrenResult* GetChildrenResult_create() {
    GetChildrenResult* result = (GetChildrenResult*) malloc(sizeof(GetChildrenResult));
    result->count = 0;
    result->done = 0;
    return result;
}

int GetChildrenResult_get_done(GetChildrenResult* result) {
    return result->done;
}

int GetChildrenResult_get_count(GetChildrenResult* result) {
    return result->count;
}

char* GetChildrenResult_get_child_at_index(GetChildrenResult* result, int index) {
    return result->children[index];
}

void get_children_cb(int rc, const struct String_vector* children, void* data) {
    int i;
    GetChildrenResult* result;

    result = (GetChildrenResult*) data;

    if (rc == ZOK) {
        result->count = children->count;
        result->children = (char**) malloc(children->count * sizeof(char*));
        for (i = 0; i < children->count; i++) {
            result->children[i] = (char*) malloc(strlen(children->data[i]) * sizeof(char));
            strcpy(result->children[i], children->data[i]);
        }
        result->done = 1;
    } else {
        result->done = -1;
    }
}
*/
import "C"
import (
	"errors"
	"time"
	"unsafe"
)

func (zk *ZooKeeper) GetChildren(path string) ([]string, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

    result := C.GetChildrenResult_create()

	r := C.zoo_aget_children(zk.zh, cPath, 0, (*[0]byte)(C.get_children_cb), unsafe.Pointer(result))
    _ = r

    for {
        done := C.GetChildrenResult_get_done(result)
        if done == 1 || done == -1 {
            break
        }
        time.Sleep(51 * time.Millisecond)
    }

    done := C.GetChildrenResult_get_done(result)
    if done == -1 {
        return []string{}, errors.New("Zookeeper Get Children failed")
    }

    count := int(C.GetChildrenResult_get_count(result))

    children := make([]string, count)
    for i := 0; i < count; i++ {
        children[i] = C.GoString(C.GetChildrenResult_get_child_at_index(result, C.int(i)))
    }

    return children, nil
}
