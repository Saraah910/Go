package main

// jobs := make(chan int, 100)
// 	results := make(chan int, 100)

// 	go worker(jobs, results)
// 	go worker(jobs, results)
// 	go worker(jobs, results)
// 	go worker(jobs, results)
// 	go worker(jobs, results)
// 	go worker(jobs, results)

// 	for i := 0; i < 100; i++ {
// 		jobs <- i
// 	}
// 	close(jobs)

// 	for j := 0; j < 100; j++ {
// 		fmt.Println(<-results)
// 	}

// // }

// func fibonacci(n int) int {
// 	if n <= 1 {
// 		return 1
// 	}
// 	return fibonacci(n-1) + fibonacci(n-2)
// }
// func worker(jobs <-chan int, results chan<- int) {
// 	for n := range jobs {
// 		results <- fibonacci(n)
// 	}
// }
// func greet(name string, c chan string) {
// 	for i := 1; i <= 5; i++ {
// 		c <- name
// 		time.Sleep(time.Millisecond * 500)
// 	}
// 	close(c)
// }

// c1 := make(chan string)
// c2 := make(chan string)
// go func() {
// 	for {
// 		c1 <- "every 500 miliSec"
// 		time.Sleep(time.Millisecond * 500)
// 	}
// }()

// go func() {
// 	for {
// 		c2 <- "Every 2 sec"
// 		time.Sleep(time.Second * 2)
// 	}
// }()

// for {
// 	select {
// 	case msg1 := <-c1:
// 		fmt.Println(msg1)
// 	case msg2 := <-c2:
// 		fmt.Println(msg2)
// 	}
// }

// c := make(chan string)
// go greet("sakshi", c)

// for msg := range c {
// 	fmt.Println(msg)
// }
// ************************
// var wg sync.WaitGroup
// wg.Add(1)
// go func() {
// 	greet("Sakshi")
// 	wg.Done()
// }()
// wg.Add(1)
// go func() {
// 	greet("Praj")
// 	wg.Done()
// }()
// wg.Wait()
// go greet("Praj")
// // time.Sleep(time.Second * 10)
// fmt.Scanln()
