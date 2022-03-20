package hooks

import (
	"fmt"
	"github.com/cloudfoundry/libbuildpack"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//const sl_cli_template = "./node_modules/.bin/slnodejs run --token  <token> --buildsessionid <id>  --useinitialcolor true  --proxy <proxy> — index.js"
const sl_cli_template = " node ./node_modules/.bin/slnodejs run "

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

func (sl *SealightsHook) AfterCompile(stager *libbuildpack.Stager) error {
	sl.Log.Info("Inside Sealights hook")

	token := os.Getenv("SL_TOKEN")
	bsid := os.Getenv("SL_BUILD_SESSION_ID")
	proxy := os.Getenv("SL_PROXY")
	if token == "" {
		sl.Log.Error("token cannot be empty (env SL_TOKEN)")
		return fmt.Errorf("token cannot be empty (env SL_TOKEN)")
	}
	if bsid == "" {
		sl.Log.Error("build session id cannot be empty (env SL_BUILD_SESSION_ID)")
		return fmt.Errorf("build session id cannot be empty (env SL_BUILD_SESSION_ID)")
	}

	bytes, err := ioutil.ReadFile(filepath.Join(stager.BuildDir(), "Procfile"))
	if err != nil {
		sl.Log.Error("failed to read Procfile")
		return err
	}

	splitted := strings.Split(string(bytes), " ")
	app := splitted[2]

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("web: node ./node_modules/.bin/slnodejs run  --useinitialcolor true --token %s --buildsessionid %s ", token, bsid))

	if proxy != "" {
		sb.WriteString(fmt.Sprintf(" --proxy %s ", proxy))
	}

	sb.WriteString(fmt.Sprintf(" %s", app))

	sl.Log.Debug("new command line: %s", sb.String())

	err = ioutil.WriteFile(filepath.Join(stager.BuildDir(), "Procfile"), []byte(sb.String()), 0755)
	if err != nil {
		sl.Log.Error("failed to update Procfile, error: %s", err)
		return err
	}

	err = sl.Command.Execute(stager.BuildDir(), os.Stdout, os.Stderr, "npm", "install", "slnodejs")
	if err != nil {
		sl.Log.Error("npm install slnodejs failed with error: " + err.Error())

		return err
	} else {
		sl.Log.Info("npm install slnodejs finished successfully")
	}
	return nil
}
