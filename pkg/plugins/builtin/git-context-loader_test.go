package builtin

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/contexts/contextsfakes"
	"github.com/rdrdog/bldr/pkg/lib"
	"github.com/rdrdog/bldr/pkg/lib/libfakes"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func getConfiguredGitContextLoader(t *testing.T) *GitContextLoader {
	logger, _ := test.NewNullLogger()
	p := &GitContextLoader{}

	p.logger = logger
	p.configuration = &config.Configuration{
		CI: false,
		Git: config.GitConfig{
			MainBranchName: "main",
		},
	}

	return p
}

func TestGitContextLoaderExecute(t *testing.T) {
	spec.Run(t, "GitContextLoader.Execute", testGitContextLoaderExecute, spec.Report(report.Terminal{}))
}

func testGitContextLoaderExecute(t *testing.T, when spec.G, it spec.S) {
	var p *GitContextLoader
	var fakeContextProvider *contextsfakes.FakeContextProvider
	var fakeLibProvider *libfakes.FakeLibProvider
	var fakeGitLib *libfakes.FakeGit
	var bc *contexts.BuildContext

	it.Before(func() {
		RegisterTestingT(t)
		p = getConfiguredGitContextLoader(t)

		fakeContextProvider = &contextsfakes.FakeContextProvider{}
		fakeLibProvider = &libfakes.FakeLibProvider{}
		fakeGitLib = &libfakes.FakeGit{}

		fakeLibProvider.GetGitLibReturns(fakeGitLib)

		bc = &contexts.BuildContext{
			GitContext: &contexts.GitContext{},
		}
		gs := &lib.GitState{
			CommitSha:              "8f825b7454e61e5c61fcf165e8299151b46d67f8",
			BranchName:             "main",
			MainBranchForkPoint:    "1231231231232135c61fcf165e8299151b46d67f8",
			ChangesSinceMainBranch: []string{"abc123"},
		}
		fakeContextProvider.GetBuildContextReturns(bc)
		fakeLibProvider.GetGitLibReturns(fakeGitLib)
		fakeGitLib.LoadRepoInformationReturns(gs)
	})

	it("execute method is called without any error and build context has correct values", func() {

		err := p.Execute(fakeContextProvider, nil, fakeLibProvider)
		assert.Nil(t, err)

		Expect(bc.GitContext.BranchName).To(Equal("main"))
		Expect(bc.GitContext.FullCommitSha).To(Equal("8f825b7454e61e5c61fcf165e8299151b46d67f8"))
		Expect(bc.GitContext.ShortCommitSha).To(Equal("8f825b7"))
		Expect(bc.GitContext.MainBranchForkPointShort).To(Equal("1231231"))
		Expect(bc.GitContext.ChangesSinceMainBranch).To(Equal([]string{"abc123"}))
	})
}
