package lurkstore

import "time"

type Registration struct {
	Id             string
	Status         string
	CSR            string
	CreationDate   time.Time
	CompletionDate *time.Time // may be Null
	ExpirationDate *time.Time // may be Null
	Lifetime       uint
	CertURL        string
}
