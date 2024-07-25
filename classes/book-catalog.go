package classes

import (
	"container/heap"
	"fmt"
	"strings"
)

type Category int

const (
	FICTION Category = iota
	SCI_FI
	MYSTERY
	FABLE
	MYTHOLOGY
)

type Book struct {
	Name        string
	Author      string
	Publisher   string
	PublishYear int
	Category    Category
	Price       float64
	Count       int
}

type Catalog struct {
	books         []Book
	authorMap     map[string][]*Book
	authorPQMap   map[string]*BookPriorityQueue
	categoryPQMap map[Category]*BookPriorityQueue
}

func NewCatalog() *Catalog {
	return &Catalog{
		authorMap:     make(map[string][]*Book),
		authorPQMap:   make(map[string]*BookPriorityQueue),
		categoryPQMap: make(map[Category]*BookPriorityQueue),
	}
}

func (c *Catalog) AddBookToCatalog(book Book) {
	c.books = append(c.books, book)
	bookPtr := &c.books[len(c.books)-1]

	// pointer to the book is stored in the authorMap
	c.authorMap[book.Author] = append(c.authorMap[book.Author], bookPtr)

	if _, exists := c.authorPQMap[book.Author]; !exists {
		c.authorPQMap[book.Author] = &BookPriorityQueue{}
		heap.Init(c.authorPQMap[book.Author])
	}
	heap.Push(c.authorPQMap[book.Author], bookPtr)

	if _, exists := c.categoryPQMap[book.Category]; !exists {
		c.categoryPQMap[book.Category] = &BookPriorityQueue{}
		heap.Init(c.categoryPQMap[book.Category])
	}
	heap.Push(c.categoryPQMap[book.Category], bookPtr)
}

func (c *Catalog) SearchBookByName(prefix string) []Book {
	var bookList []Book
	for _, book := range c.books {
		if strings.HasPrefix(book.Name, prefix) {
			bookList = append(bookList, book)
		}
	}
	return bookList
}

func (c *Catalog) SearchBookByAuthor(authorName string) []Book {
	var bookList []Book
	for _, book := range c.authorMap[authorName] {
		bookList = append(bookList, *book)
	}
	return bookList
}

func (c *Catalog) GetMostSoldBooksByAuthor(authorName string, limit int) []Book {
	var bookList []Book
	if pq, exists := c.authorPQMap[authorName]; exists {
		for limit > 0 && pq.Len() > 0 {
			bookList = append(bookList, *heap.Pop(pq).(*Book))
			limit--
		}
	}
	return bookList
}

func (c *Catalog) GetMostSoldBooksByCategory(category Category, limit int) []Book {
	var bookList []Book
	if pq, exists := c.categoryPQMap[category]; exists {
		for limit > 0 && pq.Len() > 0 {
			bookList = append(bookList, *heap.Pop(pq).(*Book))
			limit--
		}
	}
	return bookList
}

// Priority queue implementation for books
type BookPriorityQueue []*Book

func (pq BookPriorityQueue) Len() int { return len(pq) }

// see it carefully, we are using > instead of < because we want to pop the book with highest count and what are the parameters for Less function
func (pq BookPriorityQueue) Less(i, j int) bool {
	return pq[i].Count > pq[j].Count
}

func (pq BookPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *BookPriorityQueue) Push(x interface{}) {
	item := x.(*Book)
	*pq = append(*pq, item)
}

func (pq *BookPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func BookCatalogSystem() {
	book := Book{"HP & The PS", "J K Rowling", "Bloomsbury", 1997, FICTION, 200, 80}
	book1 := Book{"HP & The COS", "J K Rowling", "Bloomsbury", 1998, FICTION, 1000, 100}
	book2 := Book{"HP & The POA", "J K Rowling", "Bloomsbury", 1999, FICTION, 2000, 500}
	book3 := Book{"HP & The HBP", "J K Rowling", "Bloomsbury", 2005, FICTION, 3000, 700}
	book4 := Book{"The Immortals of Meluha", "Amish", "Westland", 2010, MYTHOLOGY, 1500, 600}
	book5 := Book{"The Secret of the Nagas", "Amish", "Westland", 2011, MYTHOLOGY, 2500, 400}
	book6 := Book{"The Oath of the Vayuputras", "Amish", "Westland", 2013, MYTHOLOGY, 3500, 200}
	book7 := Book{"Do Androids Dream of Electric Sheep", "Philip K Dick", "DoubleDay", 1968, SCI_FI, 30, 20}

	catalog := NewCatalog()
	catalog.AddBookToCatalog(book)
	catalog.AddBookToCatalog(book1)
	catalog.AddBookToCatalog(book2)
	catalog.AddBookToCatalog(book3)
	catalog.AddBookToCatalog(book4)
	catalog.AddBookToCatalog(book5)
	catalog.AddBookToCatalog(book6)
	catalog.AddBookToCatalog(book7)

	list := catalog.GetMostSoldBooksByAuthor("Amish", 2)
	for _, book := range list {
		fmt.Println(book.Name, book.Count)
	}

	fmt.Println("**************************************************************************************")

	list = catalog.GetMostSoldBooksByCategory(FICTION, 2)
	for _, book := range list {
		fmt.Println(book.Name, book.Count)
	}

	fmt.Println("**************************************************************************************")

	list = catalog.SearchBookByAuthor("Amish")
	for _, book := range list {
		fmt.Println(book.Name, book.Count)
	}

	fmt.Println("**************************************************************************************")

	list = catalog.SearchBookByName("Do")
	for _, book := range list {
		fmt.Println(book.Name, book.Count)
	}
}
