package database

const (
	ConnectionException                           = "08000"
	ConnectionDoesNotExist                        = "08003"
	ConnectionFailure                             = "08006"
	SQLClientUnableToEstablishSQLConnection       = "08001"
	SQLServerRejectedEstablishmentOfSQLConnection = "08004"
	TransactionResolutionUnknown                  = "08007"
	ProtocolViolation                             = "08P01"
)

func IsConnectionException(code string) bool {
	switch code {
	case ConnectionException, ConnectionDoesNotExist, ConnectionFailure, SQLClientUnableToEstablishSQLConnection, SQLServerRejectedEstablishmentOfSQLConnection, TransactionResolutionUnknown, ProtocolViolation:
		return true
	}
	return false
}
