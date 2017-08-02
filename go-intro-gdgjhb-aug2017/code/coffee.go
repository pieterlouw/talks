package main

import "fmt"

//Coffee struct is public scoped
type Coffee struct {
	Name        string //public
	strongLevel int    //private
}

// Drink is a func
func (c *Coffee) Drink() {
	fmt.Printf("%s hits the spot\n", c.Name)
}

func (c *Coffee) SetStrongLevel(level int) {
	c.strongLevel = level
}

func (c *Coffee) StrongLevel() int {
	return c.strongLevel
}

func DeterminCoffeeStrongness(c Coffee) error {
	if c.StrongLevel() < 5 { //don't drink anything less than 5
		return fmt.Errorf("this coffee (%s) is not strong enough", c.Name)
	}

	return nil
}

func main() {
	var frisco = Coffee{}
	var number42 int

	ricoffy, grounded := Coffee{Name: "Ricoffy"}, Coffee{Name: "Grounded coffee"}

	number42, number1 := 42, 1

	grounded.SetStrongLevel(5)
	ricoffy.SetStrongLevel(1)

	err := DeterminCoffeeStrongness(ricoffy)
	if err != nil {
		fmt.Printf("Cannot drink coffee: %+v\n", err)
	} else {
		ricoffy.Drink()
	}

	err = DeterminCoffeeStrongness(grounded)
	if err != nil {
		fmt.Printf("Cannot drink coffee: %+v\n", err)
	} else {
		grounded.Drink()
	}

}
