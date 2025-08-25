package proxy

import (
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"os"
)

var packsContentKeys = map[string]string{
	"637d8f5e-bc6d-4b2e-b396-f571a3c51542": "",
	"a6e484e6-0388-42ca-8b0b-711f1c4d2899": "",
}

func ParsePacks() []*resource.Pack {
	dir := "/shared/resources"
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var packs []*resource.Pack
	for _, entry := range entries {
		pack, err := resource.ReadPath(dir + "/" + entry.Name())
		if err != nil {
			panic(err)
		}

		key, ok := packsContentKeys[pack.UUID().String()]
		if ok {
			pack.WithContentKey(key)
		}
		packs = append(packs, pack)
	}
	return packs
}
