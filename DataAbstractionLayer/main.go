package main

import "mongodblayer"

type person struct {
	Name string
	Age  int
	City string
}

func main() {
	mongodblayer.Init("mongodb://localhost:27017", "golang")
	defer mongodblayer.Close()
	mongodblayer.TestConnection()
	//pepe := person{"Pepito", 24, "La Habana, Cuba"}
	mongodblayer.FindOneDocument("go", "5d12946f3666c63b396274ae")
}
