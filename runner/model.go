package runner

import "time"

type FunctionArg struct {
	ValuesOverride			map[string]interface{}	`json:"valuesOverride"`
	FieldTags				map[string]string		`json:"fieldTags"`
}


type RunDef struct {
	Name            string                 `json:"name"`
	ClientClassName string                 `json:"clientClassName"`
	FunctionName	string				   `json:"functionName"`
	Args			map[string]FunctionArg `json:"args"`
}


type RunSuiteDef struct {
	Tests			[]RunDef              	`json:"runDefinitions"`
	GlobalValues	map[string]interface{} 	`json:"globalValues"`
	GlobalTags		map[string]string    	`json:"globalTags"`
	Repeat			int						`json:"repeats"`
}


type RunResult struct {
	Name		string
	Passed		bool
	Error		error
	Duration	time.Duration
}

type RunSuiteResult struct {
	TotalTests		int
	Passed			int
	Failed			int
	Duration		time.Duration
}