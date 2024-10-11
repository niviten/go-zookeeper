package zookeeper

/*
#cgo LDFLAGS: -L/usr/local/lib -lzookeeper_mt
#include <zookeeper/zookeeper.h>

void watcher(zhandle_t *zzh, int type, int state, const char *path, void *watcherCtx) {
    int *ctx = (int*) watcherCtx;
    if (state == ZOO_CONNECTED_STATE) {
        printf("Zookeeper connected successfully\n");
        *ctx = 1;
    } else {
        printf("Zookeeper connection failed\n");
        *ctx = -1;
    }
}
*/
import "C"
import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

type ZooKeeper struct {
	zh        *C.zhandle_t
	connected bool
}

var instance *ZooKeeper
var once sync.Once

func New(host string) *ZooKeeper {
	once.Do(func() {
		instance = new(host)
	})
	return instance
}

func (zk *ZooKeeper) Close() {
	if zk.connected {
		C.zookeeper_close(zk.zh)
		zk.connected = false
	}
}

func (zk *ZooKeeper) IsConnected() bool {
	return zk.connected
}

func new(host string) *ZooKeeper {
	cHost := C.CString(host)
	defer C.free(unsafe.Pointer(cHost))

	ctx := C.int(0)

	zh := C.zookeeper_init(cHost, (*[0]byte)(C.watcher), 30000, nil, unsafe.Pointer(&ctx), 0)

	for {
		if ctx == 1 || ctx == -1 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Println("ctx: ", int(ctx))

	return &ZooKeeper{
		zh:        zh,
		connected: ctx == 1,
	}
}
