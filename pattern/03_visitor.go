package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*
	Посетитель - поведенический паттерн, который позволяет отвязать функциональность от объекта.
	Новые методы добавляются не для каждого типа из семейства, а для промежуточного объекта visitor, аккумулирующего функциональность.
	Типам семейства добавляется только один метод accept(visitor).
	Так проще добавлять операции к существующей базе кода без особых изменений и страха всё сломать.
	Этот паттерн чаще всего используется, когда нужно добавить функционал к объектам разного типа.

	Паттерн Посетитель используется, когда:
		нужно применить одну и ту же операцию к объектам разных типов;
		часто добавляются новые операции для объектов;
		требуется добавить новый функционал, но избежать усложнения кода объекта.

	Преимущества:
		упрощает добавление нового функционала существующим объектам, когда требуется выполнить независимые операции над объектами сложной структуры.

	Недостатки:
		не оправда если часто изменяется иерархия объектов
		может привести к нарушению инкапсуляции объектов
*/

/*
	Предположим, есть конструктор, который собирает автомобиль из колёс и двигателя. В какой-то момент нужно добавить тестирование компонентов при сборке.
*/

// CarPart — семейство типов, которым хотим добавить
// функциональность детали автомобиля.
type CarPart interface {
    Accept(CarPartVisitor)
}

// CarPartVisitor — интерфейс visitor,
// в его коде и содержится новая функциональность.
type CarPartVisitor interface {
    testWheel(wheel *Wheel)
    testEngine(engine *Engine)
}

// Wheel — реализация деталей.
type Wheel struct {
    Name string
}

// Accept — единственный метод, который нужно добавить типам семейства,
// ссылка на метод visitor.
func (w *Wheel) Accept(visitor CarPartVisitor) {
    visitor.testWheel(w)
}

type Engine struct{}

func (e *Engine) Accept(visitor CarPartVisitor) {
    visitor.testEngine(e)
}

type Car struct {
    parts []CarPart
}

// NewCar — конструктор автомобиля.
func NewCar() *Car {
    this := new(Car)
    this.parts = []CarPart{
        &Wheel{"front left"},
        &Wheel{"front right"},
        &Wheel{"rear right"},
        &Wheel{"rear left"},
        &Engine{}}
    return this
}

func (c *Car) Accept(visitor CarPartVisitor) {
    for _, part := range c.parts {
        part.Accept(visitor)
    }
}

// TestVisitor — конкретная реализация visitor, 
// которая может проверять колёса и двигатель.
type TestVisitor struct {
}

func (v *TestVisitor) testWheel(wheel *Wheel) {
    fmt.Printf("Testing the %v wheel\n", wheel.Name)
}

func (v *TestVisitor) testEngine(engine *Engine) {
    fmt.Println("Testing engine")
}

// func main() {
//     // клиентский код
//     car := NewCar()
//     visitor := new(TestVisitor)
//     car.Accept(visitor)
// }
