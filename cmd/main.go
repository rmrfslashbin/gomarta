package main

import (
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/davecgh/go-spew/spew"
	"github.com/rmrfslashbin/gomarta/pkg/buses"
	"github.com/rmrfslashbin/gomarta/pkg/gtfspec"
	"github.com/rs/zerolog"
)

const (
	// APP_NAME is the name of the application
	APP_NAME = "gomarta"

	// CONFIG_FILE is the name of the config file
	CONFIG_FILE = "config.yaml"
)

// Context is used to pass context/global configs to the commands
type Context struct {
	// log is the logger
	log *zerolog.Logger
}

// ConfigSetCmd sets a config value
type BusCmd struct {
	Vehicles string `name:"vehicles" default:"https://gtfs-rt.itsmarta.com/TMGTFSRealTimeWebService/vehicle" help:"URL for the Marta Bus Vehicles GTFS endpoint."`
	Trips    string `name:"trips" default:"https://gtfs-rt.itsmarta.com/TMGTFSRealTimeWebService/tripupdate" help:"URL for the Marta Bus Trips GTFS endpoint."`
}

// Run is the entry point for the BusCmd command
func (r *BusCmd) Run(ctx *Context) error {
	vehicles, err := buses.GetData(&buses.Input{
		Url: r.Vehicles,
		Log: ctx.log,
	})
	if err != nil {
		return err
	}

	for _, vehicle := range vehicles {
		spew.Dump(vehicle)
	}

	return nil
}

// UpdateSpecsCmd updates the GTFS feed specs
type UpdateSpecsCmd struct {
	Url  string  `name:"url" default:"https://itsmarta.com/google_transit_feed/google_transit.zip" help:"URL the GTFS feed spec zip file."`
	Json *string `name:"json" group:"output" required:"" xor:"output" help:"Output the specs as JSON to a file."`
	Gob  *string `name:"gob" group:"output" required:"" xor:"output" help:"Output the specs as Gob to a file."`
}

// Run is the entry point for the UpdateSpecsCmd command
func (r *UpdateSpecsCmd) Run(ctx *Context) error {
	spec, err := gtfspec.Update(&gtfspec.Input{
		Url: r.Url,
		Log: ctx.log,
	})
	if err != nil {
		return err
	}
	if r.Json != nil {
		fpqn := filepath.Clean(*r.Json)
		if err := spec.ToJSONFile(fpqn); err != nil {
			return err
		}
	}
	if r.Gob != nil {
		fpqn := filepath.Clean(*r.Gob)
		if err := spec.ToGOBFile(fpqn); err != nil {
			return err
		}
	}

	return nil
}

// CLI is the main CLI struct
type CLI struct {
	// Global flags/args
	LogLevel string `name:"loglevel" env:"LOGLEVEL" default:"debug" enum:"panic,fatal,error,warn,info,debug,trace" help:"Set the log level."`

	Bus    BusCmd         `cmd:"" help:"Get bus data."`
	Update UpdateSpecsCmd `cmd:"" help:"Update the GTFS feed specs."`
}

func main() {
	var err error

	// Set up the logger
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Parse the command line
	var cli CLI
	ctx := kong.Parse(&cli)

	// Set up the logger's log level
	// Default to info via the CLI args
	switch cli.LogLevel {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	// Log some start up stuff for debugging
	log.Info().
		Str("app_name", APP_NAME).
		Str("log_level", cli.LogLevel).
		Msg("starting up")

	// Call the Run() method of the selected parsed command.
	err = ctx.Run(&Context{log: &log})

	// FatalIfErrorf terminates with an error message if err != nil
	ctx.FatalIfErrorf(err)
}
