package main

import (
	"bitburst-assessment-task/internal/client"
	"bitburst-assessment-task/internal/db"
	"bitburst-assessment-task/internal/server"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// config holds configuration that comes from command flags or env vars, and is needed to configure and start the app
type config struct {
	Log struct {
		// Local path of a filename for storing logs
		Path string `mapstructure:"path"`

		// Logging level:
		// -1 for TRACE, 0 for DEBUG, 1 for INFO, 2 for WARNING, 3 for ERROR, 4 for FATAL, 5 for PANIC
		Level int `mapstructure:"level"`

		// Beautify if set to true format logs in a beautiful way instead of default json formatting, specifically intended for console
		Beautify bool `mapstructure:"beautify"`

		// file for storing logs
		file *os.File `mapstructure:"-"`
	} `mapstructure:"log"`

	Server server.Config `mapstructure:"server"`

	Client client.Config `mapstructure:"client"`

	Database db.Config `mapstructure:"database"`
}

// setConfigDetails sets command flags, defaults and envs and default config file details
func setConfigDetails(v *viper.Viper, p *pflag.FlagSet) {
	// look for config file in following directories,
	// if config path isn't supplied from command args
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("./config/")

	// replace viper keys char to env ones,
	// eg: viperkey=log.path -> env=LOG_PATH
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// only use envs that start with BITBURST prefix
	v.SetEnvPrefix("bitburst")

	v.AutomaticEnv()

	// set command flags, their defaults and bind envs

	// for logs
	p.String("log-path", "./logs.jsonl", "local path of a filename for storing logs")
	v.BindPFlag("log.path", p.Lookup("log-path"))
	v.SetDefault("log.path", "./logs.jsonl")

	p.Int("log-level", 1, "output logs that are higher than or equal to specified level: -1 for TRACE, 0 for DEBUG, 1 for INFO, 2 for WARNING, 3 for ERROR, 4 for FATAL, 5 for PANIC")
	v.BindPFlag("log.level", p.Lookup("log-level"))
	v.SetDefault("log.level", 1)

	p.Bool("log-beautify", false, "format logs in a beautiful way instead of default json formatting, specifically intended for console")
	v.BindPFlag("log.beautify", p.Lookup("log-beautify"))
	v.SetDefault("log.beautify", false)

	// for server
	p.StringP("server-listen-address", "l", "0.0.0.0:9090", "listen address for http server, port must be included")
	_ = v.BindPFlag("server.listen_address", p.Lookup("server-listen-address"))
	v.BindEnv("server.listen_address", "SERVER_LISTEN_ADDRESS")
	v.SetDefault("server.listen_address", "0.0.0.0:9090")

	p.Duration("server-read-timeout", 0, "timeout duration for the server to read request body")
	_ = v.BindPFlag("server.read_timeout", p.Lookup("server-read-timeout"))
	v.SetDefault("server.read_timeout", 0)

	p.Duration("server-write-timeout", 0, "timeout duration for the server to write response")
	_ = v.BindPFlag("server.write_timeout", p.Lookup("server-write-timeout"))
	v.SetDefault("server.write_timeout", 0)

	p.Duration("server-shutdown-timeout", time.Second*5, "timeout duration for the server to shutdown")
	_ = v.BindPFlag("server.shutdown_timeout", p.Lookup("server-shutdown-timeout"))
	v.SetDefault("server.shutdown_timeout", time.Second*5)

	// for client
	p.String("client-tester-service-address", "127.0.0.1:9010", "listen address of tester service")
	_ = v.BindPFlag("client.tester_service_address", p.Lookup("client-tester-service-address"))
	v.BindEnv("client.tester_service_address", "CLIENT_TESTER_SERVICE_ADDRESS")
	v.SetDefault("client.tester_service_address", "127.0.0.1:9010")

	// for database
	p.StringP("database-host", "h", "127.0.0.1", "database host")
	v.BindPFlag("database.host", p.Lookup("database-host"))
	v.BindEnv("database.host", "DATABASE_HOST")
	v.SetDefault("database.host", "127.0.0.1")

	p.StringP("database-port", "p", "5432", "database port")
	v.BindPFlag("database.port", p.Lookup("database-port"))
	v.BindEnv("database.port", "DATABASE_PORT")
	v.SetDefault("database.port", "5432")

	p.StringP("database-username", "u", "postgres", "database username")
	v.BindPFlag("database.username", p.Lookup("database-username"))
	v.BindEnv("database.username", "DATABASE_USERNAME")
	v.SetDefault("database.username", "postgres")

	p.String("database-password", "postgres", "database password")
	v.BindPFlag("database.password", p.Lookup("database-password"))
	v.BindEnv("database.password", "DATABASE_PASSWORD")
	v.SetDefault("database.password", "postgres")

	p.StringP("database-name", "n", "postgres", "database name")
	v.BindPFlag("database.name", p.Lookup("database-name"))
	v.BindEnv("database.name", "DATABASE_NAME")
	v.SetDefault("database.name", "postgres")

	p.String("database-sslmode", "disable", "database sslmode")
	v.BindPFlag("database.sslmode", p.Lookup("database-sslmode"))
	v.SetDefault("database.sslmode", "disable")

	p.Int("database-migration-version", 3, "database migration version")
	v.BindPFlag("database.migration_version", p.Lookup("database-migration-version"))
	v.SetDefault("database.migration_version", 3)
}

// Will be set using ldflags
var (
	version    string
	commitHash string
	buildDate  string
)

const appName = "Bitburst Assessment Task"

func main() {
	retcode := 0
	defer func() { os.Exit(retcode) }() // call os.Exit at the end of main, so other deferred calls don't get discarded

	v, p := viper.New(), pflag.NewFlagSet(appName, pflag.ExitOnError)

	setConfigDetails(v, p)

	p.StringP("config-path", "c", "", "local path to configuration file")
	p.Bool("version", false, "show version information")

	_ = p.Parse(os.Args[1:])

	// check if version wants to be printed
	if v, _ := p.GetBool("version"); v {
		fmt.Printf("%s version %s (%s) built on %s\n", appName, version, commitHash, buildDate)

		retcode = 0
		return
	}

	// check if config path is supplied
	configPath, _ := p.GetString("config-path")
	if configPath != "" {
		v.SetConfigFile(configPath)
	}

	var err error
	// try to read config file
	var conf config
	if err = v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// config file was found, but there were some errors while reading it
			if !os.IsNotExist(err) {
				log.Logger.Fatal().Err(err).Msg("failed to read config file")
			}
		}

		// config file not found, manually fill config from envs or cmd flags
		if configPath != "" {
			log.Logger.Warn().Msg("failed to find config file")
		}
	}

	// unmarshal application configuration from conf file into conf struct
	if err := v.Unmarshal(&conf); err != nil {
		log.Logger.Err(err).Msg("failed to unmarshal config file")
		retcode = -1
		return
	}

	// create new file for storing logs
	conf.Log.file, err = os.Create(conf.Log.Path)
	if err != nil {
		log.Logger.Err(err).Msg("failed to create log file")
		retcode = -1
		return
	}
	defer func() {
		if err := conf.Log.file.Close(); err != nil {
			// output file close error to stdout, because file will be unusable
			logger := log.Logger.Output(os.Stdout).With().Logger()
			logger.Err(err).Str("logs-path", conf.Log.Path).Msg("FAILED to close logs file")
		}
	}()

	// output logs both in terminal and file
	// and set beautiful logging for terminals
	var logsWriters io.Writer
	if conf.Log.Beautify {
		logsWriters = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.Stamp,
		}
	} else {
		logsWriters = os.Stdout
	}
	multi := zerolog.MultiLevelWriter(logsWriters, conf.Log.file)

	log.Logger = log.Output(multi).With().Logger()

	// set verbosity level
	zerolog.SetGlobalLevel(zerolog.Level(conf.Log.Level))

	// set up database
	log.Logger.Info().Msg("connecting to database")
	var (
		database *db.DB
		ok       bool
	)
	for i := 0; i < 3; i++ {
		database, err = db.New(&conf.Database, &log.Logger)
		if err != nil {
			// sometimes database may be off for some secs, so retrying is a good practice
			if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "failed to ping database") {
				log.Logger.Warn().Err(err).Msg("couldn't establish database connection, retrying in 5 secs...")
				time.Sleep(5 * time.Second)
				continue
			} else {
				log.Logger.Err(err).Msg("failed to establish database connection")
				retcode = -1
				return
			}
		}
		ok = true
		break
	}
	if !ok {
		log.Logger.Err(err).Msg("failed to establish database connection, database is off")
		retcode = -1
		return
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Logger.Warn().Err(err).Msg("failed to close database connection")
			retcode = -1
		}
	}()

	// set up client
	cli := client.New(&conf.Client)

	// set up server
	srv := server.New(&conf.Server, database, cli)

	// start the server
	log.Logger.Info().Str("listen-address", conf.Server.ListenAddress).Msg("starting the server")
	srvErrChan := make(chan error)
	go srv.Start(srvErrChan)

	// run a background job that will delete objects that weren't seen for 30 seconds
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go database.DeleteNotSeenObjects(ctx)

	// catch interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case q := <-quit:
		log.Logger.Info().Str("signal", q.String()).Msg("received signal, closing server and other opened resources")

		// close server
		if err := srv.Close(); err != nil {
			log.Logger.Warn().Err(err).Msg("failed to close the server")
			retcode = -1
		}
	case err := <-srvErrChan:
		log.Logger.Err(err).Msg("failed to start the server, closing other opened resources")
	}
}
