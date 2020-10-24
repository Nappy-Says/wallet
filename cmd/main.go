package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	anotherExample()
	
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

type Producer struct {
	OutChan chan int
}

func (p *Producer) produce() {
	for {
		time.Sleep(2 * time.Second)
		p.OutChan <- rand.Int()
	}
}
func (p *Producer) getOutChan() <-chan int {
	return p.OutChan
}
func anotherExample() {
	p := Producer{OutChan: make(chan int, 10)}
	go p.produce()
	prodChan := p.getOutChan()
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case a := <-prodChan:
			log.Println("got from a: ", a)
		case s := <-ticker.C:
			log.Println("got from s", s)

		}
	}
	}
}
