package cache

type RuntimeCache struct {
	values map[string][]byte
}

func (c *RuntimeCache) ReadString(key string) string {
	val, ok := c.values[key]
	if ok {
		return string(val)
	}
	return ""
}

func (c *RuntimeCache) Has(key string) bool {
	_, ok := c.values[key]
	return ok
}

func (c *RuntimeCache) EraseAll() error {
	c.values = make(map[string][]byte)
	return nil
}

func (c *RuntimeCache) Write(key string, val []byte) error {
	c.values[key] = val
	return nil
}
