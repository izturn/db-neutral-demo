package main

import (
	"flag"
	"os"

	"github.com/izturn/db-neutral-demo/pkg/infra/dbstore"
	"github.com/izturn/db-neutral-demo/pkg/infra/logger"
	"github.com/izturn/db-neutral-demo/pkg/interface/config"
	"github.com/izturn/db-neutral-demo/pkg/interface/dbadapter"
	"github.com/izturn/db-neutral-demo/pkg/interface/webserver"
	"github.com/izturn/db-neutral-demo/pkg/usecase"
)

var cfgPath = flag.String("c", "./config.json", "the path of config file")

func main() {
	flag.Parse()

	cfg := config.MustLoad(*cfgPath)
	s := dbstore.MustNew(&cfg.DB)
	l := logger.New()
	if cfg.DB.MirgateOnly {
		l.Log("mirgate is finished, see u again!")
		os.Exit(0)
	}
	us := usecase.NewBookInteractor(dbadapter.New(s))
	webserver.New(cfg.Server.Addr, us, l).Run()
}
