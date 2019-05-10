package runner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
)

// function used to build a run suite from a file, the file is a JSON file containing the definition of what to run
func buildRunSuiteFromFile(configFile string) (RunSuiteDef, error) {
	runSuiteDef := RunSuiteDef{}
	runSuiteDefTxt, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalln(err)
		return RunSuiteDef{}, err
	}

	err = json.Unmarshal(runSuiteDefTxt, &runSuiteDef)
	if err != nil {
		log.Fatalln(err)
		return RunSuiteDef{}, err
	}

	return runSuiteDef, nil
}

// function used to build a run suite given a type - this uses reflection to identify the methods to wrap
func buildRunSuiteFromType(interfaceType reflect.Type, globalValues map[string]interface{}, globalTags map[string]string, excludes []string, clients map[string]interface{}) (RunSuiteDef, error) {
	runSuite := RunSuiteDef{Tests: make([]RunDef, 0), GlobalValues: globalValues, GlobalTags: globalTags}

	for i := 0; i < interfaceType.NumMethod(); i++ {
		methodName := interfaceType.Method(i).Name
		for _, excludedMethod := range excludes {
			if methodName == excludedMethod {
				continue
			}
		}
		runSuite.Tests = append(runSuite.Tests, RunDef{Name: fmt.Sprintf(" Test %T.%s", interfaceType, methodName),
			FunctionName: methodName,
			ClientClassName: "",
			Args: make(map[string]FunctionArg),
		})

	}
	return runSuite, nil
}
