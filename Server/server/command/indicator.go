package command

import (
	"server/server"
	"server/server/database"
	"server/server/language"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type IndicatorCommand struct {
	Player string `cmd:"player"`
}

func (IndicatorCommand) Allow(src cmd.Source) bool {
	return Indicator.Test(src)
}

func (IndicatorCommand) PermissionMessage(src cmd.Source) string {
	return Indicator.PermissionMessage(src)
}

func (i IndicatorCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	pl, ok := src.(*player.Player)
	if !ok {
		o.Error(text.Colourf("<red>You must run this command in-game.</red>"))
		return
	}

	playerData, err := server.Database.FindPlayerByName(i.Player, &database.PlayerNameSearchOpts{CaseInsensitive: true, PartialMatch: true})
	if err != nil {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.PlayerNotExist))
		return
	}

	deviceName := getDeviceName(playerData.DeviceOS)

	pl.Message(text.Colourf("<green>%s's device: <yellow>%s</yellow></green>", playerData.Username, deviceName))
}

func getDeviceName(deviceOS protocol.DeviceOS) string {
	switch deviceOS {
	case protocol.DeviceAndroid:
		return "Android"
	case protocol.DeviceIOS:
		return "iOS"
	case protocol.DeviceOSX:
		return "macOS"
	case protocol.DeviceFireOS:
		return "Fire OS"
	case protocol.DeviceGearVR:
		return "Gear VR"
	case protocol.DeviceWin32:
		return "Windows"
	case protocol.DeviceDedicated:
		return "Dedicated"
	case protocol.DeviceTVOS:
		return "tvOS"
	case protocol.DeviceLinux:
		return "Linux"
	default:
		return "Unknown"
	}
}
