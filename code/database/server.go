package server


import (
	// built-in
	"context"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"

	// 3rd party
	// "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/log15adapter"
	"github.com/jackc/pgx/v4/pgxpool"
	log "gopkg.in/inconshreveable/log15.v2"

	//custom
	"masterthesis/custom"
)


func (postgres *PostgresSQL) Insert(args *InsertRequest, reply *InsertResponse) error {

	_, err := postgres.db.Exec(context.Background(), args.Query)

	if err != nil {
		log.Crit("Unable to insert records", "error", err)
		return err
	}

	return nil
}


func (postgres *PostgresSQL) Select(args *SelectRequest, reply *SelectResponse) error {
	_, err := postgres.db.Exec(context.Background(), args.Query)

	if err != nil {
		log.Crit("Unable to select records", "error", err)
		return err
	}
	
	return nil
}


func (postgres *PostgresSQL) serve() {
	// register master
	rpc.Register(postgres)
	
	// handle HTTP requests
	rpc.HandleHTTP()

	// get socket name
	socketName := os.Getenv("database_socket_name")
	
	// remove previously created socket
	os.Remove(socketName)

	// listen to the socket
	listener, err := net.Listen("unix", socketName)

	// check for error
	if err != nil {
		log.Crit("listen error:", err)
	}

	// initialize thread for listener
	go http.Serve(listener, nil)
}


func createConnection(databaseUrl string) *pgxpool.Pool {

	logger := log15adapter.NewLogger(log.New("module", "pgx"))

	poolConfig, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		log.Crit("Unable to parse DATABASE_URL", "error", err)
		os.Exit(1)
	}

	poolConfig.ConnConfig.Logger = logger

	db, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Crit("Unable to create connection pool", "error", err)
		os.Exit(1)
	}

	return db
}


func MakeDatabase() *PostgresSQL {

	// load environment
	custom.LoadEnvironment()

	// establish connection url
	user 	:= os.Getenv("DB_USER")
	pass	:= os.Getenv("DB_PASS")
	host 	:= os.Getenv("DB_HOST")
	port 	:= os.Getenv("DB_PORT")
	dbName	:= os.Getenv("DB_NAME")
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, dbName)

	// create database connection
	db := createConnection(databaseUrl)

	// create database
	postgres := PostgresSQL{
		db: 		db,
		SocketName: os.Getenv("database_socket_name"),
	}

	// start server
	postgres.serve()

	return &postgres
}
