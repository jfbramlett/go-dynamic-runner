package runner

import (
	"context"
	"reflect"
	"testing"
)

type SuiteRunner interface {
	Run() RunSuiteResult
}

// simple runner that just runs the tests and reports the results
type basicSuiteRunner struct {
	runSuiteDef   	RunSuiteDef
	runnerFactory 	RunnerFactory
	clients			map[string]interface{}
}

func (f *basicSuiteRunner) Run() RunSuiteResult {
	runSuiteResult := RunSuiteResult{}

	if f.runSuiteDef.Repeat == 0 {
		f.runSuiteDef.Repeat = 1
	}

	// prep the runners to run
	runners := []Runner{}
	for _, runDef := range f.runSuiteDef.Tests {
		runner, err := f.runnerFactory.GetRunner(f.runSuiteDef, runDef, f.clients)
		if err == nil {
			runners = append(runners, runner)
		}
	}

	runCount := 0
	for runCount < f.runSuiteDef.Repeat {
		for _, runner := range runners {
			runSuiteResult.TotalTests++
			result := runner.Run(context.Background())
			runSuiteResult.Duration = runSuiteResult.Duration + result.Duration
			if result.Passed {
				runSuiteResult.Passed++
			} else {
				runSuiteResult.Failed++
			}
		}
		runCount++
	}

	return runSuiteResult
}

// constructor to build a new RunSuite (set of things to execute). This builds it from a JSON-based config
func NewRunSuite(configFile string, clients map[string]interface{}) (SuiteRunner, error) {
	return NewTestingRunSuite(nil, configFile, clients)
}

// constructor to build a new RunSuite (set of things to execute). This builds it from a JSON-based config. Each RunDef
// when executing will be run as a Go Test
func NewTestingRunSuite(t *testing.T, configFile string, clients map[string]interface{}) (SuiteRunner, error) {
	runSuite, err := buildRunSuiteFromFile(configFile)
	if err != nil {
		return nil, err
	}

	return newRunSuite(runSuite, t, clients), nil
}

// constructor method for creating a new RunSuite, a run suite represents a set of run defs (or things to run), this builds
// the suite automatically based on the type
func NewAutoRunSuite(interfaceType reflect.Type, globalValues map[string]interface{}, globalTags map[string]string, excludes []string, clients map[string]interface{}) (SuiteRunner, error) {
	return NewAutoTestingRunSuite(nil, interfaceType, globalValues, globalTags, excludes, clients)
}

// constructor method for creating a new RunSuite, a run suite represents a set of run defs (or things to run), this builds
// the suite automatically based on the type. Each execution will be wrapped as a Go test case.
func NewAutoTestingRunSuite(t *testing.T, interfaceType reflect.Type, globalValues map[string]interface{}, globalTags map[string]string, excludes []string, clients map[string]interface{}) (SuiteRunner, error) {
	runSuite, err := buildRunSuiteFromType(interfaceType, globalValues, globalTags, excludes, clients)
	if err != nil {
		return nil, err
	}

	return newRunSuite(runSuite, t, clients), nil
}

func newRunSuite(runSuite RunSuiteDef, t *testing.T, clients map[string]interface{}) SuiteRunner {
	if t != nil {
		return &basicSuiteRunner{runSuiteDef: runSuite, runnerFactory: &testingRunnerFactory{mainTest: t}, clients: clients}
	} else {
		return &basicSuiteRunner{runSuiteDef: runSuite, runnerFactory: &defaultRunnerFactory{}, clients: clients}
	}
}
