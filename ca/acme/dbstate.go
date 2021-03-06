package acme

import (
	sqlraw "database/sql"
	"fmt"
	"time"

	"github.com/dta4/dns3l-go/state"
)

type ACMEStateManagerSQL struct {
	CAID string
	Prov state.SQLDBProvider
}

type ACMEStateManagerSQLSession struct {
	prov *ACMEStateManagerSQL
	db   *sqlraw.DB
}

func (m *ACMEStateManagerSQL) NewSession() (ACMEStateManagerSession, error) {
	db, err := m.Prov.GetNewDBConn()
	if err != nil {
		return nil, err
	}
	return &ACMEStateManagerSQLSession{db: db, prov: m}, nil
}

func (s *ACMEStateManagerSQLSession) Close() error {
	return s.db.Close()
}

func (s *ACMEStateManagerSQLSession) GetACMEUserPrivkeyByID(userid string) (string, string, error) {

	row := s.db.QueryRow(`select privatekey, registration from `+s.prov.Prov.DBName("acmeusers")+
		` where user_id = $1 AND ca_id = $2 limit 1;`, userid, s.prov.CAID)

	var keyStr string
	var registrationStr string
	err := row.Scan(&keyStr, &registrationStr)
	if err == sqlraw.ErrNoRows {
		return "", "", nil
	} else if err != nil {
		return "", "", err
	}

	return keyStr, registrationStr, nil

}

func (s *ACMEStateManagerSQLSession) PutACMEUser(userid, privatekey,
	registrationStr string, registrationDate time.Time) error {

	_, err := s.db.Exec(`INSERT INTO `+s.prov.Prov.DBName("acmeusers")+
		` (user_id, ca_id, privatekey, registration, registration_date) values ($1, $2, $3, $4, $5);`,
		userid, s.prov.CAID, privatekey, registrationStr, state.TimeToDBStr(registrationDate))

	if err != nil {
		return fmt.Errorf("problem while obtaining certificate: %v", err)
	}
	return nil
}
