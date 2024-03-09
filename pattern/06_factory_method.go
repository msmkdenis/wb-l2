package pattern

import (
	"fmt"
	"strings"
)

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

/*
	Фабрика (фабричный метод) — порождающий паттерн проектирования, который определяет интерфейс для создания объектов определённого типа.
	Паттерн предлагает создавать объекты не напрямую, а через вызов специального фабричного метода.
	Все объекты должны иметь общий интерфейс, который отвечает за их создание.
	При использовании в объектно-ориентированном программировании базовый класс делегирует создание объектов классам-наследникам.
	Так как в Go отсутствуют классы и наследование, нельзя реализовать классический вариант Фабрики.

	Паттерн используется, когда:
		заранее неизвестны типы объектов — фабричный метод отделяет код создания объектов от остального кода, где они используются;
		нужна возможность расширять части существующей системы.
*/


type Database interface {
    Insert(q string) error
	Delete(q string) error
}

type Mysql struct {
}

func newMysql(dsn string) *Mysql {
    fmt.Println("Connect to mysql")
    return &Mysql{}
}

func (c *Mysql) Insert(q string) error {
    fmt.Printf("Insert to mysql: %s\n", q)
    return nil
}

func (c *Mysql) Delete(q string) error {
    fmt.Printf("Delete from mysql: %s\n", q)
    return nil
}

type Postgresql struct {
}

func newPostgresql(dsn string) *Postgresql {
    fmt.Println("Connect to postgresql")
    return &Postgresql{}
}

func (c *Postgresql) Insert(q string) error {
    fmt.Printf("Insert to postgresql: %s\n", q)
    return nil
}

func (c *Postgresql) Delete(q string) error {
    fmt.Printf("Delete from postgresql: %s\n", q)
    return nil
}

// NewConnector реализует фабричный метод.
func NewConnector(dsn string) Database {
    switch {
    case strings.HasPrefix(dsn, "mysql://"):
        return newMysql(dsn)
    case strings.HasPrefix(dsn, "postgresql://"):
        return newPostgresql(dsn)
    default:
        panic(fmt.Sprintf("unknown dsn protocol: %s", dsn))
    }
}

// func main() {
//     mysql := NewConnector("mysql://...")
//     mysql.Insert("")
//     mysql.Delete("")

//     pg := NewConnector("postgresql://...")
//     pg.Insert("")
//     pg.Delete("")
// }