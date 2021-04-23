package db

import (
	"bitburst-assessment-task/internal/db/migrations"
	"bitburst-assessment-task/internal/db/objects"
	"context"
	"database/sql"
	"net/url"

	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/stdlib"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	pgx "github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"

	migrate "github.com/golang-migrate/migrate/v4"

	"github.com/golang-migrate/migrate/v4/database/postgres"
)

// Config contains configuration needed for constructing Postgres database
type Config struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`

	MigrationVersion int    `mapstructure:"migration_version"`
	SSLmode          string `mapstructure:"sslmode"` // enable/disable
	LogLevel         int    `mapstructure:"log_level"`
}

// DB contains Postgres database dependencies
type DB struct {
	sqlDB *sql.DB
	q     objects.Querier
}

// NewEnv initializez connection to Postgres database and injects dependencies to it, returns a structure for interacting with database.
func New(conf *Config) (db *DB, err error) {
	defer errors.Wrap(err, "db.NewEnv")

	db = &DB{}

	// constuct connection url in form of: postgresql://{username}:{password}@{host}:{port}/{database}?sslmode=true|false
	connURL := &url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(conf.Username, conf.Password),
		Host:   conf.Host + ":" + conf.Port,
		Path:   conf.Name,
	}
	connURL.Query().Add("sslmode", conf.SSLmode)

	// parse connection url
	connConfig, err := pgx.ParseConfig(connURL.String())
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse connection url")
	}

	// assign zerolog logger to pgx logger and only display error level logs
	connConfig.Logger = zerologadapter.NewLogger(log.Logger.With().Logger().Level(zerolog.Level(conf.LogLevel)))

	// open connection to postgres
	db.sqlDB = stdlib.OpenDB(*connConfig)
	err = db.sqlDB.Ping()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to ping database")
	}

	// migrate database schemas
	if err = db.migrate(conf.MigrationVersion); err != nil {
		return nil, errors.WithMessage(err, "failed to migrate database schemas")
	}

	// prepare sql statements
	if err = db.prepareStmts(); err != nil {
		return nil, errors.WithMessage(err, "failed to prepare sql statements")
	}

	return db, nil
}

// migrate migrates database schemas to specified version.
func (db *DB) migrate(version int) (err error) {
	defer errors.Wrap(err, "db.DB.migrate")

	sourceInstance, err := bindata.WithInstance(bindata.Resource(migrations.AssetNames(), migrations.Asset))
	if err != nil {
		return errors.WithMessage(err, "failed to get migrations schemas")
	}

	targetInstance, err := postgres.WithInstance(db.sqlDB, &postgres.Config{
		SchemaName:      "bitburst",
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return
	}

	m, err := migrate.NewWithInstance("go-bindata", sourceInstance, "postgres", targetInstance)
	if err != nil {
		return err
	}

	err = m.Migrate(uint(version))
	if err != nil && err != migrate.ErrNoChange {
		return errors.WithMessagef(err, "failed to migrate to version: %d", version)
	}

	return sourceInstance.Close()
}

// prepareStmts prepares sql statements
func (db *DB) prepareStmts() (err error) {
	defer errors.Wrap(err, "db.DB.prepareStmts")

	// prepare database queries
	pq, err := objects.Prepare(context.Background(), db.sqlDB)
	if err != nil {
		return err
	}

	db.q = pq

	return nil
}

// Close closes database connection
func (db *DB) Close() error {
	return errors.Wrap(db.sqlDB.Close(), "db.Env.Close")
}
