package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/calmitchell617/reserva/internal/data"
	"github.com/calmitchell617/reserva/internal/jsonlog"
	"github.com/calmitchell617/reserva/internal/mailer"
	"github.com/calmitchell617/reserva/internal/vcs" // New import

	_ "github.com/lib/pq"
)

var (
	version = vcs.Version()
)

type config struct {
	port int
	env  string
	db   struct {
		write struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  string
		}
		read struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  string
		}
	}
	limiter struct {
		enabled bool
		rps     float64
		burst   int
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}

	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 80, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.write.dsn, "write-db-dsn", "", "PostgreSQL write DSN")
	flag.IntVar(&cfg.db.write.maxOpenConns, "write-db-max-open-conns", 25, "PostgreSQL write max open connections")
	flag.IntVar(&cfg.db.write.maxIdleConns, "write-db-max-idle-conns", 25, "PostgreSQL write max idle connections")
	flag.StringVar(&cfg.db.write.maxIdleTime, "write-db-max-idle-time", "15m", "PostgreSQL write max connection idle time")

	flag.StringVar(&cfg.db.read.dsn, "read-db-dsn", "", "PostgreSQL read DSN")
	flag.IntVar(&cfg.db.read.maxOpenConns, "read-db-max-open-conns", 25, "PostgreSQL read max open connections")
	flag.IntVar(&cfg.db.read.maxIdleConns, "read-db-max-idle-conns", 25, "PostgreSQL read max idle connections")
	flag.StringVar(&cfg.db.read.maxIdleTime, "read-db-max-idle-time", "15m", "PostgreSQL read max connection idle time")

	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 0, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "", "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	if cfg.db.write.dsn == "" {
		fmt.Println("You must enter a write node DSN to start the server. It can be the same as the write node.")
		os.Exit(1)
	}

	if cfg.db.read.dsn == "" {
		fmt.Println("You must enter a read node DSN to start the server. It can be the same as the write node.")
		os.Exit(1)
	}

	if cfg.smtp.host == "" {
		fmt.Println("You must enter an SMTP host to start the server")
		os.Exit(1)
	}
	if cfg.smtp.port == 0 {
		fmt.Println("You must enter an SMTP port to start the server")
		os.Exit(1)
	}
	if cfg.smtp.username == "" {
		fmt.Println("You must enter an SMTP username to start the server")
		os.Exit(1)
	}
	if cfg.smtp.password == "" {
		fmt.Println("You must enter an SMTP password to start the server")
		os.Exit(1)
	}

	if cfg.smtp.sender == "" {
		fmt.Println("You must enter an SMTP sender to start the server")
		os.Exit(1)
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	writeDb, readDb, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer writeDb.Close()
	defer readDb.Close()

	logger.PrintInfo("database connection pool established", nil)

	expvar.NewString("version").Set(version)

	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("write database", expvar.Func(func() interface{} {
		return writeDb.Stats()
	}))

	expvar.Publish("read database", expvar.Func(func() interface{} {
		return readDb.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(writeDb, readDb),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender, cfg.env),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, *sql.DB, error) {
	writeDb, err := sql.Open("postgres", cfg.db.write.dsn)
	if err != nil {
		return nil, nil, err
	}

	writeDb.SetMaxOpenConns(cfg.db.write.maxOpenConns)
	writeDb.SetMaxIdleConns(cfg.db.write.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.write.maxIdleTime)
	if err != nil {
		return nil, nil, err
	}

	writeDb.SetConnMaxIdleTime(duration)

	readDb, err := sql.Open("postgres", cfg.db.read.dsn)
	if err != nil {
		return nil, nil, err
	}

	readDb.SetMaxOpenConns(cfg.db.read.maxOpenConns)
	readDb.SetMaxIdleConns(cfg.db.read.maxIdleConns)

	readDb.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = writeDb.PingContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	err = readDb.PingContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	return writeDb, readDb, nil
}
