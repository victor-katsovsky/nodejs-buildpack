package hooks

import (
	"fmt"
	"github.com/cloudfoundry/libbuildpack"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//const sl_cli_template = "./node_modules/.bin/slnodejs run --token  <token> --buildsessionid <id>  --useinitialcolor true  --proxy <proxy> — index.js"
const sl_cli_template = " ./node_modules/.bin/slnodejs run "

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
	sl.Log.Info("Hello from Sealights hook")

	token := os.Getenv("SL_TOKEN")
	bsid := os.Getenv("SL_BUILD_SESSION_ID")
	proxy := os.Getenv("SL_PROXY")
	if token == "" {
		sl.Log.Error("token cannot be empty (env SL_TOKEN)")
		return errors.New("empty token")
	}
	if bsid == "" {
		sl.Log.Error("token cannot be empty (env SL_BUILD_SESSION_ID)")
		return errors.New("empty bsid")
	}

	fmt.Println(proxy)

	bytes, err := ioutil.ReadFile(filepath.Join(stager.BuildDir(), "Procfile"))
	if err != nil {
		sl.Log.Error("failed to read Procfile")
		return err
	}

	splitted := strings.Split(string(bytes), " ")
	web := splitted[0]
	app := splitted[2]

	fmt.Println("web: ", web, "app: ", app)

	final := web + sl_cli_template + fmt.Sprintf(" --token %s --buildsessionid %s ", token, bsid)
	if proxy != "" {
		final += fmt.Sprintf(" --proxy %s ", proxy)
	}

	final += " --useinitialcolor true "
	final += app

	err = ioutil.WriteFile(filepath.Join(stager.BuildDir(), "Procfile"), []byte(final), 0700)
	if err != nil {
		sl.Log.Error("failed to write Procfile, error: %s", err)
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
