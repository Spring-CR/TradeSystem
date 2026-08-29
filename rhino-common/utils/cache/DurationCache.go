package cache

import (
	"sync"
	"time"
	"log"
)

type DurationCache struct{
	sync.RWMutex
	liftCyc time.Duration
	onClearValue func(key string, value interface{})
	mapData map[string]interface{}
}

func NewDurationCache(liftCyc time.Duration, onClearValue func(key string, value interface{}))*DurationCache{
	return &DurationCache{liftCyc:liftCyc, onClearValue:onClearValue, mapData:make(map[string]interface{})}
}

func (c *DurationCache)Put(key string, value interface{}){
	c.Lock()
	c.mapData[key] = value
	c.Unlock()
	go func(c *DurationCache, k string){
		time.Sleep(c.liftCyc)
		c.Lock()
		v,ok:=c.mapData[k]
		log.Printf("remove expired value with key:%v", k)
		delete(c.mapData, k)
		c.Unlock()
		if ok{
			c.onClearValue(k, v)
		}
	}(c,key)
}

func (c *DurationCache)Delete(key string){
	c.Lock()
	defer c.Unlock()

	delete(c.mapData, key)
}

func (c *DurationCache)Get(key string)(interface{}, bool){
	c.RLock()
	defer c.RUnlock()

	v, ok := c.mapData[key]
	return v, ok
}