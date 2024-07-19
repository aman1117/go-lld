package classes

import (
	"fmt"
)

type Meeting struct {
	start, end int
	room       *Room
}

func (m *Meeting) GetStart() int {
	return m.start
}
func (m *Meeting) GetEnd() int {
	return m.end
}

type Room struct {
	name     string
	calendar map[int][]*Meeting
}

func NewRoom(name string) *Room {
	return &Room{name: name, calendar: make(map[int][]*Meeting)}
}
func (r *Room) Book(day, start, end int) bool {
	for _, meeting := range r.calendar[day] {
		if start < meeting.GetEnd() && end > meeting.GetStart() {
			return false
		}
	}
	r.calendar[day] = append(r.calendar[day], &Meeting{start: start, end: end, room: r})
	return true
}

type Scheduler struct {
	rooms []*Room
}

func NewScheduler(rooms []*Room) *Scheduler {
	return &Scheduler{rooms: rooms}
}
func (r *Room) GetName() string {
	return r.name
}
func (s *Scheduler) Book(day, start, end int) string {
	for _, room := range s.rooms {
		if room.Book(day, start, end) {
			return room.GetName()
		}
	}
	return "No room available"
}
func MeetingScheduler() {
	room1 := NewRoom("Atlas")
	room2 := NewRoom("Nexus")
	room3 := NewRoom("HolyCow")

	rooms := []*Room{room1, room2, room3}

	scheduler := NewScheduler(rooms)

	fmt.Println(scheduler.Book(15, 2, 5)) // Atlas
	fmt.Println(scheduler.Book(15, 5, 8)) // Atlas
	fmt.Println(scheduler.Book(15, 4, 8)) // Nexus
	fmt.Println(scheduler.Book(15, 3, 6)) // HolyCow
	fmt.Println(scheduler.Book(15, 7, 8)) // HolyCow
	fmt.Println(scheduler.Book(16, 6, 9)) // Atlas
}
