package runner

import (
	"testing"
)

type RunnerFactory interface {
	GetRunner(runSuiteDef RunSuiteDef, runDef RunDef, clients map[string]interface{}) (Runner, error)
}

// runner factory that just creates a basic runner instance
type defaultRunnerFactory struct {}

func (d *defaultRunnerFactory) GetRunner(runSuiteDef RunSuiteDef, runDef RunDef, clients map[string]interface{}) (Runner, error) {
	return NewRunner(runSuiteDef, runDef, clients)
}

// runner factory that creates new runner instances wrapped in standard Go testing
type testingRunnerFactory struct {
	mainTest		*testing.T
}

func (d *testingRunnerFactory) GetRunner(runSuiteDef RunSuiteDef, runDef RunDef, clients map[string]interface{}) (Runner, error) {
	underlying, err := NewRunner(runSuiteDef, runDef, clients)
	if err != nil {
		return nil, err
	}
	return NewTestingRunner(d.mainTest, runDef.Name, underlying), nil
}

