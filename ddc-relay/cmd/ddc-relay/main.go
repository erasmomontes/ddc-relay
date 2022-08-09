package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/NativeMetrics/ddc-relay/pkg/agent"
	"github.com/NativeMetrics/ddc-relay/pkg/settings"
	"github.com/bugsnag/bugsnag-go"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Info().Msg(".env file not found, using environment variables")
	}

	if err = settings.Load(); err != nil {
		log.Error().
			Err(err).
			Msg("failed to load settings")

		return
	}

	// configure logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "run":
		setupBugsnag("agent")

		if err := agent.Run(); err != nil {
			log.Error().
				Err(err).
				Msg("failed to run agent")

			return
		}
	case "debug":
		handleDebug(args)
	default:
		c := []string{"run", "debug"}
		s := strings.Join(c, "\n- ")
		fmt.Printf("invalid command. \navailable commands are:\n- %s\n", s)
	}
}

func handleDebug(args []string) bool {
	return true
}

func setupBugsnag(appType string) {
	if k := settings.Values.BugsnagAPIKey; k != "" {
		bugsnag.Configure(bugsnag.Configuration{
			APIKey:          k,
			AppVersion:      os.Getenv("VERSION"),
			ReleaseStage:    "production",
			AppType:         appType,
			ProjectPackages: []string{"main", "github.com/NativeMetrics/ddc-relay"},
		})
		log.Info().Msg("bugsnag enabled")
	}
}
