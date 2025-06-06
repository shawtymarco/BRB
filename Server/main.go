package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	core "server/server"
	"server/server/command"
	"server/server/database"
	"server/server/game"
	"server/server/games/buildffa"
	"server/server/games/lobby"
	"server/server/language"
	"server/server/user"
	"server/server/utils"
	"server/server/worldmanager"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/npc"

	"github.com/df-mc/dragonfly/server/player/chat"

	_ "server/server/api"
	_ "server/server/games"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	chat.Global.Subscribe(chat.StdoutSubscriber{})

	log := slog.Default()

	serverConf := utils.Panics(utils.ReadConfig[core.Server](path.Join(".", "config", "server.json")))
	core.Config = serverConf

	game.Maps = utils.Panics(utils.ReadConfig[map[string]game.MapData](path.Join(".", "config", "maps.json")))

	language.RegisterLanguages(serverConf.Languages)

	if conn := os.Getenv("DATABASE_URI"); conn != "" { //"mongodb://localhost:27017" OR "mongodb://root:password@mongodb:27017"
		core.Database = utils.Panics(database.NewMongoDBDatabase(conn))
	} else {
		core.Database = database.NewLocalDatabase()
	}

	log.Info("Successfully connected to the database!", "type", core.Database.String())

	command.RegisterCommands()

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
	core.MCServer = srv

	srv.World().Exec(func(tx *world.Tx) {
		bot := npc.Create(npc.Settings{
			Name:     fmt.Sprintf("Mark"),
			Skin:     npc.MustSkin(npc.MustParseTexture(path.Join(".", "config", "skins", "mark.png")), npc.DefaultModel),
			Position: core.Config.Hub.SpawnPoint,
			Scale:    1,
		}, tx, nil)

		utils.Panics(user.New(bot, true))
		lobby.Join(bot)

		core.Bot = bot
	})

	core.WorldManager = utils.Panics(worldmanager.ManagerSettings{
		Folder:   path.Join(".", "maps"),
		Fallback: srv.World(),
		Logger:   log,
	}.NewManager())

	buildffa.NewBuildFFA()

	for pl := range srv.Accept() {
		utils.Panics(user.New(pl, false))
		lobby.Join(pl)
	}

	for identifier, err := range core.Database.SaveAll() {
		errorCode := utils.RandString(6)
		log.With("code", errorCode).With("identifier", identifier).Error(err.Error())
	}
}
