# AOI
###调研学习实现一些AOI兴趣区算法

#### 目前实现：
- 九宫格
- 四叉树


```shell
go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: code.corp.ecoplants.tech/cloud/meta/atools/aoi
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkGridManger-12    	      30	  43034352 ns/op	18078476 B/op	  130843 allocs/op
BenchmarkQuadtree-12      	     100	  31404461 ns/op	 9761675 B/op	  112047 allocs/op
```

