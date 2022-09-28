package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	handler "github.com/kuzkuss/VK_DB_Project/api"
	forums "github.com/kuzkuss/VK_DB_Project/internal/app/models/forumsRepository"
	"github.com/kuzkuss/VK_DB_Project/internal/server"
	storeDB "github.com/kuzkuss/VK_DB_Project/internal/store"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/configs.toml", "path to config file")
}

func main() {
	// var errtest error
	// fmt.Println(errtest.Error())

	// return

	flag.Parse()

	db, err := server.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	dbForums := storeDB.NewDataBaseForums(db)
	dbUsers := storeDB.NewDataBaseUsers(db)
	fr := forums.NewForumsRep(dbForums, dbUsers)
	h := handler.NewForumRouter(fr)
	
	config := server.NewConfig()
	_, err = toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(h, config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}