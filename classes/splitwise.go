package classes

import (
	"fmt"
	"math"
	"sync"
)

type Split int

const (
	EQUAL Split = iota
	EXACT
	PERCENT
)

type User struct {
	ID   int
	Name string
	// total expense so far means the total amount the user has spent so far ( means he has to give this amount to others) if positive
	// and if negative, it means the user has to get back this amount from others
	TotalExpenseSoFar float64
	UserExpenseSheet  map[*User]float64
	mu                sync.Mutex
}

var userIDCounter int = 1
var userIDLock sync.Mutex

func NewUser(name string) *User {
	userIDLock.Lock()
	id := userIDCounter
	userIDCounter++
	userIDLock.Unlock()
	return &User{
		ID:               id,
		Name:             name,
		UserExpenseSheet: make(map[*User]float64),
	}
}

func (u *User) AddToUserExpenseSheet(user *User, value float64) {
	if u == user {
		return
	}

	u.mu.Lock()
	defer u.mu.Unlock()
	// adding value to user's total expense so far and to the user's expense sheet
	u.TotalExpenseSoFar += value
	u.UserExpenseSheet[user] += value
}

func (u *User) PrintTotalBalance() {
	if u.TotalExpenseSoFar > 0 {
		fmt.Printf("%s owes a total of %.2f\n", u.Name, u.TotalExpenseSoFar)
	} else {
		fmt.Printf("%s gets back a total of %.2f\n", u.Name, -u.TotalExpenseSoFar)
	}
}

type Expense struct {
	ID                  int
	Description         string
	Split               Split
	PercentDistribution []float64
	ExactDistribution   []float64
	Creditor            *User
	Defaulters          []*User
	ExactTotalAmount    float64
}

var expenseIDCounter int = 1
var expenseIDLock sync.Mutex

func NewExpense(creditor *User, split Split, defaulters []*User, exactTotalAmount float64) *Expense {
	expenseIDLock.Lock()
	id := expenseIDCounter
	expenseIDCounter++
	expenseIDLock.Unlock()
	return &Expense{
		ID:               id,
		Creditor:         creditor,
		Split:            split,
		Defaulters:       defaulters,
		ExactTotalAmount: exactTotalAmount,
	}
}

type Splitwise struct {
	users     []*User
	userIDMap map[int]*User
	mu        sync.Mutex
}

func NewSplitwise() *Splitwise {
	return &Splitwise{
		userIDMap: make(map[int]*User),
	}
}

func (s *Splitwise) RegisterUser(user *User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.userIDMap[user.ID]; !exists {
		s.users = append(s.users, user)
		s.userIDMap[user.ID] = user
	}
}

func (s *Splitwise) AddExpense(expense *Expense) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.verifyUsers(expense.Creditor, expense.Defaulters) {
		fmt.Println("Can't process expense. Kindly register all users and retry")
		return
	}
	s.calculateExpenses(expense)
}

func (s *Splitwise) verifyUsers(user *User, users []*User) bool {
	usersMap := make(map[int]bool)
	for _, usr := range users {
		usersMap[usr.ID] = true
	}
	usersMap[user.ID] = true

	for id := range usersMap {
		if _, exists := s.userIDMap[id]; !exists {
			return false
		}
	}
	return true
}

