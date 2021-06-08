package cache

var defaultCache *MultiLevelCache

func GetObject(key string, obj interface{}, ttl int, f LoadFunc) error {
	return defaultCache.GetObject(key, obj, ttl, f)
}

func Delete(key string) error {
	return defaultCache.Delete(key)
}

func NewDefaultCache(option ...Option) *MultiLevelCache {
	if defaultCache == nil {
		defaultCache = NewMultiLevelCache(option...)
	}
	return defaultCache
}

func RegisterCache(cache ...Cache) {
	if defaultCache == nil {
		defaultCache = NewDefaultCache()
	}
	defaultCache.RegisterCache(cache...)
}
