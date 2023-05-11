package service

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type RunContainer interface {
	RunProgram(resp container.CreateResponse, input string, binaryName string) (string, error)
	CompileProgram(resp container.CreateResponse, compileCommand string) (types.IDResponse, error)
	CreateContainer(image string, AbsPath string) (container.CreateResponse, error)
	RemoveContainer(resp container.CreateResponse)
}

type Service struct {
	RunContainer
}

func NewService() *Service {
	return &Service{
		RunContainer: NewClient(),
	}
}
