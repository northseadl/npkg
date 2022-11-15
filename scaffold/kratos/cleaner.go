package kratos

// Cleaner 存储清理方法的结构体
type Cleaner struct {
	cleans []func()
}

// NewCleaner 创建一个存储清理方法的结构体
func NewCleaner() *Cleaner {
	return &Cleaner{}
}

// AddFunc 添加清理方法
func (c *Cleaner) AddFunc(f func()) *Cleaner {
	c.cleans = append(c.cleans, f)
	return c
}

func (c *Cleaner) Clean() {
	for _, clean := range c.cleans {
		clean()
	}
}
