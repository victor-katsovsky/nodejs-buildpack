package hooks

import (
	"github.com/cloudfoundry/libbuildpack"
	"io/ioutil"
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

func (sl *SealightsHook) AfterCompile(stager *libbuildpack.Stager) error {
	sl.Log.Info("Hello from Sealights hook")

	//token := os.Getenv("SL_TOKEN")
	//bsid := os.Getenv("SL_BUILD_SESSION_ID")
	//proxy := os.Getenv("SL_PROXY")
	//if token == "" {
	//	return fmt.Errorf("token cannot be empty (env SL_TOKEN)")
	//}
	//if bsid == "" {
	//	return fmt.Errorf("token cannot be empty (env SL_BUILD_SESSION_ID)")
	//}
	//
	//fmt.Println(proxy)

	hn := os.Getenv("NODE_HOME")
	files, err := ioutil.ReadDir(hn)
	if err != nil {
		sl.Log.Info("no folder %s", hn)
		return err
	}
	sl.Log.Info("listing content of: %s", hn)
	for _, file := range files {
		sl.Log.Debug(file.Name())
	}
	return nil

	err = sl.Command.Execute(stager.BuildDir(), os.Stdout, os.Stderr, "npm", "install", "slnodejs")
	if err != nil {
		sl.Log.Error("npm install failed with error: " + err.Error())

		return err
	} else {
		sl.Log.Error("npm install finished sucessfully")
	}
	return nil
}
