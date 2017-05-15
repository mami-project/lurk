package starstore

import (
	"database/sql"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func setupTempDb() (db *sql.DB, dbFileName string, err error) {
	tmpDbFile, err := ioutil.TempFile("./", "temp-db")
	if err != nil {
		return
	}

	db, err = DbInit(tmpDbFile.Name())
	if err != nil {
		return
	}

	err = DbCreateRegistrationTable(db)
	if err != nil {
		return
	}

	return db, tmpDbFile.Name(), err
}

func TestInitDB(t *testing.T) {
	db, fname, err := setupTempDb()

	defer os.Remove(fname)

	if err != nil {
		t.Errorf("%s", err)
	}

	if db == nil {
		t.Errorf("Got nil while expecting non-nil DB")
	}
}

func TestDbCreateRegistrationTable(t *testing.T) {
	db, fname, err := setupTempDb()
	if err != nil {
		t.Errorf("%s", err)
	}

	defer os.Remove(fname)

	err = DbCreateRegistrationTable(db)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestDbAddRegistration(t *testing.T) {
	db, fname, err := setupTempDb()
	if err != nil {
		t.Errorf("%s", err)
	}

	defer os.Remove(fname)

	r := Registration{CSR: "a csr", Lifetime: 1234}
	_, err = DbAddRegistration(db, r)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestDbGetRegistrationById(t *testing.T) {
	db, fname, err := setupTempDb()
	if err != nil {
		t.Errorf("%s", err)
	}

	defer os.Remove(fname)

	r := Registration{CSR: "a csr", Lifetime: 1234}
	id, err := DbAddRegistration(db, r)
	if err != nil {
		t.Errorf("%s", err)
	}

	r2, err := DbGetRegistrationById(db, id)
	if err != nil {
		t.Errorf("%s", err)
	}

	if r2.Status != "new" {
		t.Errorf("want: status=new, got status=%s", r2.Status)
	}

	if r2.CertURL != "" {
		t.Errorf("want: certURL=\"\", got certURL=%s", r2.CertURL)
	}

	// Want a reasonably recent timestamp
	delta := r2.CreationDate.Sub(time.Now()) / time.Second
	if delta < -5 {
		t.Errorf("want creation date at most 5s in the past, got: %v", delta)
	}
}

// add a registration
// dequeue a registration
// successfully finalise the dequeued registration
func TestDbUpdateSuccessfulRegistration(t *testing.T) {
	db, fname, err := setupTempDb()
	if err != nil {
		t.Errorf("%s", err)
	}

	defer os.Remove(fname)

	r := Registration{CSR: "a csr", Lifetime: 1234}
	_, err = DbAddRegistration(db, r)
	if err != nil {
		t.Errorf("%s", err)
	}

	reg, err := DbDequeueRegistration(db)
	if err != nil {
		t.Errorf("%s", err)
	}

	ttl := "+48 hours"
	var lifetime uint = 365
	certURL := "http://acme.example.org/path/to/certs"

	err = DbUpdateSuccessfulRegistration(db, reg.Id, certURL, lifetime, ttl)
	if err != nil {
		t.Errorf("%s", err)
	}

	r2, err := DbGetRegistrationById(db, reg.Id)
	if err != nil {
		t.Errorf("%s", err)
	}

	if r2.Status != "done" {
		t.Errorf("want: status done, got %s", r2.Status)
	}

	delta := r2.ExpirationDate.Sub(*r2.CompletionDate)
	if delta != 48*time.Hour {
		t.Errorf("want: delta %s, got %v", ttl, delta)
	}

	if r2.Lifetime != lifetime {
		t.Errorf("want: lifetime %d, got %d", lifetime, r2.Lifetime)
	}

	if r2.CertURL != certURL {
		t.Errorf("want: cert URL %s, got %s", lifetime, r2.Lifetime)
	}
}

func TestDbUpdateFailedRegistration(t *testing.T) {
	db, fname, err := setupTempDb()
	if err != nil {
		t.Errorf("%s", err)
	}

	defer os.Remove(fname)

	r := Registration{CSR: "a csr", Lifetime: 1234}
	_, err = DbAddRegistration(db, r)
	if err != nil {
		t.Errorf("%s", err)
	}

	reg, err := DbDequeueRegistration(db)
	if err != nil {
		t.Errorf("%s", err)
	}

	errMsg := "this and that happened"

	err = DbUpdateFailedRegistration(db, reg.Id, errMsg)
	if err != nil {
		t.Errorf("%s", err)
	}

	r2, err := DbGetRegistrationById(db, reg.Id)
	if err != nil {
		t.Errorf("%s", err)
	}

	if r2.Status != "failed" {
		t.Errorf("want: status failed, got %s", r2.Status)
	}

	if !r2.ErrMsg.Valid {
		t.Errorf("want: errmsg %s, got NULL", errMsg)
	} else if r2.ErrMsg.String != errMsg {
		t.Errorf("want: errmsg %s, got %s", errMsg, r2.ErrMsg.String)
	}
}

func TestDbListRegistrations(t *testing.T) {
	db, fname, err := setupTempDb()
	if err != nil {
		t.Errorf("%s", err)
	}

	defer os.Remove(fname)

	var wanted = []Registration{
		Registration{
			CSR:      "CSR 1",
			Lifetime: 1,
		},
		Registration{
			CSR:      "CSR 2",
			Lifetime: 2,
		},
	}

	for _, r := range wanted {
		_, err = DbAddRegistration(db, r)
		if err != nil {
			t.Errorf("%s", err)
		}
	}

	got, err := DbListRegistrations(db)
	if err != nil {
		t.Errorf("%s", err)
	}

	if len(got) != len(wanted) {
		t.Errorf("wanted %d results, got %d", len(wanted), len(got))
	}
}
