// Main binary of the MiCasa project
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/derWhity/micasa"
	"github.com/derWhity/micasa/internal/fsutils"
	"github.com/derWhity/micasa/internal/log"
	"github.com/derWhity/micasa/internal/migrate"
	kitlog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/kardianos/osext"
	_ "github.com/mattn/go-sqlite3" // Just needed for the sqlite driver
)

const (
	appName    = "MiCasa"
	appVersion = "0.1.0"
	dbFile     = "micasa.db"
)

func main() {
	execDir, err := osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}

	configFile := flag.String(
		"config",
		filepath.Join(execDir, "config.json"),
		"The configuration file to load the application's configruation from",
	)

	// Initialize the logger
	var logger log.Logger
	{
		logger = log.New(kitlog.NewLogfmtLogger(os.Stdout), log.LvlDebug)
		logger = logger.With(
			log.FldTimestamp, kitlog.DefaultTimestampUTC,
			log.FldVersion, appVersion,
		)
	}
	logger.Info(fmt.Sprintf("%s version %s is starting up...", appName, appVersion))

	// Load the main configuration file
	cs := micasa.NewConfigService(*configFile, logger)
	if err := cs.Load(); err != nil {
		logger.Error("Cannot load config. Using defaults", err)
	}

	// Load configuration and prepare data directory
	conf := cs.GetConfig()
	logger.Info(fmt.Sprintf("Using '%s' as data directory", conf.DataDir))
	fsutils.CheckAndCreateDir(conf.DataDir, logger)

	// Set up the database connection and perform pending migrations
	var db *sqlx.DB
	{
		dbFileName := path.Join(conf.DataDir, dbFile)

		if db, err = sqlx.Open("sqlite3", dbFileName); err != nil {
			logger.Crit("Failed to open database connection:", log.FldError, err)
			panic("Startup failed")
		}
		logger.Info("Performing database migrations...")
		if err = migrate.ExecuteMigrationsOnDb(db, logger); err != nil {
			logger.Crit("Database migration has failed", log.FldError, err)
			panic("Cannot continue. Please check database for consistency and try again")
		}
	}
}
