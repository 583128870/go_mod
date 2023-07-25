package cache

import (
	"errors"
	"log"
	"time"
)

type TestData struct {
	Key  string `json:"key"`
	Date string `json:"date"`
}

//testStructModel 结构体json字符串类型的缓存
type testStructModel struct {
}

func (t *testStructModel) GetCacheKey(key interface{}, data interface{}) (CacheKey, error) {
	if info, ok := data.(TestData); ok {
		return CacheKey(info.Key), nil
	} else {
		return "", errors.New("类型错误")
	}
}

//EliminatePolicyFunc 淘汰策略
func (t *testStructModel) EliminatePolicyFunc(key CacheKey, data interface{}) bool {
	dataTime, err := time.ParseInLocation(time.DateTime, string(key), time.Local)
	if err != nil {
		log.Printf("时间转换错误：%s", err.Error())
		return true
	}
	//淘汰20秒前的数据
	if dataTime.Add(1 * time.Second).Before(time.Now()) {
		return true
	} else {
		return false
	}
}

//PollUpdatePolicyFunc 巡检更新策略 参数：oldList 之前的缓存数据  响应：需要缓存的新数据（差异数据）
func (t *testStructModel) PollUpdatePolicyFunc(oldList interface{}) interface{} {

	nowList := make(map[string]TestData)
	key := time.Now().Format(time.DateTime)
	nowList[key] = TestData{
		Key:  key,
		Date: key,
	}
	if _oldList, ok := oldList.(map[string]interface{}); ok {
		if _oldList[key] == nil {
			return nowList
		}
	}
	return nowList
}

//字符串类型的缓存
type testStringModel struct {
}

func (t *testStringModel) GetCacheKey(key interface{}, data interface{}) (CacheKey, error) {
	if info, ok := key.(string); ok {
		return CacheKey(info), nil
	} else {
		return "", errors.New("类型错误")
	}
}

//EliminatePolicyFunc 淘汰策略
func (t *testStringModel) EliminatePolicyFunc(key CacheKey, data interface{}) bool {
	dataTime, err := time.ParseInLocation(time.DateTime, string(key), time.Local)
	if err != nil {
		log.Printf("时间转换错误：%s", err.Error())
		return true
	}
	//淘汰20秒前的数据
	if dataTime.Add(1 * time.Second).Before(time.Now()) {
		return true
	} else {
		return false
	}
}

//PollUpdatePolicyFunc 巡检更新策略 参数：oldList 之前的缓存数据  响应：需要缓存的新数据（差异数据）
func (t *testStringModel) PollUpdatePolicyFunc(oldList interface{}) interface{} {
	nowList := make(map[string]string)
	key := time.Now().Format(time.DateTime)
	nowList[key] = "123123"
	if _oldList, ok := oldList.(map[string]interface{}); ok {
		if _oldList[key] == nil {
			return nowList
		}
	}
	return nowList
}
