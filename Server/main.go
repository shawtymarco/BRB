package main

import (
	"log/slog"
	"os"
	"path"
	core "server/server"
	"server/server/command"
	"server/server/database"
	"server/server/game/maps"
	"server/server/language"
	"server/server/user"
	"server/server/utils"

	"github.com/df-mc/dragonfly/server/player/chat"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	chat.Global.Subscribe(chat.StdoutSubscriber{})

	log := slog.Default()

	serverConf := utils.Panics(utils.ReadConfig[core.Server](path.Join(".", "config", "server.json")))
	core.Config = serverConf

	mapfig := utils.Panics(utils.ReadConfig[maps.MapCollection](path.Join(".", "config", "maps.json")))

	maps.MapPath = "./server/data/maps/twmap.zip"
	maps.MapConfig = mapfig

	language.RegisterLanguages(serverConf.Languages)

	if conn := os.Getenv("DATABASE_URI"); conn != "" { //"mongodb://localhost:27017" OR "mongodb://root:password@mongodb:27017"
		core.Database = utils.Panics(database.NewMongoDBDatabase(conn))
	} else {
		core.Database = database.NewLocalDatabase()
	}

	log.Info("Successfully connected to the database!", "type", core.Database.String())

	command.GlobalCommands()

	c := core.DefaultConfig()
	conf := utils.Panics(c.Config(log))
	conf.ShutdownMessage = chat.Translate(language.TranslateString("%disconnect.disconnected"), 1, "")

	conf.ReadOnlyWorld = true
	srv := conf.New()
	utils.SetServer(srv)
	srv.Listen()
	srv.World().StopWeatherCycle()
	srv.World().StopRaining()
	srv.World().StopThundering()
	srv.World().StopTime()
	srv.CloseOnProgramEnd()

	for pl := range srv.Accept() {
		utils.Panics(user.New(pl, false))
	}

	for identifier, err := range core.Database.SaveAll() {
		errorCode := utils.RandString(6)
		log.With("code", errorCode).With("identifier", identifier).Error(err.Error())
	}
}
