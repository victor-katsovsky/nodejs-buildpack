package hooks

import (
	"fmt"
	"github.com/cloudfoundry/libbuildpack"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type SealightsHook struct {
	libbuildpack.DefaultHook
	Log             *libbuildpack.Logger
	Command         *libbuildpack.Command
	EmptyTokenError string `default:"token cannot be empty (env SL_TOKEN)"`
	EmptyBuildError string `default:"build session id cannot be empty (env SL_BUILD_SESSION_ID)"`
}

func init() {
	logger := libbuildpack.NewLogger(os.Stdout)
	command := &libbuildpack.Command{}
	libbuildpack.AddHook(&SealightsHook{
		Log:     logger,
		Command: command,
	})
}

func (sl *SealightsHook) AfterCompile(stager *libbuildpack.Stager) error {
	sl.Log.Info("Inside Sealights hook")

	token := os.Getenv("SL_TOKEN")
	bsid := os.Getenv("SL_BUILD_SESSION_ID")
	proxy := os.Getenv("SL_PROXY")
	if token == "" {
		sl.Log.Error(sl.EmptyTokenError)
		return fmt.Errorf(sl.EmptyTokenError)
	}
	if bsid == "" {
		sl.Log.Error(sl.EmptyBuildError)
		return fmt.Errorf(sl.EmptyBuildError)
	}

	err := sl.installAgent(stager)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadFile(filepath.Join(stager.BuildDir(), "Procfile"))
	if err != nil {
		sl.Log.Error("failed to read Procfile")
		return err
	}

	split := strings.Split(string(bytes), " ")
	app := split[2]

	newCmd := sl.createAppStartCommandLine(app, token, bsid, proxy)

	sl.Log.Debug("new command line: %s", newCmd)

	err = ioutil.WriteFile(filepath.Join(stager.BuildDir(), "Procfile"), []byte(newCmd), 0755)
	if err != nil {
		sl.Log.Error("failed to update Procfile, error: %s", err)
		return err
	}

	return nil
}

func (sl *SealightsHook) installAgent(stager *libbuildpack.Stager) error {
	err := sl.Command.Execute(stager.BuildDir(), os.Stdout, os.Stderr, "npm", "install", "slnodejs")
	if err != nil {
		sl.Log.Error("npm install slnodejs failed with error: " + err.Error())
		return err
	}
	sl.Log.Info("npm install slnodejs finished successfully")
	return nil
}

func (sl *SealightsHook) createAppStartCommandLine(app, token, bsid, proxy string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("web: node ./node_modules/.bin/slnodejs run  --useinitialcolor true --token %s --buildsessionid %s ", token, bsid))

	if proxy != "" {
		sb.WriteString(fmt.Sprintf(" --proxy %s ", proxy))
	}

	sb.WriteString(fmt.Sprintf(" %s", app))
	return sb.String()
}
