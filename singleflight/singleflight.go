package singleflight

import "sync"

// call 代表正在进行中，或已经结束的请求，使用sync.WaitGroup锁避免重入
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group 是singlefight的主数据结构,管理不同的key的请求(call)
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// Do 接收两个参数，第一个参数是key 第二个参数是一个函数fn 针对相同的key，无论Do被调用多少次fn只会被调用一次，等fn调用结束了，返回返回值或错误
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	//没有 创建
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	//有 直接返回
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()         //如果请求在进行中，则等待
		return c.val, c.err //请求结束，返回结果
	}
	//没有把这个key创建穿进去一系列的
	c := new(call)
	c.wg.Add(1)  //发起加锁请求
	g.m[key] = c //添加到g.m 表明key已经有对应的请求在处理
	g.mu.Unlock()

	c.val, c.err = fn() //调用fn,发起请求
	c.wg.Done()         //请求结束
	delete(g.m, key)
	g.mu.Lock()

	return c.val, c.err

}
