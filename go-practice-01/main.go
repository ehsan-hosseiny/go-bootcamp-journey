package main

import "fmt"

type Employee struct {
	Name   string
	Salary int
}

type Contractor struct {
	HourlyRate  int
	HoursWorked int
}


type Payroll interface {
	Pay() int
}

func (e Employee) Pay() int {
	return e.Salary
}

func (c Contractor) Pay() int {
	return c.HourlyRate * c.HoursWorked
}

func ProcessPayment(p Payroll) {
	fmt.Printf("Paying amount: %d\n", p.Pay())
}

func (e *Employee) Rais(amount int){
	e.Salary += amount
}

func main() {
	employee := Employee{
		Name:   "Saeed",
		Salary: 3000,
	}

	contractor := Contractor{
		HourlyRate:  50,
		HoursWorked: 40,
	}

	employee.Rais(500)

	ProcessPayment(employee)
	ProcessPayment(contractor)
}