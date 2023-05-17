package config

import (
	"encoding/json"
	"os"

	"github.com/izturn/db-neutral-demo/pkg/infra/algoutil"
	"github.com/izturn/db-neutral-demo/pkg/infra/dbstore"
)

type Config struct {
	Comments string         `json:"comments"`
	Server   Server         `json:"server"`
	DB       dbstore.Config `json:"db"`
}
type Server struct {
	Addr string `json:"addr"`
}

func MustLoad(path string) *Config {
	if path == "" {
		panic("the config path is empty")
	}

	var v Config
	raw := algoutil.Must1(os.ReadFile(path))
	algoutil.Must(json.Unmarshal(raw, &v))
	return &v
}
