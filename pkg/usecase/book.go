package usecase

import domain "github.com/izturn/db-neutral-demo/pkg/domain"

type Book struct {
	Id          int
	Title       string
	Author      string
	Description string
}

type BookInteractor struct {
	repo domain.BookRepository
}

func NewBookInteractor(repo domain.BookRepository) *BookInteractor {
	b := &BookInteractor{
		repo: repo,
	}
	return b
}

func (b *BookInteractor) CreateBook(book Book) error {
	bk := domain.Book{
		Author: book.Author,
		Title:  book.Title,
		Desc:   book.Description,
	}
	_, err := b.repo.Create(bk)
	return err

}
func (b *BookInteractor) UpdateBook(id int, book Book) error {
	bk := domain.Book{
		Id:     0,
		Author: book.Author,
		Title:  book.Title,
		Desc:   book.Description,
	}
	return b.repo.Update(id, bk)
}
func (b *BookInteractor) DeleteBook(id int) error {
	return b.repo.Delete(id)
}
func (b *BookInteractor) GetBook(id int) (Book, error) {
	bk, err := b.repo.Get(id)
	if err != nil {
		return Book{}, err
	}
	return Book{
		Id:          bk.Id,
		Title:       bk.Title,
		Author:      bk.Author,
		Description: bk.Desc,
	}, nil
}
func (b *BookInteractor) ListBooks(page, pageSize int) ([]Book, error) {
	bb, err := b.repo.List(page, pageSize)
	if err != nil {
		return nil, err
	}

	var books []Book
	for _, v := range bb {
		books = append(books, Book{
			Id:          v.Id,
			Title:       v.Title,
			Author:      v.Author,
			Description: v.Desc,
		})
	}

	return books, nil

}
