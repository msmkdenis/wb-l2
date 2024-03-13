package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:


- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*

Так же требуется поддерживать функционал fork/exec-команд

Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).


*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор, пока не будет введена команда выхода (например \quit).
*/

func main() {
	scanner := bufio.NewScanner(os.Stdin) //чтение со стандартного потока ввода
	fmt.Println(`print "quit" to exit`)
	for {
		currentDir, _ := os.Getwd()    //получаем текущую директорию
		fmt.Printf("%s> ", currentDir) //выводим предложение к вводу команд с указанием текущей директории
		ok := scanner.Scan()
		if !ok {
			log.Fatal(errors.New("can't read from stdin"))
		}
		input := scanner.Text() //считываем ввод пользователя
		if input == "quit" {    //если ввели строку exit - выходим из всей программы
			os.Exit(1)
		}
		commands := strings.Split(input, "|") // разделитель для использования множественных команд в одной строке
		err := process(commands)
		if err != nil {
			fmt.Println(err)
			fmt.Println("try again")
		}
	}
}

func process(commands []string) error {
	for _, command := range commands {
		a := strings.TrimPrefix(command, " ")
		a = strings.TrimSuffix(a, " ")
		c := strings.Split(a, " ")
		switch c[0] {
		case "pwd":
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			fmt.Println(wd)
		case "cd":
			if len(c) != 2 {
				return errors.New(`"cd" command must have one argument`)
			}
			err := os.Chdir(c[1])
			if err != nil {
				return err
			}
		case "echo":
			if len(c) < 2 {
				return errors.New(`not enough arguments in "echo" command`)
			}
			fmt.Println(strings.Join(c[1:], " "))
		case "kill":
			if len(c) != 2 {
				return errors.New(`"kill" command must have one argument`)
			}
			pid, err := strconv.Atoi(c[1])
			if err != nil {
				return errors.New(`pid in "kill" command must be integer`)
			}
			pr, err := os.FindProcess(pid)
			if err != nil {
				return errors.New(`wrong pid in "kill" command`)
			}
			err = pr.Kill()
			if err != nil {
				return errors.New(`failed to "kill" process`)
			}
		case "ps":
			pr, err := ps.Processes()
			if err != nil {
				return err
			}
			for _, v := range pr {
				fmt.Println(v.Pid(), v.Executable(), v.PPid())
			}
		default:
			cmd := exec.Command(c[0], strings.Join(c[1:], " "))
			err := cmd.Run()
			if err != nil {
				return err
			}
			cmd.Wait()
		}
	}
	return nil
}
