// Этот файл с программным кодом должен принадлежать какому-нибудь пакету (в данном случае - main)
package main

// В файле может использоваться функционал из других пакетов.
// В этом случае используемые пакеты надо импортировать с помощью ключевого слова import
import (
	"fmt"
	"list/storage/list"
)

type Person struct {
	Name string
	Age  int
}

func PersonPrinter(p Person) {
	fmt.Printf("Name:\t'%s'\nAge:\t'%d'\n", p.Name, p.Age)
}

func sum(a, b int) (sum int) {
	sum = a + b
	return
}

func main() {
	l := list.List{}
	l.Add(4)
	l.Print()
}
