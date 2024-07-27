package classes

import (
	"fmt"
	"sync"
)

type FoodItems int
type BeverageItems int

const (
	Sandwich FoodItems = iota
	Poha
	Vada
	Burger
)

const (
	Tea BeverageItems = iota
	Coffee
	Water
)

type Store struct {
	id                string
	foodSupply        map[FoodItems]int
	beverageSupply    map[BeverageItems]int
	foodUnitsSold     map[FoodItems]int
	beverageUnitsSold map[BeverageItems]int
	foodRates         map[FoodItems]int
	beverageRates     map[BeverageItems]int
	mu                sync.Mutex
}

func NewStore(foods, beverages []pair) *Store {
	s := &Store{
		id:                getUniqueID("store"),
		foodSupply:        make(map[FoodItems]int),
		beverageSupply:    make(map[BeverageItems]int),
		foodUnitsSold:     make(map[FoodItems]int),
		beverageUnitsSold: make(map[BeverageItems]int),
	}
	for _, food := range foods {
		s.foodSupply[FoodItems(food.first)] = food.second
	}
	for _, beverage := range beverages {
		s.beverageSupply[BeverageItems(beverage.first)] = beverage.second
	}
	return s
}

func (s *Store) purchaseFood(foodItem FoodItems, qty int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.foodSupply[foodItem] < qty {
		fmt.Println("Not enough stocks")
		return
	}

	fmt.Println("Purchasing ...")
	fmt.Printf("[Before Purchase] foodUnitsSold[foodItem]: %d\n", s.foodUnitsSold[foodItem])

	s.foodSupply[foodItem] -= qty
	s.foodUnitsSold[foodItem] += qty

	fmt.Printf("[After Purchase] foodUnitsSold[foodItem]: %d\n", s.foodUnitsSold[foodItem])
}

func (s *Store) purchaseBeverage(beverageItem BeverageItems, qty int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.beverageSupply[beverageItem] < qty {
		fmt.Println("Not enough stocks")
		return
	}
	fmt.Println("Purchasing ...")
	fmt.Printf("[Before Purchase] beverageUnitsSold[foodItem]: %d\n", s.beverageUnitsSold[beverageItem])
	s.beverageSupply[beverageItem] -= qty
	s.beverageUnitsSold[beverageItem] += qty
	fmt.Printf("[After Purchase] beverageUnitsSold[foodItem]: %d\n", s.beverageUnitsSold[beverageItem])
}

func (s *Store) setBeverageRates(beverageRates map[BeverageItems]int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.beverageRates = beverageRates
}

func (s *Store) setFoodRates(foodRates map[FoodItems]int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.foodRates = foodRates
}

func (s *Store) getFoodUnitsSold() map[FoodItems]int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.foodUnitsSold
}

func (s *Store) getId() string {
	return s.id
}

// city records the prices of food and beverages
type City struct {
	id             string
	foodPrices     map[FoodItems]int
	beveragePrices map[BeverageItems]int
	stores         []*Store
	mu             sync.Mutex
}

func NewCity(foodPrices, beveragePrices []pair) *City {
	c := &City{
		id:             getUniqueID("city"),
		foodPrices:     make(map[FoodItems]int),
		beveragePrices: make(map[BeverageItems]int),
	}
	for _, food := range foodPrices {
		c.foodPrices[FoodItems(food.first)] = food.second
	}
	for _, beverage := range beveragePrices {
		c.beveragePrices[BeverageItems(beverage.first)] = beverage.second
	}
	return c
}

func (c *City) addStore(store *Store) {
	c.mu.Lock()
	defer c.mu.Unlock()

	store.setFoodRates(c.foodPrices)
	store.setBeverageRates(c.beveragePrices)
	c.stores = append(c.stores, store)
}

func (c *City) purchaseFood(storeId string, foodItem FoodItems, qty int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, store := range c.stores {
		if store.getId() == storeId {
			store.purchaseFood(foodItem, qty)
			break
		}
	}
}

