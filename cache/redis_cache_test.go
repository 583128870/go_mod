package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"testing"
	"time"
)

func TestRedisCache(t *testing.T) {
	redisConn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Fatalf("redis连接错误：%s", err.Error())
	}
	statistics := NewRedisCache(&testStringModel{}, redisConn, "test")
	//value, rel, err := statistics.GetCacheInfo("5")
	//if err != nil {
	//	fmt.Println(err)
	//}
	for {
		time.Sleep(1 * time.Second)
		allData, err := statistics.GetCacheData(false)
		if err != nil {
			log.Printf("获取缓存数据错误：%s", err.Error())
			continue
		}
		fmt.Println(allData)
	}

}
