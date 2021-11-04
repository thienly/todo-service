package main

import "fmt"

type celcius float64
type temperature struct {
	celcius
}

func (c *celcius) change()  {
	*c = 20
}
func main()  {
	c:= celcius(10.0)
	t:= temperature{c}
	t.change()
	fmt.Println(t.celcius)
}