func (c *City) purchaseBeverage(storeId string, beverageItem BeverageItems, qty int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, store := range c.stores {
		if store.getId() == storeId {
			store.purchaseBeverage(beverageItem, qty)
			break
		}
	}
}

func (c *City) getId() string {
	return c.id
}

func (c *City) getStores() []*Store {
	return c.stores
}

type State struct {
	id     string
	cities []*City
	mu     sync.Mutex
}

func NewState() *State {
	return &State{id: getUniqueID("state")}
}

func (s *State) addCity(city *City) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cities = append(s.cities, city)
}

func (s *State) purchaseFood(cityId, storeId string, foodItem FoodItems, qty int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, city := range s.cities {
		if city.getId() == cityId {
			city.purchaseFood(storeId, foodItem, qty)
			break
		}
	}
}

func (s *State) purchaseBeverage(cityId, storeId string, beverageItem BeverageItems, qty int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, city := range s.cities {
		if city.getId() == cityId {
			city.purchaseBeverage(storeId, beverageItem, qty)
			break
		}
	}
}

func (s *State) getCities() []*City {
	return s.cities
}

func (s *State) getId() string {
	return s.id
}

type FoodSystem struct {
	states []*State
	mu     sync.Mutex
}

func NewFoodSystem() *FoodSystem {
	return &FoodSystem{}
}

func (sys *FoodSystem) addState(state *State) {
	sys.mu.Lock()
	defer sys.mu.Unlock()
	sys.states = append(sys.states, state)
}

func (sys *FoodSystem) purchaseFood(stateId, cityId, storeId string, foodItem FoodItems, qty int) {
	sys.mu.Lock()
	defer sys.mu.Unlock()

	for _, state := range sys.states {
		if state.getId() == stateId {
			state.purchaseFood(cityId, storeId, foodItem, qty)
			break
		}
	}
}

func (sys *FoodSystem) purchaseBeverage(stateId, cityId, storeId string, beverageItem BeverageItems, qty int) {
	sys.mu.Lock()
	defer sys.mu.Unlock()

	for _, state := range sys.states {
		if state.getId() == stateId {
			state.purchaseBeverage(cityId, storeId, beverageItem, qty)
			break
		}
	}
}

func (sys *FoodSystem) getStates() []*State {
	return sys.states
}

type pair struct {
	first, second int
}

var uniqueIDCounter = make(map[string]int)
var uniqueIDMutex sync.Mutex

func getUniqueID(prefix string) string {
	uniqueIDMutex.Lock()
	defer uniqueIDMutex.Unlock()

	id := uniqueIDCounter[prefix] + 1
	uniqueIDCounter[prefix] = id
	return fmt.Sprintf("%s%d", prefix, id)
}

func SalesManagement() {
	foodSupply := []pair{{0, 1}, {1, 2}, {2, 3}, {3, 4}}
	beverageSupply := []pair{{0, 1}, {1, 2}, {2, 3}}

	store := NewStore(foodSupply, beverageSupply)
	city := NewCity(foodSupply, beverageSupply)
	city.addStore(store)

	state := NewState()
	state.addCity(city)

	foodSystem := NewFoodSystem()
	foodSystem.addState(state)

	state1 := "state1"
	city1 := "city1"
	store1 := "store1"

	foodSystem.purchaseFood(state1, city1, store1, Burger, 2)
	foodSystem.purchaseBeverage(state1, city1, store1, Coffee, 1)

	for _, state := range foodSystem.getStates() {
		if state.getId() == state1 {
			for _, city := range state.getCities() {
				if city.getId() == city1 {
					for _, store := range city.getStores() {
						if store.getId() == store1 {
							for foodItem, qty := range store.getFoodUnitsSold() {
								fmt.Printf("%d %d\n", foodItem, qty)
							}
							break
						}
					}
					break
				}
			}
			break
		}
	}

	foodSystem.purchaseFood(state1, city1, store1, Burger, 3)
}
