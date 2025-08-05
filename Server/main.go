package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"reflect"
	core "server/server"
	"server/server/capes"
	"server/server/command"
	"server/server/database"
	"server/server/game"
	"server/server/games/bedwars"
	"server/server/games/buildffa"
	"server/server/games/lobby"
	"server/server/language"
	"server/server/listener"
	"server/server/living/npc"
	"server/server/user"
	"server/server/utils"
	"server/server/worldmanager"
	"slices"
	"strings"
	"time"
	"unsafe"

	"github.com/bedrock-gophers/intercept/intercept"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/google/uuid"

	"github.com/df-mc/dragonfly/server/player"

	"github.com/df-mc/dragonfly/server/world"

	"github.com/df-mc/dragonfly/server/player/chat"

	_ "net/http/pprof"
	_ "server/server/api"
	_ "server/server/item_behavior"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	chat.Global.Subscribe(chat.StdoutSubscriber{})

	log := slog.Default()

	go func() {
		log.Info("Starting pprof on :6060")
		utils.Panic(http.ListenAndServe(":6060", nil))
	}()

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

	//multiversion.ListenerFunc(&conf, c.Network.Address, []minecraft.Protocol{
	//	v486.New(true),
	//})

	intercept.Hook(listener.PacketHandler{})
	srv := conf.New()
	utils.SetServer(srv)
	srv.Listen()
	srv.World().StopWeatherCycle()
	srv.World().StopRaining()
	srv.World().StopThundering()
	srv.World().StopTime()
	srv.World().Handle(listener.WorldHandler{})
	srv.CloseOnProgramEnd()
	core.MCServer = srv

	//srv.World().Exec(func(tx *world.Tx) {
	//	initBots(tx)
	//})

	worldsRoot := path.Join(".", "server", "worlds")
	for _, entry := range utils.Panics(os.ReadDir(worldsRoot)) {
		if entry.IsDir() && entry.Name() != "hub" {
			p := path.Join(worldsRoot, entry.Name())
			_ = os.RemoveAll(p)
		}
	}

	core.WorldManager = utils.Panics(worldmanager.ManagerSettings{
		Folder:   path.Join(".", "maps"),
		Fallback: srv.World(),
		Logger:   log,
	}.NewManager())

	buildffa.NewBuildFFA()
	bw := bedwars.NewBedWars(game.TypeBedWars, 1, 2, false)
	bw.UsersToJoin = []string{"944163286591602688", "436765918169792524"}

	srv.World().Exec(func(tx *world.Tx) {
		txtPos := mgl64.Vec3{-36.5, 99.0, -143.5}
		tx.Block(cube.PosFromVec3(txtPos))
		tx.AddEntity(entity.NewText(text.Colourf("<green>Welcome to BRBW!</green>\n<grey>The #1 Ranked Bedwars server on Bedrock</grey>\n§0\n<white>Join our Discord!.</white>\n<aqua>discord.gg/brbw</aqua>"), txtPos))
	})

	for pl := range srv.Accept() {
		intercept.Intercept(pl)
		u := user.GetUser(pl)

		allData := utils.Panics(core.Database.FindAllPlayers())
		for _, d := range allData {
			if d.UUID == u.Data.UUID {
				continue
			}

			if d.DeviceID == u.Data.DeviceID || time.Since(d.IPStoredSince) <= 24*time.Hour && d.HashedIP == u.Data.HashedIP {
				u.Data.AddAlt(d)
				utils.Panic(core.Database.SavePlayer(d))
			}
		}

		lobby.Join(pl)
		joinRankedBedWars(pl)
	}

	for identifier, err := range core.Database.SaveAll() {
		errorCode := utils.RandString(6)
		log.With("code", errorCode).With("identifier", identifier).Error(err.Error())
	}
}

func joinRankedBedWars(pl *player.Player) bool {
	u := user.GetUser(pl)
	for _, g := range bedwars.Games {
		if slices.Contains(g.UsersToJoin, u.Data.UserId) {
			pl.Handler().HandleQuit(pl)
			bedwars.Join(pl, pl.Tx(), g.TeamSize, g.TeamCount, g.Type(), false, g)
			return true
		}
	}

	return false
}

func registerCapes() {
	database.RegisterCape(capes.CreeperCape{})
}

func initBots(tx *world.Tx) {
	core.BotMark = createBot("Mark", tx)
	//core.BotSam = createBot("Sam", tx)
	//core.BotSteven = createBot("Steven", tx)
	lobby.Join(core.BotMark)
	//lobby.Join(core.BotSam)
	//lobby.Join(core.BotSteven)
}

func createBot(name string, tx *world.Tx) *player.Player {
	var id uuid.UUID
	bpd, err := core.Database.FindPlayerByName(name, &database.PlayerNameSearchOpts{CaseInsensitive: false, PartialMatch: false})
	if err != nil {
		id, _ = uuid.NewUUID()
	} else {
		id = bpd.UUID
	}
	bot := npc.Create(npc.Settings{
		UUID:     id,
		Name:     name,
		Skin:     npc.MustSkin(npc.MustParseTexture(path.Join(".", "config", "skins", fmt.Sprintf("%v.png", strings.ToLower(name)))), npc.DefaultModel),
		Position: core.Config.Hub.SpawnPoint,
		Scale:    1,
	}, tx, nil)

	utils.Panic(setMapPrivateFieldKey(
		core.MCServer,
		"p",
		bot.UUID(),
		bot.H(),
		bot.XUID(),
		bot.Name(),
	))

	return bot
}

func setMapPrivateFieldKey(structPtr any, fieldName string, key, handle any, xuid, name string) error {
	structValue := reflect.ValueOf(structPtr).Elem()
	field := structValue.FieldByName(fieldName)

	if !field.IsValid() {
		return fmt.Errorf("field %s not found", fieldName)
	}
	if field.Kind() != reflect.Map {
		return fmt.Errorf("field %s is not a map", fieldName)
	}

	field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	mapValueType := field.Type().Elem()

	entryPtr := reflect.New(mapValueType.Elem())
	entry := entryPtr.Elem()

	setUnexportedField := func(structVal reflect.Value, fieldName string, value any) error {
		f := structVal.FieldByName(fieldName)
		if !f.IsValid() {
			return fmt.Errorf("field %s not found", fieldName)
		}
		ptr := unsafe.Pointer(f.UnsafeAddr())
		reflect.NewAt(f.Type(), ptr).Elem().Set(reflect.ValueOf(value))
		return nil
	}

	if err := setUnexportedField(entry, "handle", handle); err != nil {
		return err
	}
	if err := setUnexportedField(entry, "xuid", xuid); err != nil {
		return err
	}
	if err := setUnexportedField(entry, "name", name); err != nil {
		return err
	}

	field.SetMapIndex(reflect.ValueOf(key), entryPtr)
	return nil
}
