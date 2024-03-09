Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	// type=nil, val=nil != type=*customError val=""
	// fmt.Println(reflect.TypeOf(err)) = *main.customError
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
В конструкции if выполняется если интерфейс не пустой, в нашем случае он не пустой, 
чтобы он был пустым надо что и тип и указатель были nil, а после функции у нас 
стал nil указатель, а тип - *customError
```