func (s *Splitwise) calculateExpenses(expense *Expense) bool {
	creditor := expense.Creditor
	defaulters := expense.Defaulters
	var amtPerHead []float64

	switch expense.Split {
	case EQUAL:
		amtPerHead = divideEqually(expense.ExactTotalAmount, len(defaulters))
		for i, defaulter := range defaulters {
			// creditor gets amtPerHead[i] from defaulter (negative means creditor gets)
			s.userIDMap[creditor.ID].AddToUserExpenseSheet(defaulter, -amtPerHead[i])
			// defaulter owes the creditor amtPerHead[i] (positive means defaulter owes)
			s.userIDMap[defaulter.ID].AddToUserExpenseSheet(creditor, amtPerHead[i])
		}
	case EXACT:
		amtPerHead = expense.ExactDistribution
		if expense.ExactTotalAmount != sum(amtPerHead) {
			fmt.Println("Can't create expense. Total amount doesn't equal sum of individual amounts. Please try again!")
			return false
		}

		if len(amtPerHead) != len(defaulters) {
			fmt.Println("The amounts and value numbers don't match. Expense can't be created. Please try again!")
			return false
		}
		for i, defaulter := range defaulters {
			s.userIDMap[creditor.ID].AddToUserExpenseSheet(defaulter, -amtPerHead[i])
			s.userIDMap[defaulter.ID].AddToUserExpenseSheet(creditor, amtPerHead[i])
		}
	case PERCENT:
		amtPerHead = expense.PercentDistribution
		if sum(amtPerHead) != 100 {
			fmt.Println("Can't create expense. All percentages don't add to 100. Please try again!")
			return false
		}

		if len(amtPerHead) != len(defaulters) {
			fmt.Println("The percents and value numbers don't match. Expense can't be created. Please try again!")
			return false
		}
		for i, defaulter := range defaulters {
			// why we do below calculation?
			amount := (amtPerHead[i] * expense.ExactTotalAmount) / 100.0
			amount = math.Floor((amount*100.0)+0.5) / 100.0
			s.userIDMap[creditor.ID].AddToUserExpenseSheet(defaulter, -amount)
			s.userIDMap[defaulter.ID].AddToUserExpenseSheet(creditor, amount)
		}
	default:
		break
	}
	return true
}

func (s *Splitwise) PrintBalanceForAllUsers() {
	for _, user := range s.users {
		user.PrintTotalBalance()
	}
}

func divideEqually(totalAmount float64, memberCount int) []float64 {
	parts := make([]float64, memberCount)
	for i := 0; i < memberCount; i++ {
		part := math.Trunc((100.0*totalAmount)/float64(memberCount-i)) / 100.0
		parts[i] = part
		totalAmount -= part
	}
	return parts
}

func sum(arr []float64) float64 {
	total := 0.0
	for _, value := range arr {
		total += value
	}
	return total
}

// leave it for now
func (s *Splitwise) SimplifyExpenses() {
	amounts := make([]int, len(s.users))
	for i, user := range s.users {
		amounts[i] = int(-user.TotalExpenseSoFar * 100)
	}

	for {
		minIdx, maxIdx := minMaxIdx(amounts)
		if amounts[minIdx] == 0 && amounts[maxIdx] == 0 {
			break
		}

		minAmount := min(-amounts[minIdx], amounts[maxIdx])
		amounts[minIdx] += minAmount
		amounts[maxIdx] -= minAmount

		fmt.Printf("%s pays the amount %.2f to %s\n",
			s.users[minIdx].Name, float64(minAmount)/100.0, s.users[maxIdx].Name)
	}
}

func minMaxIdx(arr []int) (minIdx, maxIdx int) {
	minIdx = 0
	maxIdx = 0
	for i := range arr {
		if arr[i] < arr[minIdx] {
			minIdx = i
		}
		if arr[i] > arr[maxIdx] {
			maxIdx = i
		}
	}
	return
}

func SplitwiseExpense() {
	u1 := NewUser("Jitu")
	u2 := NewUser("Navin")
	u3 := NewUser("Yogi")
	u4 := NewUser("Mandal")

	users := []*User{u1, u2, u3, u4}

	sp := NewSplitwise()
	sp.RegisterUser(u1)
	sp.RegisterUser(u2)
	sp.RegisterUser(u3)
	sp.RegisterUser(u4)

	expense := NewExpense(u1, EQUAL, users, 2000)
	sp.AddExpense(expense)
	sp.PrintBalanceForAllUsers()

	users2 := []*User{u2, u3}
	expense2 := NewExpense(u1, EXACT, users2, 1400)
	expense2.ExactDistribution = []float64{500, 900}
	sp.AddExpense(expense2)
	sp.PrintBalanceForAllUsers()

	expense3 := NewExpense(u4, PERCENT, users, 1200)
	expense3.PercentDistribution = []float64{40, 20, 20, 20}
	sp.AddExpense(expense3)
	sp.PrintBalanceForAllUsers()

	fmt.Println()

	for _, user := range sp.users {
		for debtor, amount := range user.UserExpenseSheet {
			if amount > 0 {
				fmt.Printf("%s owes a total of %.2f to %s\n", user.Name, amount, debtor.Name)
			} else {
				fmt.Printf("%s gets back a total of %.2f from %s\n", user.Name, -amount, debtor.Name)
			}
		}
	}

	fmt.Println()
	sp.SimplifyExpenses()
}
