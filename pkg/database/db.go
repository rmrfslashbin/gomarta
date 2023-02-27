package database

import (
	"os"
	"path/filepath"

	"github.com/rmrfslashbin/gomarta/pkg/gtfspec"
	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Options for the bus instance
type Option func(c *Database)

// Database for the app instance
type Database struct {
	log    *zerolog.Logger
	sqlite *string
	mysql  *string
	pgsql  *string
	db     *gorm.DB
}

// New creates a new mastoclinet instance
func New(opts ...Option) (*Database, error) {
	cfg := &Database{}

	// apply the list of options to Bus
	for _, opt := range opts {
		opt(cfg)
	}

	// set up logger if not provided
	if cfg.log == nil {
		log := zerolog.New(os.Stderr).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		cfg.log = &log
	}

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	if cfg.sqlite != nil {
		fqpn := filepath.Clean(*cfg.sqlite)
		cfg.log.Info().Str("dbFile", fqpn).Msg("using sqlite database")
		db, err := gorm.Open(sqlite.Open(fqpn), config)
		if err != nil {
			return nil, &ErrSqliteOpen{Err: err, Filename: fqpn}
		}
		cfg.db = db
	} else if cfg.mysql != nil {
		db, err := gorm.Open(mysql.Open(*cfg.mysql), config)
		if err != nil {
			return nil, &ErrMySqlOpen{Err: err, Dsn: *cfg.mysql}
		}
		cfg.db = db
	} else if cfg.pgsql != nil {
		db, err := gorm.Open(postgres.Open(*cfg.pgsql), config)
		if err != nil {
			return nil, &ErrPGSqlOpen{Err: err, Dsn: *cfg.pgsql}
		}
		cfg.db = db
	} else {
		return nil, &ErrNoDatabase{}
	}

	if err := cfg.db.AutoMigrate(
		&gtfspec.Agency{},
		&gtfspec.Calendar{},
		&gtfspec.CalendarDate{},
		&gtfspec.Route{},
		&gtfspec.Shape{},
		&gtfspec.Stop{},
		&gtfspec.StopTime{},
		&gtfspec.Trip{},
	); err != nil {
		return nil, err
	}

	return cfg, nil
}

// WithLogger sets the logger for the database instance
func WithLogger(log *zerolog.Logger) Option {
	return func(c *Database) {
		c.log = log
	}
}

// WithSqlite sets the sqlite connection string for the database instance
func WithSqlite(sqlite *string) Option {
	return func(c *Database) {
		c.sqlite = sqlite
	}
}

// WithMysql sets the mysql connection string for the database instance
func WithMysql(mysql *string) Option {
	return func(c *Database) {
		c.mysql = mysql
	}
}

// WithPgsql sets the pgsql connection string for the database instance
func WithPgsql(pgsql *string) Option {
	return func(c *Database) {
		c.pgsql = pgsql
	}
}

func (d *Database) Create(value interface{}) (*gorm.DB, error) {
	tx := d.db.CreateInBatches(value, 100)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (d *Database) GetAgency(agencyId string) (*gtfspec.Agency, error) {
	agency := &gtfspec.Agency{}
	if err := d.db.First(agency, "agency_id = ?", agencyId).Error; err != nil {
		return nil, err
	}
	return agency, nil
}

func (d *Database) GetRoute(routeId int) (*gtfspec.Route, error) {
	route := &gtfspec.Route{}
	if err := d.db.First(route, "route_id = ?", routeId).Error; err != nil {
		return nil, err
	}
	return route, nil
}

func (d *Database) GetStop(stopId int) (*gtfspec.Stop, error) {
	stop := &gtfspec.Stop{}
	if err := d.db.First(stop, "stop_id = ?", stopId).Error; err != nil {
		return nil, err
	}
	return stop, nil
}

func (d *Database) GetTrip(tripId int, RouteId int) (*gtfspec.Trip, error) {
	trip := &gtfspec.Trip{}
	if err := d.db.First(trip, "trip_id = ? AND route_id = ?", tripId, RouteId).Error; err != nil {
		return nil, err
	}
	return trip, nil
}
