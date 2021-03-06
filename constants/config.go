package constants

// Configuration-related constants
const (
	CmdRoot                 = "dsv"
	CliConfigName           = ".thy"
	CliConfigType           = "yml"
	EnvVarPrefix            = "thy"
	CmdFilePrefix           = "@"
	OutFilePrefix           = "file:"
	DontAutocompleteGlobals = false
	DefaultProfile          = "default"
	DefaultThyOneName       = "thy-one"
	DefaultCallback         = "localhost:8072"
)

// Cmd metadata
const (
	LastCommandKey = "lastCommand"
	FullCommandKey = "fullCommand"
	MainCommand    = "mainCommand"
)
