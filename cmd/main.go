package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	data := make([]int, 1_000_000)
	for _, i := range data {
		data[i] = i
	}

	ch := make(chan int)
	defer close(ch)
	parts := 10
	size := len(data) / parts

	for i := 0; i < parts; i++ {
		go func(ch chan<- int, data []int) {
			sum := 0
			for _, v := range data {
				sum += v
			}
			ch <- sum
		}(ch, data[i*size:(i+1)*size])
	}
	total := 0
	for i := 0; i < parts; i++ {
		total += <-ch
	}
	log.Print(total)
}

func sum() {
	data := make([]int, 1_000_00)
	for i := range data {
		data[i] = i
	}
	parts := 10
	size := len(data) / parts
	chs := make([]<-chan int, parts)
	for i := 0; i < parts; i++ {
		ch := make(chan int)
		chs[i] = ch
		go func(ch chan<- int, data []int) {
			defer close(ch)
			sum := 0
			for _, v := range data {
				sum += v
			}
			ch <- sum
		}(ch, data[i*size:(i+1)*size])
	}
	total := 0
	for value := range merge(chs) {
		total += value
	}
	log.Print(total)
}

func merge(chs []<-chan int) <-chan int {
	wg := sync.WaitGroup{}
	wg.Add(len(chs))
	merged := make(chan int)

	for _, ch := range chs {
		go func(ch <-chan int) {
			defer wg.Done()
			for val := range ch {
				merged <- val
			}
		}(ch)
	}
	go func() {
		defer close(merged)
		wg.Wait()
	}()
	return merged
}

func tick(done <-chan struct{}) {
	for {
		select {
		case <-done:
			return
		case <-time.After(time.Second):
			log.Print("tick")
		}
	}
}
