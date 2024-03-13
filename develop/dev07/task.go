package main

import (
	"fmt"
	"sync"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

func main() {
	sig := func(after time.Duration) <-chan interface{} { //реализуем функцию из условия задания
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-or(
		sig(2*time.Second), //время уменьшили
		sig(5*time.Second),
		sig(1*time.Second),
		sig(1*time.Second),
		sig(1*time.Second),
	)
	fmt.Printf("Завершил работу спустя %v\n", time.Since(start)) //должен вывести время самого долгого выполнения,
	//то есть 5 сек, так как ждем выполнения всех горутин
}

func or(channels ...<-chan interface{}) <-chan interface{} {
	c := make(chan interface{})
	wg := &sync.WaitGroup{}
	for _, ch := range channels { // итерируемся по каналам
		wg.Add(1)
		go func(ch <-chan interface{}, wg *sync.WaitGroup) { // горутины блокируются пока из закрытых каналов не пойдут zero value
			defer wg.Done()
			for range ch {
				<-ch
			}
		}(ch, wg)
	}
	wg.Wait() // ждем пока закроются все каналы
	close(c)  // закрываем канал
	return c
}
