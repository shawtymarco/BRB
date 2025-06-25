package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	core "server/server"
	"server/server/capes"
	"server/server/command"
	"server/server/database"
	"server/server/game"
	"server/server/games/bedwars"
	"server/server/games/buildffa"
	"server/server/games/lobby"
	"server/server/language"
	"server/server/user"
	"server/server/utils"
	"server/server/worldmanager"
	"slices"
	"strings"

	"github.com/df-mc/dragonfly/server/player"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/npc"

	"github.com/df-mc/dragonfly/server/player/chat"

	_ "server/server/api"
	_ "server/server/item_behavior"
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
	registerCapes()

	c := core.DefaultConfig()
	conf := utils.Panics(c.Config(log))
	conf.Entities = conf.Entities.Config().New([]world.EntityType{&bedwars.GeneratorBlockType{}})
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
		initBots(tx)
	})

	core.WorldManager = utils.Panics(worldmanager.ManagerSettings{
		Folder:   path.Join(".", "maps"),
		Fallback: srv.World(),
		Logger:   log,
	}.NewManager())

	buildffa.NewBuildFFA()
	bedwars.NewBedWars(game.TypeBedWars, 2, 2, false) // TODO: REMOVE DEBUG

	for pl := range srv.Accept() {
		utils.Panics(user.New(pl, false))
		lobby.Join(pl)
		joinRankedBedWars(pl)
	}

	for identifier, err := range core.Database.SaveAll() {
		errorCode := utils.RandString(6)
		log.With("code", errorCode).With("identifier", identifier).Error(err.Error())
	}
}

func registerCapes() {
	database.RegisterCape(capes.CreeperCape{})
}

func initBots(tx *world.Tx) {
	core.BotMark = createBot("Mark", tx)
	utils.Panics(user.New(core.BotMark, true))
	lobby.Join(core.BotMark)

	core.BotSam = createBot("Sam", tx)
	utils.Panics(user.New(core.BotSam, true))
	lobby.Join(core.BotSam)

	core.BotSteven = createBot("Steven", tx)
	utils.Panics(user.New(core.BotSteven, true))
	lobby.Join(core.BotSteven)
}

func createBot(name string, tx *world.Tx) *player.Player {
	return npc.Create(npc.Settings{
		Name:     name,
		Skin:     npc.MustSkin(npc.MustParseTexture(path.Join(".", "config", "skins", fmt.Sprintf("%v.png", strings.ToLower(name)))), npc.DefaultModel),
		Position: core.Config.Hub.SpawnPoint,
		Scale:    1,
	}, tx, nil)
}

func joinRankedBedWars(pl *player.Player) {
	u := user.LookupPlayer(pl)
	for _, g := range bedwars.Games {
		if slices.Contains(g.UsersToJoin, u.Data.UserId) {
			bedwars.Join(pl, pl.Tx(), g.TeamSize, g.TeamCount, g.Type(), false, g)
			break
		}
	}
}
