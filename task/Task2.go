package task

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

func Task2() {
	// 指针
	num := 5
	fmt.Println("Before add:", num)
	add(&num)
	fmt.Println("After add:", num)
	// 数组
	arrs := []int{1, 2, 3, 4, 5}
	add2(&arrs)
	fmt.Println("After add2:", arrs)
	// go routine
	goRoutine()
	// channel
	channel()
	//
	task1 := func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("任务1完成")
	}
	task2 := func() {
		time.Sleep(700 * time.Millisecond)
		fmt.Println("任务2完成")
	}
	task3 := func() {
		time.Sleep(300 * time.Millisecond)
		fmt.Println("任务3完成")
	}
	scheduler := newScheduler(task1, task2, task3)
	scheduler.run()

	// 接口
	interfaceTest()

	// 结构体
	printInfoTest()

	// 带缓冲channel
	channelTest()

	//sync.Mutex锁计数器
	incrementTest()

}

func add(a *int) {
	*a += 10
}

func add2(arrs *[]int) {
	for i, arr := range *arrs {
		(*arrs)[i] = arr * 2
	}
}

func goRoutine() {
	go func() {
		for i := 1; i < 10; i = i + 2 {
			fmt.Println("线程1", i)
		}
	}()

	go func() {
		for i := 0; i <= 10; i = i + 2 {
			fmt.Println("线程2", i)
		}
	}()
	time.Sleep(1 * time.Second)
}

func send(ch chan int, value int) {
	fmt.Println("发送数据:", value)
	ch <- value
}

func receive(ch chan int) {
	fmt.Println("等待数据...")
	time.Sleep(1 * time.Second)
	value := <-ch
	fmt.Println("接收到数据:", value)
}

func channel() {
	ch := make(chan int)
	go send(ch, 88)
	go receive(ch)
	time.Sleep(2 * time.Second)
}

type TaskFunc func()

type Scheduler struct {
	tasks []TaskFunc
}

func newScheduler(tasks ...TaskFunc) *Scheduler {
	return &Scheduler{tasks: tasks}
}

func executeTask(task TaskFunc) time.Duration {
	statrTime := time.Now()
	task()
	return time.Since(statrTime)
}

func (s *Scheduler) run() {
	var wg = sync.WaitGroup{}
	ch := make(chan int, len(s.tasks))

	for _, task := range s.tasks {
		wg.Add(1)
		go func(t TaskFunc) {
			defer wg.Done()
			duration := executeTask(t)
			ch <- int(duration.Milliseconds())
		}(task)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for duration := range ch {
		fmt.Println("任务执行完成:", duration, "ms")
	}

}

type Shape interface {
	Area()
	Perimeter()
}

type Rectangle struct {
	Name   string
	Length int
	Width  int
}

type Circle struct {
	Name   string
	Radius float64
}

func (r *Rectangle) Area() {
	s := r.Length * r.Width
	fmt.Println(r.Name+"的面积:", s)
}

func (c *Circle) Area() {
	area := math.Pi * c.Radius * c.Radius
	fmt.Println(c.Name+"的面积:", math.Round(area*100)/100)
}
func (r *Rectangle) Perimeter() {
	fmt.Println(r.Name+"的周长", 2*(r.Length+r.Width))
}

func (c *Circle) Perimeter() {
	perimeter := math.Pi * c.Radius
	fmt.Println(c.Name+"的周长", math.Round(perimeter*100)/100)
}

func interfaceTest() {
	var shape Shape
	shape = &Rectangle{Name: "矩形", Length: 10, Width: 20}
	shape.Area()
	shape.Perimeter()

	shape = &Circle{Name: "圆形", Radius: 5}
	shape.Area()
	shape.Perimeter()
}

type Person struct {
	Name string
	Age  int
}

type Staff struct {
	Person
	StaffID string
}

func (e *Staff) PrintInfo() {
	fmt.Println("员工ID:", e.StaffID, "员工姓名:", e.Name, "员工年龄:", e.Age)
}

func printInfoTest() {
	person := &Person{Name: "张三", Age: 18}
	employee := &Staff{Person: *person, StaffID: "E001"}
	employee.PrintInfo()
}

func channelTest() {
	nums := 100
	ch := make(chan int, nums)
	for i := 0; i < 20; i++ {
		go receiveBatch(ch)
	}
	go sendBatch(ch, nums)
	time.Sleep(2 * time.Second)
	close(ch)
}

func sendBatch(ch chan int, nums int) {
	for i := 1; i <= nums; i++ {
		ch <- i
	}
}

func receiveBatch(ch chan int) {
	for value := range ch {
		fmt.Println("接收到的值为:", value)
	}
}

var (
	count       = 0
	countAtomic int32
	lock        sync.Mutex
)

func increment(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		lock.Lock()
		{
			temp := count
			temp = temp + 1
			count = temp
		}
		lock.Unlock()
	}
}

func incrementAtomic() {
	for i := 0; i < 1000; i++ {
		atomic.AddInt32(&countAtomic, 1)
	}
}

func incrementTest() {

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go increment(&wg)
	}
	wg.Wait()
	fmt.Println("count:", count)

	for i := 0; i < 10; i++ {
		go incrementAtomic()
	}
	time.Sleep(1 * time.Second)
	fmt.Println("countAtomic:", countAtomic)
}
