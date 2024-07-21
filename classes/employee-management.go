package classes

import (
	"fmt"
	"strings"
)

type Employee struct {
	id           int
	name         string
	managerId    int
	subordinates []*Employee
}

type System struct {
	employees   []*Employee
	employeeMap map[int]*Employee
}

func NewEmployee(name string) *Employee {
	employee := &Employee{
		id:        getUniqueId(),
		name:      name,
		managerId: -1,
	}
	return employee
}

func getUniqueId() int {
	staticId++
	return staticId
}

var staticId = 0

func (e *Employee) GetId() int {
	return e.id
}

func (e *Employee) GetManagerId() int {
	return e.managerId
}

func (e *Employee) SetManagerId(managerId int) {
	e.managerId = managerId
}

func (e *Employee) GetName() string {
	return e.name
}

func (e *Employee) GetSubordinates() []*Employee {
	return e.subordinates
}

func (e *Employee) AddSubordinate(subordinate *Employee) {
	e.subordinates = append(e.subordinates, subordinate)
}

func NewSystem() *System {
	return &System{
		employeeMap: make(map[int]*Employee),
	}
}

func (s *System) RegisterEmployee(employee *Employee) {
	s.employees = append(s.employees, employee)
	s.employeeMap[employee.GetId()] = employee
}

func (s *System) RegisterManager(empId int, managerId int) {
	employee, empOk := s.employeeMap[empId]
	manager, mgrOk := s.employeeMap[managerId]
	if !empOk || !mgrOk {
		fmt.Println("Either Employee or Manager not registered! Please provide correct registered identifiers to continue")
		return
	}

	employee.SetManagerId(managerId)
	manager.AddSubordinate(employee)
}

func (s *System) PrintDetails(empId int) {
	employee, ok := s.employeeMap[empId]
	if !ok {
		fmt.Println("Employee not registered! Please provide correct Id and retry")
		return
	}

	fmt.Printf("Id: %d\tName: %s\t", empId, employee.GetName())
	if employee.GetManagerId() != 0 {
		manager := s.employeeMap[employee.GetManagerId()]
		fmt.Printf("Manager: %s\n", manager.GetName())
	} else {
		fmt.Println("Manager: None")
	}
}

func (s *System) PrintDetailsByPrefix(prefix string) {
	for _, employee := range s.employees {
		if strings.HasPrefix(employee.GetName(), prefix) {
			fmt.Printf("Id: %d\tName: %s\t", employee.GetId(), employee.GetName())
			if employee.GetManagerId() != 0 {
				manager := s.employeeMap[employee.GetManagerId()]
				fmt.Printf("Manager: %s\n", manager.GetName())
			} else {
				fmt.Println("Manager: None")
			}
		}
	}
}

func (s *System) GetSubordinates(empId int) []*Employee {
	employee, ok := s.employeeMap[empId]
	if !ok {
		fmt.Println("Employee not registered! Please provide correct Id and retry")
		return nil
	}
	return employee.GetSubordinates()
}

func (s *System) GetSubordinatesByName(name string) []*Employee {
	for _, employee := range s.employees {
		if employee.GetName() == name {
			return employee.GetSubordinates()
		}
	}
	return nil
}

func EmployeeManagement() {
	employee := NewEmployee("Achilles")
	employee1 := NewEmployee("Hector")
	employee2 := NewEmployee("Paris")
	employee3 := NewEmployee("Helen")

	system := NewSystem()
	system.RegisterEmployee(employee)
	system.RegisterEmployee(employee1)
	system.RegisterEmployee(employee2)
	system.RegisterEmployee(employee3)

	system.RegisterManager(employee1.GetId(), employee.GetId())
	system.RegisterManager(employee2.GetId(), employee.GetId())
	system.RegisterManager(employee3.GetId(), employee.GetId())

	system.PrintDetails(employee1.GetId())
	fmt.Println("********************************************************************")

	system.PrintDetailsByPrefix("He")
	fmt.Println("********************************************************************")

	subordinates := system.GetSubordinatesByName(employee.GetName())
	for _, e := range subordinates {
		fmt.Printf("%s %d\n", e.GetName(), e.GetId())
	}
	fmt.Println("********************************************************************")

	subordinates = system.GetSubordinates(employee.GetId())
	for _, e := range subordinates {
		fmt.Printf("%s %d\n", e.GetName(), e.GetId())
	}
	fmt.Println("********************************************************************")

	subordinates = system.GetSubordinates(employee1.GetId())
	for _, e := range subordinates {
		fmt.Printf("%s %d\n", e.GetName(), e.GetId())
	}
}
