package domain

type Book struct {
	Id     int    `json:"id" gorm:"primaryKey"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Desc   string `json:"desc"`
}

type BookRepository interface {
	Create(book Book) (int, error)
	Update(id int, book Book) error
	Delete(id int) error
	Get(id int) (Book, error)

	List(page, pageSize int) ([]Book, error)
}
