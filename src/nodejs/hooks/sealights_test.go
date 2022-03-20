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
)

var _ = Describe("Sealights hook", func() {
	var (
		err       error
		logger    *libbuildpack.Logger
		stager    *libbuildpack.Stager
		buffer    *bytes.Buffer
		sealights hooks.SealightsHook
		token     string
		build     string
		proxy     string
		expected  = "node ./node_modules/.bin/slnodejs run  --useinitialcolor true --token application_token --buildsessionid build312 --proxy http://localhost:1886 index.js"
	)

	BeforeEach(func() {
		token = os.Getenv("SL_TOKEN")
		build = os.Getenv("SL_BUILD_SESSION_ID")
		proxy = os.Getenv("SL_PROXY")
		buffer = new(bytes.Buffer)
		logger = libbuildpack.NewLogger(buffer)
	})

	AfterEach(func() {
		err = os.Setenv("SL_TOKEN", token)
		Expect(err).NotTo(HaveOccurred())
		err = os.Setenv("SL_BUILD_SESSION_ID", build)
		Expect(err).NotTo(HaveOccurred())
		err = os.Setenv("SL_PROXY", proxy)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("AfterCompile", func() {
		var (
			token = "application_token"
			bsid  = "build312"
			proxy = "http://localhost:1886"
		)
		Context("VCAP_SERVICES contains seeker service - as a user provided service", func() {
			BeforeEach(func() {
				err = os.Setenv("SL_TOKEN", token)
				Expect(err).NotTo(HaveOccurred())
				err = os.Setenv("SL_BUILD_SESSION_ID", bsid)
				Expect(err).NotTo(HaveOccurred())
				err = os.Setenv("SL_PROXY", proxy)
				Expect(err).NotTo(HaveOccurred())
			})
			It("test application run cmd creation", func() {
				err = sealights.AfterCompile(stager)
				Expect(err).NotTo(HaveOccurred())
				bytes, err := ioutil.ReadFile(filepath.Join(stager.BuildDir(), "Procfile"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(bytes)).To(Equal(expected))
			})
			It("hook fails with empty token", func() {
				err = os.Setenv("SL_TOKEN", "")
				Expect(err).NotTo(HaveOccurred())
				err = sealights.AfterCompile(stager)
				Expect(err).To(MatchError(ContainSubstring(sealights.EmptyTokenError)))
			})
			It("hook fails with empty build session id", func() {
				err = os.Setenv("SL_BUILD_SESSION_ID", "")
				Expect(err).NotTo(HaveOccurred())
				err = sealights.AfterCompile(stager)
				Expect(err).To(MatchError(ContainSubstring(sealights.EmptyBuildError)))
			})
		})
	})
})
