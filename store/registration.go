package lurkstore

import "time"

type Registration struct {
	Id             string
	Status         string
	CreationDate   time.Time
	CompletionDate time.Time
	ExpirationDate time.Time
	CSR            string
	Lifetime       uint
	CertURL        string
}
