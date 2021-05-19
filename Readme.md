# memorycache 

This package provides to save any value in memory 

**Important!**
Be careful to use this cache module in different processes when creating group of microservices
The cache table is only valid in one process!

### Usage:
```go
cache := New(10*time.Minute, 1*time.Hour)
cache.Set("simple_key", "value", 1*time.Minute)
...
v := cache.Get("simple_key")
...
cache.Delete("simple_key")
```
