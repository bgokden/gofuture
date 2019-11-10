# go-future

Go-future gives an implementation similar to Java/Scala Futures.

Although there are many ways to handle this behaviour in Golang.
This library is useful for people who got used to Java/Scala Future implementation.


#### Import:
```golang
import gofuture "github.com/bgokden/gofuture"
```

#### Usage:

```golang
future := gofuture.FutureFunc(func() int {
  time.Sleep(5 * time.Second)
  return x * 10
})
// do something else here
// get result when needed
result := future.Get()
```

Also it is possible to use timeouts on Get
```golang
result := future.GetWithTimeout(3 * time.Second)
```

#### Note:
This is a very basic implementation where only Get and GetWithTimeout functions are implemented.

Future Get returns an interface so type casting should be done by user.

Every future creates a channel but it is not closed so it is better allow garbage collection of Future after usage.



Java Futures: https://docs.oracle.com/javase/8/docs/api/index.html?java/util/concurrent/Future.html

Scala Futures: https://docs.scala-lang.org/overviews/core/futures.html
