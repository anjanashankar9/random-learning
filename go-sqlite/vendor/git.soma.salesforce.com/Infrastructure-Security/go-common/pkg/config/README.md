# config

The `config` package provides a common way for reading properties files.

Using the following sample config:

```properties
    foo.bar=test1
xrd.foo.bar=test2
    foo.baz=test3
iad.foo.baz=test4
```

```go
cfg := config.New("myServiceName", os.Hostname(), "path/to/config/", "xrd")

fmt.Println(cfg.GetString("foo.bar")) // Outputs test2 because an xrd override exists
fmt.Println(cfg.GetString("foo.baz")) // Outputs test3 because an xrd override does not exist
```