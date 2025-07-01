package utils

import (
	"log/slog"
	"os"
	"reflect"
	"runtime/debug"
	"sync"
	"syscall"

	srv "github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
)

var (
	logger             = slog.Default()
	server *srv.Server = nil
)

func SetServer(s *srv.Server) {
	server = s
}

func Panics[T any](t T, err error) T {
	Panic(err)
	return t
}

func EnumPanic[T any](err error) (t T) {
	Panic(err)
	return t
}

func Panic(err error) {
	if err != nil {
		errorCode := RandString(6)
		logger.With("code", errorCode).Error(err.Error())
		debug.PrintStack()

		if server != nil {
			var wg sync.WaitGroup
			wg.Add(server.PlayerCount())
			go func() {
				iter := reflect.ValueOf(FetchPrivateField[any](server, "p")).MapRange()
				for iter.Next() {
					p := iter.Value().Interface()
					handle := FetchPrivateField[*world.EntityHandle](p, "handle")
					data := FetchPrivateField[world.EntityData](handle, "data")
					s := FetchPrivateField[*session.Session](data.Data, "s")
					if s != nil {
						s.Disconnect("Disconnected [code: " + errorCode + "]")
					}
					wg.Done()
				}
			}()
			wg.Wait()
		}
		os.Exit(int(syscall.SIGTERM)) // closes the server from this
	}
}
