package classes

import (
	"fmt"
)

type RideStatus int

const (
	IDLE RideStatus = iota
	CREATED
	WITHDRAWN
	COMPLETED
)

type Ride struct {
	id         int
	origin     int
	dest       int
	seats      int
	rideStatus RideStatus
}

const AMT_PER_KM = 20.0

func (r *Ride) CalculateFare(isPriorityRider bool) float64 {

	dist := float64(r.dest - r.origin)
	baseFare := dist * AMT_PER_KM
	if r.seats < 2 {
		if isPriorityRider {
			return baseFare * 0.75
		}
		return baseFare
	}

	if isPriorityRider {
		return baseFare * float64(r.seats) * 0.5
	}
	return baseFare * float64(r.seats) * 0.75
}

func (r *Ride) SetDest(dest int) {
	r.dest = dest
}

func (r *Ride) GetId() int {
	return r.id
}

func (r *Ride) SetId(id int) {
	r.id = id
}

func (r *Ride) SetOrigin(origin int) {
	r.origin = origin
}

func (r *Ride) GetRideStatus() RideStatus {
	return r.rideStatus
}

func (r *Ride) SetRideStatus(rideStatus RideStatus) {
	r.rideStatus = rideStatus
}

func (r *Ride) SetSeats(seats int) {
	r.seats = seats
}

type Person struct {
	name string
}

type Driver struct {
	Person
}

func NewDriver(name string) Driver {
	return Driver{Person{name: name}}
}

type Rider struct {
	Person
	id             int
	completedRides []Ride
	currentRide    Ride
}

// this is obsolete for this implementation
func NewRider(id int, name string) Rider {
	return Rider{id: id, Person: Person{name: name}}
}

func (r *Rider) CreateRide(id, origin, dest, seats int) {
	if origin >= dest {
		fmt.Println("Wrong values of Origin and Destination provided. Can't create ride")
		return
	}

	r.currentRide.SetId(id)
	r.currentRide.SetOrigin(origin)
	r.currentRide.SetDest(dest)
	r.currentRide.SetSeats(seats)
	r.currentRide.SetRideStatus(CREATED)
}

func (r *Rider) UpdateRide(id, origin, dest, seats int) {
	if r.currentRide.GetRideStatus() == WITHDRAWN {
		fmt.Println("Can't update ride. Ride was withdrawn")
		return
	}
	if r.currentRide.GetRideStatus() == COMPLETED {
		fmt.Println("Can't update ride. Ride already complete")
		return
	}

	r.CreateRide(id, origin, dest, seats)
}

// means we can only withdraw the current ride
func (r *Rider) WithdrawRide(id int) {
	if r.currentRide.GetId() != id {
		fmt.Println("Wrong ride Id as input. Can't withdraw current ride")
		return
	}
	if r.currentRide.GetRideStatus() != CREATED {
		fmt.Println("Ride wasn't in progress. Can't withdraw ride")
		return
	}

	r.currentRide.SetRideStatus(WITHDRAWN)
}

func (r *Rider) GetId() int {
	return r.id
}

// when we close the ride we calculate the fare and add the ride to the completed rides
func (r *Rider) CloseRide() float64 {
	if r.currentRide.GetRideStatus() != CREATED {
		fmt.Println("Ride wasn't in progress. Can't close ride")
		return 0
	}

	r.currentRide.SetRideStatus(COMPLETED)
	r.completedRides = append(r.completedRides, r.currentRide)
	return r.currentRide.CalculateFare(len(r.completedRides) >= 10)
}

type RideSystem struct {
	drivers int
	riders  []Rider
}

func NewRideSystem(drivers int, riders []Rider) *RideSystem {
	if drivers < 1 || len(riders) < 1 {
		fmt.Println("Not enough drivers or riders")
	}

	return &RideSystem{drivers: drivers, riders: riders}
}

func (s *RideSystem) CreateRide(riderId, rideId, origin, dest, seats int) {
	if s.drivers == 0 {
		fmt.Println("No drivers around. Can't create ride")
		return
	}

	for i := range s.riders {
		if s.riders[i].GetId() == riderId {
			s.riders[i].CreateRide(rideId, origin, dest, seats)
			s.drivers--
			break
		}
	}
}

func (s *RideSystem) UpdateRide(riderId, rideId, origin, dest, seats int) {
	for i := range s.riders {
		if s.riders[i].GetId() == riderId {
			s.riders[i].UpdateRide(rideId, origin, dest, seats)
			break
		}
	}
}

func (s *RideSystem) WithdrawRide(riderId, rideId int) {
	for i := range s.riders {
		if s.riders[i].GetId() == riderId {
			s.riders[i].WithdrawRide(rideId)
			s.drivers++
			break
		}
	}
}

func (s *RideSystem) CloseRide(riderId int) float64 {
	for i := range s.riders {
		if s.riders[i].GetId() == riderId {
			s.drivers++
			return s.riders[i].CloseRide()
		}
	}
	return 0
}

func RideSystemClass() {

	aman := NewRider(1708, "Aman")
	aakanksha := NewRider(817, "Aakanksha")
	couples := []Rider{aman, aakanksha}

	rideSystem := NewRideSystem(3, couples)

	fmt.Println("*****************************************************************")

	rideSystem.CreateRide(1708, 1, 50, 60, 1)
	fmt.Println(rideSystem.CloseRide(1708))

	fmt.Println("*****************************************************************")

	rideSystem.CreateRide(817, 2, 50, 60, 1)

	fmt.Println(rideSystem.CloseRide(817))
}
