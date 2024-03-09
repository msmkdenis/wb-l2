package pattern

import "fmt"

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

/*
	Стратегия определяет семейство похожих алгоритмов и свой объект для каждого из них.
	Это позволяет клиенту выбирать подходящий алгоритм на этапе выполнения кода.
	Например, клиенту может быть предоставлен итератор по графу и возможность выбрать стратегию обхода:
		в глубину (Depth-First Search) — рекурсивный перебор всех дочерних элементов;
		в ширину (Breadth-First Search) — последовательный перебор всех элементов на расстоянии k, затем k + 1 и т. д.

	Клиент освобождён от деталей реализации алгоритмов.
	Это помогает улучшать алгоритмы в объекте Стратегия или добавлять новые, не требуя изменений клиентского кода.
	С другой стороны, клиент должен знать, чем различаются существующие алгоритмы, чтобы выбрать подходящий вариант.

	Паттерн Стратегия используется, когда:
		нужно применять разные варианты одного и того же алгоритма;
		нужно выбирать алгоритм во время выполнения программы;
		нужно скрывать детали реализации алгоритмов.

	Классическая формулировка паттерна включает три типа участников:
		объект strategy умеет работать с алгоритмами семейства;
		объекты concreteStrategy содержат частные реализации этих алгоритмов;
		объект context взаимодействует с клиентским кодом, формирует стратегический запрос.

	Преимущества
		Можно динамически определять, какой алгоритм будет запущен
		Соблюдается инкапсуляция - код алгоритмов отделен и скрыт от остального кода,
		Алгоритмы вызываются единообразно: без if-ов и других подобных конструкций
	
	Недостатки:
		Усложняет программу за счёт дополнительных классов.
		Клиент должен знать, в чём состоит разница между стратегиями, чтобы выбрать подходящую.

	Пример реализации шаблона с кешированием данных в оперативной памяти.
	Поскольку размер памяти ограничен, нужно обеспечить механизм ротации кеша.
	При освобождении памяти могут применяться стратегии:
		LRU (Least Recently Used) — вычищаем элементы, которые использовались давно.
		FIFO (First In First Out) — удаляем элементы, которые были созданы раньше остальных.
		LFU (Least Frequently Used) — чистим записи, которые использовались редко.
*/

// evictionAlgo — интерфейс strategy.
type evictionAlgo interface {
    evict(c *cache)
}

// реализация concreteStrategy

type fifo struct {}

func (l *fifo) evict(c *cache) {
    fmt.Println("Evicting by fifo strategy")
}

type lru struct {}

func (l *lru) evict(c *cache) {
    fmt.Println("Evicting by lru strategy")
}

type lfu struct {}

func (l *lfu) evict(c *cache) {
    fmt.Println("Evicting by lfu strategy")
}

// cache содержит контекст.
type cache struct {
    storage      map[string]string
    evictionAlgo evictionAlgo
    capacity     int
    maxCapacity  int
}

func initCache(e evictionAlgo) *cache {
    storage := make(map[string]string)
    return &cache{
        storage:      storage,
        evictionAlgo: e,
        capacity:     0,
        maxCapacity:  2,
    }
}

// setEvictionAlgo определяет алгоритм освобождения памяти.
func (c *cache) setEvictionAlgo(e evictionAlgo) {
    c.evictionAlgo = e
}

func (c *cache) add(key, value string) {
    if c.capacity == c.maxCapacity {
        c.evict()
    }
    c.capacity++
    c.storage[key] = value
}

func (c *cache) get(key string) {
    delete(c.storage, key)
}

func (c *cache) evict() {
    c.evictionAlgo.evict(c)
    c.capacity--
}

// func main() {
//     // клиентский код
//     lfu := &lfu{}
//     cache := initCache(lfu)
//     cache.add("a", "1")
//     cache.add("b", "2")
//     cache.add("c", "3")
//     lru := &lru{}
//     cache.setEvictionAlgo(lru)
//     cache.add("d", "4")
//     fifo := &fifo{}
//     cache.setEvictionAlgo(fifo)
//     cache.add("e", "5")
// } 