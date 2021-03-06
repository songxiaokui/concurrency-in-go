package main

/**
Golang默认所有的任务都在一个cpu核里，如果想使用多核来跑goroutine的任务，需要配置runtime.GOMAXPROCS。
GOMAXPROCS 的数目根据自己任务量分配就可以了，有个前提是不要大于你的cpu核数。
并行比较适合那种cpu密集型计算，如果是IO密集型使用多核的化会增加cpu切换的成本。
如果程序是ＩＯ为主的，启用多核心反而有上下文切换
所以对于以涉及ＩＯ操作的主的程序启用多核对于加速程序意义不大。
*/
import (
	"fmt"
	"runtime"
)

func test(c chan bool, n int) {
	x := 0

	var num int = 1000 * 1e6
	for i := 0; i < num; i++ {
		x += i
	}

	fmt.Println(n, x)

	//当n=10的时候就退出
	if n == 10 {
		c <- true
	}
}

func main() {
	//设置多个cpu进行运算
	fmt.Println("cpu num: ", runtime.NumCPU()) //4

	// runtime.GOMAXPROCS(3) //同时有3个goroutine在运行，分在3个不同的cpu对应的线程上进行

	runtime.GOMAXPROCS(4)

	c := make(chan bool)

	// 在独立携程中进行
	for i := 0; i < 200; i++ {
		go test(c, i)
	}

	<-c

	fmt.Println("main end...")
}

/**
$ time go run more_cpu.go
cpu num:  4
199 499999999500000000
3 499999999500000000
1 499999999500000000
4 499999999500000000
0 499999999500000000
2 499999999500000000
...

5 499999999500000000
102 499999999500000000
53 499999999500000000
116 499999999500000000
54 499999999500000000
103 499999999500000000
6 499999999500000000
55 499999999500000000
104 499999999500000000
56 499999999500000000
7 499999999500000000
105 499999999500000000
57 499999999500000000
8 499999999500000000
106 499999999500000000
58 499999999500000000
9 499999999500000000
59 499999999500000000
107 499999999500000000
10 499999999500000000
main end...

real    0m53.031s
user    2m33.619s
sys     0m0.257s
从打印的结果来看，每次有三个独立携程在运行

当加大cpu运行的核数
main end...

real    0m14.182s
user    0m52.987s
sys     0m0.149

通过top查看进程，发现more_cpu大量消耗cpu
21839 heige     20   0  102708   2216   1324 R 391.4  0.0   1:43.56 more_cpu
*/
