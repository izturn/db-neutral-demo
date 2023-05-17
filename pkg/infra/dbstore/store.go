package dbstore

import (
	"fmt"

	"github.com/izturn/db-neutral-demo/pkg/domain"
	"github.com/izturn/db-neutral-demo/pkg/infra/algoutil"
	"github.com/izturn/db-neutral-demo/pkg/infra/errutil"

	dm "github.com/izturn/gorm2-dm8"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	typMysql    = "mysql"
	typPostgres = "postgres"
	typDMv8     = "dm_v8"

	typKBv8r6  = "kb_v8r6"
	typeSqlite = "sqlite"
)

type Store struct {
	db *gorm.DB
}

func openMySQL(dbName, dsnFmt string, createDBIfNeeds bool, cfg *gorm.Config) (*gorm.DB, error) {
	if createDBIfNeeds {
		defDBDsn := fmt.Sprintf(dsnFmt, "mysql")
		defDB, err := gorm.Open(mysql.Open(defDBDsn), cfg)
		if err != nil {
			return nil, err
		}

		res := defDB.Exec("CREATE DATABASE IF NOT EXISTS " + dbName + ";")
		if res.Error != nil {
			return nil, fmt.Errorf("create db: %s is failed: %w", dbName, err)
		}
		dbInst, _ := defDB.DB()
		_ = dbInst.Close()

	}

	return gorm.Open(mysql.Open(fmt.Sprintf(dsnFmt, dbName)), cfg)
}

// https://dev.to/karanpratapsingh/connecting-to-postgresql-using-gorm-24fj
func openPostgres(dbName, dsnFmt string, createDBIfNeeds bool, cfg *gorm.Config) (*gorm.DB, error) {
	if createDBIfNeeds {
		defDBDsn := fmt.Sprintf(dsnFmt, "postgres")
		defDB, err := gorm.Open(postgres.Open(defDBDsn), cfg)
		if err != nil {
			return nil, err
		}

		count := 0
		defDB.Raw("SELECT COUNT(1) FROM pg_database WHERE datname = ?", dbName).Scan(&count)
		if count == 0 {
			res := defDB.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName))
			if res.Error != nil {
				return nil, fmt.Errorf("create db: %s is failed: %w", dbName, err)
			}
			dbInst, _ := defDB.DB()
			_ = dbInst.Close()
		}

	}
	return gorm.Open(postgres.Open(fmt.Sprintf(dsnFmt, dbName)), cfg)
}

func openDM8(dsnFmt string, cfg *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(dm.Open(dsnFmt), cfg)
}

func MustNew(cfg *Config) *Store {
	gCfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	var db *gorm.DB
	switch cfg.Typ {
	case typMysql:
		db = algoutil.Must1(openMySQL(cfg.Name, cfg.Mysql.FmtDsn, cfg.CreateDbIfNotExists, gCfg))

	case typPostgres:
		db = algoutil.Must1(openPostgres(cfg.Name, cfg.Postgres.FmtDsn, cfg.CreateDbIfNotExists, gCfg))

		// https://eco.dameng.com/document/dm/zh-cn/pm/go-rogramming-guide.html#11.7%20ORM%20%E6%96%B9%E8%A8%80%E5%8C%85
		//https://github.com/housepower/ckman/blob/main/repository/dm8/dm8.go
	case typDMv8:
		db = algoutil.Must1(openDM8(cfg.DMv8.FmtDsn, gCfg))

	case typKBv8r6:

	case typeSqlite:

	default:

	}

	if cfg.AutoMirgate {
		algoutil.Must(db.AutoMigrate(&domain.Book{}))
	}

	return &Store{
		db: db,
	}
}

func (s *Store) GetBookByID(id int) (domain.Book, error) {
	book := domain.Book{
		Id: id,
	}

	res := s.db.First(&book)
	if res.Error == gorm.ErrRecordNotFound {
		return book, errutil.ErrNotFound
	}
	return book, res.Error
}
func (s *Store) UpdateBookWithID(id int, book domain.Book) error {
	book.Id = id
	res := s.db.Save(&book)
	return res.Error
}
func (s *Store) DeleteBookByID(id int) error {
	book := domain.Book{
		Id: id,
	}
	res := s.db.Delete(&book)
	return res.Error
}
func (s *Store) CreateBook(book domain.Book) (int, error) {
	book.Id = 0
	res := s.db.Save(&book)
	return book.Id, res.Error
}
func (s *Store) ListBooks(page, pageSize int) ([]domain.Book, error) {
	if page <= 0 {
		page = 1
	}

	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var books []domain.Book
	res := s.db.Limit(pageSize).Offset(offset).Find(&books)

	return books, res.Error
}
