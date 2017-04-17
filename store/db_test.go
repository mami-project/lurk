package lurkstore

import (
	"database/sql"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func initTempDb() (*sql.DB, string, error) {
	tmpDbFile, err := ioutil.TempFile("./", "temp-db")
	if err != nil {
		return nil, "", err
	}

	db, err := DbInit(tmpDbFile.Name())

	return db, tmpDbFile.Name(), err
}

func TestInitDB(t *testing.T) {
	db, fname, err := initTempDb()

	defer os.Remove(fname)

	if err != nil {
		t.Errorf("%s", err)
	}

	if db == nil {
		t.Errorf("Got nil while expecting non-nil DB")
	}
}

func TestDbCreateRegistrationTable(t *testing.T) {
	db, fname, err := initTempDb()
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
	db, fname, err := initTempDb()
	if err != nil {
		t.Errorf("%s", err)
	}

	defer os.Remove(fname)

	err = DbCreateRegistrationTable(db)
	if err != nil {
		t.Errorf("%s", err)
	}

	_, err = DbAddRegistration(db, "a csr", 1234)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestDbGetRegistrationById(t *testing.T) {
	db, fname, err := initTempDb()
	if err != nil {
		t.Errorf("%s", err)
	}

	defer os.Remove(fname)

	err = DbCreateRegistrationTable(db)
	if err != nil {
		t.Errorf("%s", err)
	}

	id, err := DbAddRegistration(db, "a csr", 1234)
	if err != nil {
		t.Errorf("%s", err)
	}

	reg, err := DbGetRegistrationById(db, id)
	if err != nil {
		t.Errorf("%s", err)
	}

	if reg.Status != "new" {
		t.Errorf("want: status=new, got status=%s", reg.Status)
	}

	if reg.CertURL != "" {
		t.Errorf("want: certURL=\"\", got certURL=%s", reg.CertURL)
	}

	// Want a reasonably recent timestamp
	delta := reg.CreationDate.Sub(time.Now()) / time.Second
	if delta < -5 {
		t.Errorf("want creation date at most 5s in the past, got: %v", delta)
	}
}
