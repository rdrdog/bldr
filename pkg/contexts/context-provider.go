package contexts

import (
	"github.com/sirupsen/logrus"
)

type ContextProvider struct {
	BuildContext  *BuildContext
	DeployContext *DeployContext
	logger        *logrus.Logger
}

func NewContextProvider(logger *logrus.Logger) *ContextProvider {
	return &ContextProvider{
		BuildContext:  CreateBuildContext(),
		DeployContext: &DeployContext{},
		logger:        logger,
	}
}
