package hooks

import (
	"github.com/cloudfoundry/libbuildpack"
	"os"
)

type SealightsHook struct {
	libbuildpack.DefaultHook
	Log     *libbuildpack.Logger
	Command *libbuildpack.Command
}

func init() {
	logger := libbuildpack.NewLogger(os.Stdout)
	command := &libbuildpack.Command{}
	libbuildpack.AddHook(&SealightsHook{
		Log:     logger,
		Command: command,
	})
}

func (ls *SealightsHook) AfterCompile(compiler *libbuildpack.Stager) error {
	ls.Log.Info("Hello from Sealights hook")
	return nil
}
