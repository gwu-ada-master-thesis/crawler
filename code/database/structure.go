package database



type PostgresSQL struct {
	db 			*pgxpool.Pool
	SocketName 	string
}

