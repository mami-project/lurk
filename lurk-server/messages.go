package main

// JSON response messages
type MsgRegistrationPending struct {
	Status     string `json:"status"`
	RetryAfter uint64 `json:"retry-after"`
}

type MsgRegistrationDone struct {
	Status   string `json:"status"`
	Lifetime uint   `json:"lifetime"`
	CertURL  string `json:"certificates"`
}

type MsgRegistrationFailed struct {
	Status  string `json:"status"`
	Details string `json:"details"`
}

type MsgError struct {
	Details string `json:"error"`
}
