package dbstore

type Mysql struct {
	FmtDsn string `json:"fmt_dsn"`
}
type Postgres struct {
	FmtDsn string `json:"fmt_dsn"`
}
type DMv8 struct {
	FmtDsn string `json:"fmt_dsn"`
}

type KBv8r6 struct {
	FmtDsn string `json:"fmt_dsn"`
}

type Sqlite struct {
	Dsn string `json:"dsn"`
}
type Config struct {
	Name                string   `json:"name"`
	CreateDbIfNotExists bool     `json:"create_db_if_not_exists"`
	AutoMirgate         bool     `json:"auto_mirgate"`
	MirgateOnly         bool     `json:"mirgate_only"`
	Typ                 string   `json:"typ"`
	Mysql               Mysql    `json:"mysql"`
	Postgres            Postgres `json:"postgres"`
	DMv8                DMv8     `json:"dm_v8"`
	KBv8r6              KBv8r6   `json:"kb_v8r6"`
	Sqlite              Sqlite   `json:"sqlite"`
}
