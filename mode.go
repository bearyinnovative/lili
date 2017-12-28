package main

import "os"

const (
	UnknownMode = "unknown"
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)
const (
	unknownCode = iota
	debugCode
	releaseCode
	testCode
)

var mode = debugCode
var modeName = DebugMode

func init() {
	mode := os.Getenv("LILI_MODE")
	SetMode(mode)
}

func SetMode(value string) {
	switch value {
	case "":
		mode = unknownCode
	case DebugMode:
		mode = debugCode
	case ReleaseMode:
		mode = releaseCode
	case TestMode:
		mode = testCode
	default:
		panic("lili mode type error: " + value)
	}
	modeName = value
}

func Mode() string {
	return modeName
}

func IsDebug() bool {
	return mode == debugCode
}
