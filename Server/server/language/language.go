package language

import (
	"path"
	"server/server/utils"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"golang.org/x/text/language"
)

var Langfigs = map[string]LangConfig{}

func Translate(pl *player.Player) LangConfig {
	return Langfigs[fileCodes[pl.Locale().String()]]
}

func TranslateWithH(entityHandle *world.EntityHandle, tx *world.Tx) LangConfig {
	if e, ok := entityHandle.Entity(tx); ok {
		return Langfigs[fileCodes[e.(*player.Player).Locale().String()]]
	}
	return LangConfig{}
}

type LangConfig struct {
	Misc struct {
		NowPlaying       string
		SelectedWoodSkin string
		SelectedCape     string
	}

	Commands struct {
		Success struct {
			Hub          string
			Ping         string
			PingSelf     string
			GiveRank     string
			Link         string
			Warp         string
			YouGotWarped string
		}

		Error struct {
			Permission          string
			CoolDown            string
			OnlyOneTarget       string
			RankHierarchy       string
			LinkExpired         string
			NoMorePlayersToWarp string
			NoGameToJoin        string
		}
	}

	Game struct {
		JoinGame    string
		QuitGame    string
		Countdown   string
		ChangeTeams string

		Error struct {
			NotInAGame                             string
			TeamIsFull                             string
			GameAlreadyStarted                     string
			CannotBreakBlocksBecauseGameNotStarted string
		}
	}

	BedWars struct {
		YouWonTitle     string
		YouLostTitle    string
		TeamWonSubTitle string
		YouDiedTitle    string
		YouDiedSubTitle string
		RespawningIn    string
		KilledBy        string
		VoidDeath       string
		BedBreak        string
		GiveIron        string
		GiveGold        string
		GiveDiamond     string
		GiveEmerald     string
		Error           struct {
			CannotBreakBed string
			CannotBreakMap string
		}
	}

	BuildFFA struct {
		JoinMessage string
		YouDied     string
		KilledBy    string
		VoidDeath   string
	}

	Error struct {
		InventoryFull string
	}
}

var fileCodes = make(map[string]string)

func RegisterLanguages(langs map[string][]string) {
	for lang, codes := range langs {
		utils.Panic(RegisterLanguage(lang))
		for _, code := range codes {
			fileCodes[code] = lang
		}
	}
}

func RegisterLanguage(langFile string) error {
	lang, err := utils.ReadConfig[LangConfig](path.Join(".", "translations", langFile))
	if err != nil {
		return err
	}
	Langfigs[langFile] = lang
	return nil
}

type TranslateString string

func (s TranslateString) Resolve(language.Tag) string { return string(s) }
