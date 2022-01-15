package providers

import (
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/sirupsen/logrus"
)

type DefaultContextProvider struct {
	buildContext  *contexts.BuildContext
	deployContext *contexts.DeployContext
	logger        *logrus.Logger
}

func NewContextProvider(logger *logrus.Logger) *DefaultContextProvider {
	return &DefaultContextProvider{
		buildContext:  contexts.CreateBuildContext(),
		deployContext: &contexts.DeployContext{},
		logger:        logger,
	}
}

func (p *DefaultContextProvider) GetBuildContext() *contexts.BuildContext {
	return p.buildContext
}
func (p *DefaultContextProvider) GetDeployContext() *contexts.DeployContext {
	return p.deployContext
}
