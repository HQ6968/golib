package monitor

import (
	"time"
	"sync"
	"core"
	"github.com/verystar/golib/logger"
	"sync/atomic"
	"encoding/json"
)

type RedisMonitor struct {
	data chan map[string]interface{}
	dbPrefix   string
}

var redisMonitor *RedisMonitor
var once sync.Once

func NewRedisMonitor(dbPrefix string) *RedisMonitor {
	//Single
	once.Do(func() {
		redisMonitor = new(RedisMonitor)
		redisMonitor.data = make(chan map[string]interface{}, 100)
		redisMonitor.dbPrefix = dbPrefix
	})
	return redisMonitor
}

func (this *RedisMonitor) Run() {
	var err error
	var marshaledBytes []byte

	client, ok := core.Redis("stat")
	if !ok {
		logger.Fatal("Redis not found")
		return
	}
	var count uint32
	pipe := client.Pipeline()
	for data := range this.data {
		n := atomic.AddUint32(&count, 1)
		if marshaledBytes, err = json.Marshal(data); err != nil {
			return
		}
		pipe.RPush("__stat__", string(marshaledBytes))
		if n == 100 {
			atomic.StoreUint32(&count, 0)
			_, err = pipe.Exec()
			if err != nil {
				logger.Error("Stat Pipeline error", err)
			}
		}
	}
}

func Stat(num int64, v1, v2, v3 string) {
	if num < 0 || v1 == "" || v2 == "" || v3 == "" {
		return
	}

	data := map[string]interface{}{
		"dbf":     redisMonitor.dbPrefix,
		"num":     num,
		"v1":      v1,
		"v2":      v2,
		"v3":      v3,
		"v4":      nil,
		"replace": false,
		"time":    time.Now().Unix(),
	}

	redisMonitor.data <- data
}
