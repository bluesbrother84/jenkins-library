package aakaas

import (
	"time"

	abapbuild "github.com/SAP/jenkins-library/pkg/abap/build"
	"github.com/SAP/jenkins-library/pkg/abaputils"
	"github.com/SAP/jenkins-library/pkg/command"
)

type AakUtils interface {
	command.ExecRunner
	abapbuild.HTTPSendLoader
	ReadAddonDescriptor(FileName string) (abaputils.AddonDescriptor, error)
	GetMaxRuntime() time.Duration
	GetPollingInterval() time.Duration
}
