package webserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/izturn/db-neutral-demo/pkg/infra/errutil"
	"github.com/izturn/db-neutral-demo/pkg/infra/mux"
	"github.com/izturn/db-neutral-demo/pkg/usecase"
)

type Logger interface {
	Log(args ...interface{})
}

type webHandler interface {
	CreateBook(book usecase.Book) error
	UpdateBook(id int, book usecase.Book) error
	DeleteBook(id int) error
	GetBook(id int) (usecase.Book, error)
	ListBooks(page, pageSize int) ([]usecase.Book, error)
}
type webServer struct {
	*mux.Router

	hldr   webHandler
	addr   string
	logger Logger
}

func New(addr string, hldr webHandler, logger Logger) *webServer {
	router := mux.NewRouter()
	b := &webServer{
		addr:   addr,
		Router: router,
		logger: logger,
		hldr:   hldr,
	}

	router.HandleFunc("GET", "/books/:id", b.GetBook)
	router.HandleFunc("PUT", "/books/:id", b.UpdateBook)
	router.HandleFunc("DELETE", "/books/:id", b.DeleteBook)
	router.HandleFunc("POST", "/books/", b.CreateBook)
	router.HandleFunc("GET", "/books/", b.ListBooks)

	return b
}

func (b *webServer) GetBook(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.IntParam(r.Context(), "id")
	book, err := b.hldr.GetBook(id)
	if err != nil {
		b.logger.Log(fmt.Sprintf("get book: %d is failed: %v", id, err))
		if err != errutil.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

		return
	}
	json.NewEncoder(w).Encode(&book)
}

func (b *webServer) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.IntParam(r.Context(), "id")
	err := b.hldr.DeleteBook(id)
	if err != nil {
		b.logger.Log(fmt.Sprintf("delete book: %d is failed: %v", id, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("delete is ok"))
}

func (b *webServer) CreateBook(w http.ResponseWriter, r *http.Request) {
	buf, _ := io.ReadAll(r.Body)

	var book usecase.Book
	err := json.Unmarshal(buf, &book)
	if err != nil {
		b.logger.Log("unmarshal is failed with error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = b.hldr.CreateBook(book)
	if err != nil {
		b.logger.Log(fmt.Sprintf("create book: %s is failed: %v", book.Title, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("create is ok"))

}
func (b *webServer) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.IntParam(r.Context(), "id")

	buf, _ := io.ReadAll(r.Body)

	var book usecase.Book
	err := json.Unmarshal(buf, &book)
	if err != nil {
		b.logger.Log("unmarshal is failed with error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = b.hldr.UpdateBook(id, book)
	if err != nil {
		b.logger.Log(fmt.Sprintf("update book: %d is failed: %v", id, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("update is ok"))
}

func (b *webServer) ListBooks(w http.ResponseWriter, r *http.Request) {
	page, _ := mux.IntQuery(r.Context(), "page")
	pageSize, _ := mux.IntQuery(r.Context(), "page_size")

	books, err := b.hldr.ListBooks(page, pageSize)
	if err != nil {
		b.logger.Log(fmt.Sprintf("list books is failed: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(books) == 0 {
		books = []usecase.Book{}
	}

	json.NewEncoder(w).Encode(&books)
}

func (ws *webServer) Run() {
	ws.logger.Log("the server is running at:", ws.addr)
	http.ListenAndServe(ws.addr, ws)
}
