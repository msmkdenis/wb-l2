package pattern

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
	Строитель — порождающий паттерн проектирования, который позволяет пошагово создавать объекты.
	Когда создаётся сложный объект, у конструктора может быть слишком много параметров, причём часть из них используются редко.
	Вместо этого можно определить несколько методов-конструкторов, которые будут отвечать за инициализацию определённых параметров.
	Тогда сложный объект будет создаваться в несколько шагов.

	Паттерн используется, когда:
		нужно избавиться от одного сложного конструктора или целого набора;
		код должен реализовывать разные представления объекта, при этом этапы создания представлений одинаковы;
		надо пошагово создать сложный объект, в том числе с использованием рекурсии.

	Недостатки
		Усложняет код программы из-за введения дополнительных классов.
*/

// Object — объект с параметром.
type Object struct {
    // данные объекта
    // ...
    // настраиваемые поля объекта
    Mode int
    Path string
}

// SetMode — пример функции, которая присваивает поле Mode.
func (o *Object) SetMode(mode int) *Object {
    o.Mode = mode
    return o
}

// SetPath — пример функции, которая присваивает поле Path.
func (o *Object) SetPath(path string) *Object {
    o.Path = path
    return o
}

// NewObject — функция-конструктор объекта.
func NewObject() *Object {
    return &Object{}
}

// func main() {
//     o := NewObject().SetMode(10).SetPath(`root`)
// 	fmt.Println(o)
// }
