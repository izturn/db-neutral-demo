package dbadapter

import (
	"github.com/izturn/db-neutral-demo/pkg/domain"
	"github.com/izturn/db-neutral-demo/pkg/infra/dbstore"
)

type DBStoreAdapter struct {
	store *dbstore.Store
}

func New(s *dbstore.Store) *DBStoreAdapter {
	return &DBStoreAdapter{
		store: s,
	}
}
func (r *DBStoreAdapter) Create(book domain.Book) (int, error) {
	return r.store.CreateBook(book)

}
func (r *DBStoreAdapter) Update(id int, book domain.Book) error {
	return r.store.UpdateBookWithID(id, book)
}
func (r *DBStoreAdapter) Delete(id int) error {
	return r.store.DeleteBookByID(id)
}
func (r *DBStoreAdapter) Get(id int) (domain.Book, error) {
	return r.store.GetBookByID(id)
}
func (r *DBStoreAdapter) List(page, pageSize int) ([]domain.Book, error) {
	return r.store.ListBooks(page, pageSize)
}
