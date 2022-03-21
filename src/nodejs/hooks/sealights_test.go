package hooks_test

import (
	"bytes"
	"github.com/cloudfoundry/libbuildpack"
	"github.com/cloudfoundry/nodejs-buildpack/src/nodejs/hooks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var _ = Describe("Sealights hook", func() {
	var (
		err          error
		buildDir     string
		logger       *libbuildpack.Logger
		buffer       *bytes.Buffer
		stager       *libbuildpack.Stager
		sealights    *hooks.SealightsHook
		token        string
		build        string
		proxy        string
		procfile     string
		testProcfile = "web: node index.js --build 192 --name Good"
		expected     = strings.ReplaceAll("web: node ./node_modules/.bin/slnodejs run  --useinitialcolor true --token application_token --buildsessionid build312 --proxy http://localhost:1886 index.js --build 192 --name Good", " ", "")
	)

	BeforeEach(func() {
		buildDir, err = ioutil.TempDir("", "nodejs-buildpack.build.")
		Expect(err).To(BeNil())

		buffer = new(bytes.Buffer)
		logger = libbuildpack.NewLogger(buffer)
		args := []string{buildDir, ""}
		stager = libbuildpack.NewStager(args, logger, &libbuildpack.Manifest{})

		token = os.Getenv("SL_TOKEN")
		build = os.Getenv("SL_BUILD_SESSION_ID")
		proxy = os.Getenv("SL_PROXY")
		err = ioutil.WriteFile(filepath.Join(stager.BuildDir(), "Procfile"), []byte(testProcfile), 0755)
		Expect(err).To(BeNil())

		sealights = &hooks.SealightsHook{
			libbuildpack.DefaultHook{},
			logger,
			&libbuildpack.Command{},
		}
	})

	AfterEach(func() {
		err = os.Setenv("SL_TOKEN", token)
		Expect(err).To(BeNil())
		err = os.Setenv("SL_BUILD_SESSION_ID", build)
		Expect(err).To(BeNil())
		err = os.Setenv("SL_PROXY", proxy)
		Expect(err).To(BeNil())
		err = ioutil.WriteFile(filepath.Join(stager.BuildDir(), "Procfile"), []byte(procfile), 0755)
		Expect(err).To(BeNil())

		err = os.RemoveAll(buildDir)
		Expect(err).To(BeNil())
	})

	Describe("AfterCompile", func() {
		var (
			token = "application_token"
			bsid  = "build312"
			proxy = "http://localhost:1886"
		)
		Context("build new application run command in Procfile", func() {
			BeforeEach(func() {
				err = os.Setenv("SL_TOKEN", token)
				Expect(err).NotTo(HaveOccurred())
				err = os.Setenv("SL_BUILD_SESSION_ID", bsid)
				Expect(err).NotTo(HaveOccurred())
				err = os.Setenv("SL_PROXY", proxy)
				Expect(err).NotTo(HaveOccurred())
			})
			It("test application run cmd creation", func() {
				err = sealights.SetApplicationStart(stager)
				Expect(err).NotTo(HaveOccurred())
				bytes, err := ioutil.ReadFile(filepath.Join(stager.BuildDir(), "Procfile"))
				Expect(err).NotTo(HaveOccurred())
				cleanResult := strings.ReplaceAll(string(bytes), " ", "")
				Expect(cleanResult).To(Equal(expected))
			})
			It("hook fails with empty token", func() {
				err = os.Setenv("SL_TOKEN", "")
				Expect(err).NotTo(HaveOccurred())
				err = sealights.SetApplicationStart(stager)
				Expect(err).To(MatchError(ContainSubstring(hooks.EmptyTokenError)))
			})
			It("hook fails with empty build session id", func() {
				err = os.Setenv("SL_BUILD_SESSION_ID", "")
				Expect(err).NotTo(HaveOccurred())
				err = sealights.SetApplicationStart(stager)
				Expect(err).To(MatchError(ContainSubstring(hooks.EmptyBuildError)))
			})
		})
	})
})
