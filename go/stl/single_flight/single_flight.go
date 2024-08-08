package singleFlight

import "sync"

type call struct {
	wg  sync.WaitGroup //可重入锁
	val interface{}    //请求结果
	err error          //错误反馈
}

type Group struct {
	m  map[string]*call //一类请求与同一类呼叫的映射表
	mu sync.Mutex       //并发控制锁,保证线程安全
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	//判断以key为关键词的该类请求是否存在
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		// 如果请求正在进行中，则等待
		c.wg.Wait()
		return c.val, c.err
	}
	//该类请求不存在,创建个请求
	c := new(call)
	// 发起请求前加锁,并将请求添加到请求组内以表示该类请求正在处理
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	//调用请求函数获取内容
	c.val, c.err = fn()
	//请求结束
	c.wg.Done()
	g.mu.Lock()
	//从请求组中删除该呼叫请求
	delete(g.m, key)
	g.mu.Unlock()
	return c.val, c.err
}

func (g *Group) DoChan(key string, fn func() (interface{}, error)) (ch chan interface{}) {
	ch = make(chan interface{}, 1)
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if _, ok := g.m[key]; ok {
		g.mu.Unlock()
		return ch
	}
	c := new(call)
	c.wg.Add(1)  // 发起请求前加锁
	g.m[key] = c // 添加到 g.m，表明 key 已经有对应的请求在处理
	g.mu.Unlock()
	go func() {
		c.val, c.err = fn() // 调用 fn，发起请求
		c.wg.Done()         // 请求结束
		g.mu.Lock()
		delete(g.m, key) // 更新 g.m
		ch <- c.val
		g.mu.Unlock()
	}()
	return ch
}

func (g *Group) ForgetUnshared(key string) {
	g.mu.Lock()
	_, ok := g.m[key]
	if ok {
		delete(g.m, key)
	}
	g.mu.Unlock()
}
