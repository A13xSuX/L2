package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err) //интерфэйс содержит тип *os.PathError
	fmt.Printf("%T\n", err)
	fmt.Println(err == nil) //type is not nil
}

//интерфейс считается nil, если 2 компонента == nil(dynamic type,dinamic value)
//обычный интерфейс имеет набор методов
//пустой интерфейс не имеет методов, может хранить значения любого типа
