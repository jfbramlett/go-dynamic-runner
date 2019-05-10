package runner

import (
	"context"
	"fmt"
	"github.com/jfbramlett/faker/fakegen"
	"reflect"
	"testing"
	"time"
)

type Runner interface {
	Run(ctx context.Context) RunResult
}

func NewRunner(runSuiteDef RunSuiteDef, runDef RunDef, clients map[string]interface{}) (Runner, error) {
	client, found := clients[runDef.ClientClassName]
	if !found {
		return nil, fmt.Errorf("failed to find client with name %s", runDef.ClientClassName)
	}

	// get the method for this test
	method, err := getFunction(runDef.FunctionName, client)
	if err != nil {
		return nil, err
	}

	errorReturnIdx := getErrorReturnIndex(method)

	return &basicRunner{client: client, method: method, errorIdx: errorReturnIdx, runSuiteDef: runSuiteDef, runDef: runDef}, nil
}

type basicRunner struct {
	client       interface{}
	method       reflect.Value
	errorIdx     int

	runSuiteDef   RunSuiteDef
	runDef        RunDef
}

func (b *basicRunner) Run(ctx context.Context) RunResult {
	start := time.Now()
	params, err := b.getParams(ctx, b.method)
	if err != nil {
		return b.failedRun(err, start)
	}

	return b.getErrorResponse(b.method.Call(params), start)
}

func (b *basicRunner) getErrorResponse(results []reflect.Value, start time.Time) RunResult {
	if b.errorIdx >= 0 {
		errResult := results[b.errorIdx]
		if errResult.IsNil() {
			return b.passedRun(start)
		} else {
			if err, ok := errResult.Interface().(error); ok {
				return b.failedRun(err, start)
			}
		}
	}

	return b.passedRun(start)
}


func (b *basicRunner) getParams(ctx context.Context, method reflect.Value) ([]reflect.Value, error) {
	methodType := method.Type()

	argCount := methodType.NumIn()
	if methodType.IsVariadic() {
		argCount--
	}

	argDef := b.runDef.Args
	if argDef == nil {
		argDef = make(map[string]FunctionArg)
	}

	contextInterface := reflect.TypeOf((*context.Context)(nil)).Elem()

	in := make([]reflect.Value, argCount)
	for i := 0; i < argCount; i++ {
		methodArgType := methodType.In(i)

		if methodArgType.Implements(contextInterface) {
			in[i] = reflect.ValueOf(ctx)
		} else {
			newInstance, err := b.createParamValue(methodArgType)
			if err != nil {
				return []reflect.Value{}, err
			}
			in[i] = reflect.ValueOf(newInstance)
		}
	}

	return in, nil
}

func (b *basicRunner) createParamValue(argType reflect.Type) (interface{}, error) {
	newInstance := b.newInstance(argType)

	generator := b.getFakeGeneratorFor(argType)

	err := generator.FakeData(newInstance)

	return newInstance, err
}

// NewInstance creates a new instance of the "type" given
func (b *basicRunner) newInstance(typ reflect.Type) interface{} {
	if isPointer(typ) {
		return reflect.New(typ.Elem()).Interface()
	} else {
		return reflect.New(typ).Interface()
	}
}

func (b *basicRunner) getFakeGeneratorFor(argType reflect.Type) *fakegen.FakeGenerator {
	argName := getTypeName(argType)

	generator := fakegen.NewFakeGenerator()
	generator.AddFieldFilter("XXX_.*")

	tags := b.getTags(argName)
	for k, v := range tags {
		generator.AddFieldTag(k, v)
	}

	values := b.getValues(argName)
	for k, v := range values {
		generator.AddProvider(k, StaticTagProvider{val: v}.GetTaggedValue)
		generator.AddFieldTag(k, k)
	}
	return generator
}

func (b *basicRunner) getTags(argName string) map[string]string {
	tags := make(map[string]string)
	for k, v := range b.runSuiteDef.GlobalTags {
		tags[k] = v
	}

	if argDescription, found := b.runDef.Args[argName]; found && argDescription.FieldTags != nil {
		for k, v := range argDescription.FieldTags {
			tags[k] = v
		}
	}

	return tags
}

func (b *basicRunner) getValues(argName string) map[string]interface{} {
	values := make(map[string]interface{})
	for k, v := range b.runSuiteDef.GlobalValues {
		values[k] = v
	}

	if argDescription, found := b.runDef.Args[argName]; found && argDescription.ValuesOverride != nil {
		for k, v := range argDescription.ValuesOverride {
			values[k] = v
		}
	}

	return values
}



func (b *basicRunner) failedRun(err error, start time.Time) RunResult {
	return RunResult{Name: b.runDef.Name, Passed: false, Error: err, Duration: time.Since(start)}
}

func (b *basicRunner) passedRun(start time.Time) RunResult {
	return RunResult{Name: b.runDef.Name, Passed: true, Duration: time.Since(start)}
}


func isPointer(typ reflect.Type) bool {
	return typ.Kind() == reflect.Ptr
}

func getTypeName(typ reflect.Type) string {
	typeName := fmt.Sprintf("%s", typ)
	if isPointer(typ) {
		typeName = typeName[1:]
	}
	return typeName
}

func getFunction(funcName string, typ interface{}) (reflect.Value, error) {
	function := reflect.ValueOf(typ).MethodByName(funcName)
	if !function.IsValid() || function.IsNil() {
		return reflect.ValueOf(""), fmt.Errorf("failed to find method %s", funcName)
	}

	return function, nil
}

func getErrorReturnIndex(method reflect.Value) int {
	methodType := method.Type()

	errorInterface := reflect.TypeOf((*error)(nil)).Elem()

	argCount := methodType.NumOut()
	for i := 0; i < argCount; i++ {
		resultType := methodType.Out(i)
		if resultType.Kind() == reflect.Interface {
			if resultType.Implements(errorInterface) {
				return i
			}
		}
	}

	return -1
}



type StaticTagProvider struct {
val			interface{}
}

func (s StaticTagProvider) GetTaggedValue(v reflect.Value) (interface{}, error) {
	return s.val, nil
}


// Variation of our runner that runs the test as a sub-test of the given test
type testingRunner struct {
	underlying		Runner
	name			string
	mainTest		*testing.T
}

func (b *testingRunner) Run(ctx context.Context) RunResult {
	var result RunResult
	b.mainTest.Run(b.name, func(t *testing.T){
		result = b.underlying.Run(ctx)
		if !result.Passed {
			t.Fail()
		}
	})
	return result
}

func NewTestingRunner(t *testing.T, name string, underlying Runner) Runner {
	return &testingRunner{underlying: underlying, name: name, mainTest: t}
}
