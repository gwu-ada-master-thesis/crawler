package database


import (
	"database/sql"
	"os"
	"strconv"
)

// -------------------------------------------------------------------------------------
// ------------------------------------------------------------------ REQUEST

type RequestArgs struct {
	PublicKey                 string
	PublicKeyByteArr          []byte
	SymmetricKeyCipher        string
	SymmetricKeyCipherByteArr []byte
}

type RequestReply struct {
	Success bool
}

// -------------------------------------------------------------- END INSERT REQUEST
// -------------------------------------------------------------------------------------
// ----------------------------------------------------- SINGLE RECORD QUERY REQUEST

type RequestSingleUserRecordArgs struct {
	PublicKey string
}

type RequestSingleBlockchainRecordArgs struct {
	CurrentBlockHash string
}

type RequestSingleRecordReply struct {
	Row     *sql.Row
	Success bool
}

// ------------------------------------------------- END SINGLE RECORD QUERY REQUEST
// -------------------------------------------------------------------------------------
// ----------------------------------------------------- ALL RECORD QUERY REQUEST

type RequestAllUserFileRecordArgs struct {
	UserId int
}

type RequestAllUserFileRecordReply struct {
	Row     *sql.Rows
	Success bool
}

// ------------------------------------------------- END ALL RECORD QUERY REQUEST
// -------------------------------------------------------------------------------------
// ---------------------------------------------------------- USER EXISTENCE REQUEST

type RequestUserExistenceArgs struct {
	PublicKey string
}

type RequestUserExistenceReply struct {
	Row     *sql.Row
	Success bool
}

// ------------------------------------------------------ END USER EXISTENCE REQUEST
// -------------------------------------------------------------------------------------

// UNIX-domain socket name in /var/tmp, for the database
func databaseSock() string {
	return "/var/tmp/821-db-" + strconv.Itoa(os.Getuid())
}