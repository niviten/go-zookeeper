package zookeeper

/*
#cgo LDFLAGS: -L/usr/local/lib -lzookeeper_mt
#include <zookeeper/zookeeper.h>

void zk_set_logger(const char* time) {
    char log_file_path[51];
    sprintf(log_file_path, "logs/zookeeper_%s.log", time);
    FILE *logFile = fopen(log_file_path, "a");
    if (logFile == NULL) {
        perror("zk set logger failed");
        return;
    }
    zoo_set_log_stream(logFile);
}

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
        timeStr := time.Now().Format("2006-01-02-15h-04m-05s")
        cTimeStr := C.CString(timeStr)
        defer C.free(unsafe.Pointer(cTimeStr))
        C.zk_set_logger(cTimeStr)
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
		time.Sleep(51 * time.Millisecond)
	}

	return &ZooKeeper{
		zh:        zh,
		connected: ctx == 1,
	}
}
