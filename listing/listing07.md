Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
Сначала выведутся цифры от 1 до 9 в случайном порядке, затем 0 пока вручную не прекратим работу программы.

Происходит потому, что:
	каналы а и b отрабатывают заданное время и закрываются
	чтение из закрытого канала возможно, из закрытого канала приходят zero value значения
	конструкция for ... := range chan читает из канала с, до момента пока он не закроется
	канал с НЕ закрывается в отличии от каналов а и b

Т.е. каналы b c в итоге, т.к. закрываются, постоянно шлют zero value, а канал c (merge) бесконечно их обрабатывает
и читается через far ... range, который работает до тех пор пока !читаемый канал не будет закрыт.
```
