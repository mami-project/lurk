package starstore

import (
	"database/sql"
	"time"
)

const DefaultLifetime = 365

type Registration struct {
	Id             string
	Status         string
	CSR            string
	CreationDate   time.Time
	CompletionDate *time.Time // may be Null
	ExpirationDate *time.Time // may be Null
	Lifetime       uint
	CertURL        string
	ErrMsg         sql.NullString
}
