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
	Global struct {
		Misc struct {
			NowPlaying string
		}
		Commands struct {
			Success struct {
				Hub      string
				Ping     string
				PingSelf string
				GameMode string
				GiveRank string
			}
			Error struct {
				Permission    string
				CoolDown      string
				OnlyOneTarget string
				RankHierarchy string
			}
		}

		WorldEdit struct {
			Success struct {
				Cylinder string
				Fill     string
				Mirror   string
				Paste    string
				PosSet   string
				Redo     string
				Replace  string
				Rotate   string
				Sphere   string
				Undo     string
				Up       string
				Wall     string
				Wand     string
			}

			Error struct {
				PosNotSet     string
				BlockNotExist string
				NothingToUndo string
				NothingToRedo string
			}
		}

		Game struct {
			WaitingForPlayers string
			JoinGame          string
			QuitGame          string
			Countdown         string
			ForceStartGame    string
			ChangeTeams       string
			NicknameSaved     string
			ChatColorSaved    string
			CosmeticEquipped  string
			CosmeticOwned     string
			CosmeticUnowned   string

			Error struct {
				NotInAGame         string
				TeamIsFull         string
				GameAlreadyStarted string
				CosmeticNotOwned   string
			}
		}

		Error struct {
			InventoryFull string
		}
	}
	TowerWars struct {
		YouDied      string
		RespawningIn string
		ShieldMana   string
		Error        struct {
			CannotAttackSummoner  string
			CannotAttackAllyTower string
		}
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
