package mdb

import (
	"database/sql"
	"log"
	"time"

	"github.com/mattn/go-sqlite3"
)

type EmailEntry struct {
	Id          int64
	Email       string
	ConfirmedAt *time.Time
	OptOut      bool
}

// Não sei tratar bem esses erros ainda :(
func TryCreate(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE emails (
			id				INTEGER PRIMARY KEY,
			email			TEXT UNIQUE,
			confirmed_at	INTEGER,
			opt_out			INTEGER
		)
	`)
	if err != nil {
		if sqlError, ok := err.(sqlite3.Error); ok {
			// code == 1 <-> db já existe
			if sqlError.Code != 1 {
				log.Println(sqlError)

			}
			log.Println(err)
		}
	}
}

func emailEntryFromRow(row *sql.Row) (*EmailEntry, error) {
	var id int64
	var email string
	var confirmedAt int64
	var optOut bool

	// Captura dados de uma linha
	err := row.Scan(&id, &email, &confirmedAt, &optOut)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// timestamp Unix -> time.Time Go
	t := time.Unix(confirmedAt, 0)
	return &EmailEntry{
		Id:          id,
		Email:       email,
		ConfirmedAt: &t,
		OptOut:      optOut,
	}, nil
}

func CreateEmail(db *sql.DB, email string) error {
	_, err := db.Exec(`
		INSERT INTO emails(email, confirmed_at, opt_out)
		VALUES(?, 0, false)
	`, email)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetEmail(db *sql.DB, email string) (*EmailEntry, error) {
	row := db.QueryRow(`
		SELECT id, email, confirmed_at, opt_out
		FROM emails
		WHERE email = ?
	`, email)

	emailEntry, err := emailEntryFromRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return nil, err
		}
		return nil, err
	}
	return emailEntry, nil
}

func UpdateEmail(db *sql.DB, entry *EmailEntry) error {
	unixTime := entry.ConfirmedAt.Unix()

	_, err := db.Exec(`
	INSERT INTO
		emails(emails, confirmed_at, opt_ot)
			VALUES(?, ?, ?)
		ON CONFLIT(email) DO UPDATE SET
			confirmed_at=?
			opt_out=?
	`, entry.Email, unixTime, entry.OptOut, unixTime, entry.OptOut)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DeleteEmail(db *sql.DB, email string) error {
	_, err := db.Exec(`
		UPDATE emails
			SET opt_out=true
		WHERE email=?
	`, email)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

type GetEmailBathQueryParams struct {
	Page  int
	Count int
}

func GetEmailBatch(db *sql.DB, params GetEmailBathQueryParams) ([]EmailEntry, error) {
	var empty []EmailEntry

	row, err := db.Query(`
		SELECT id, email, confirmed_at, opt_out
		FROM emails
		WHERE opt_out = false
		ORDER BY id ASC
		LIMIT ? OFFSET ?`, params.Count, (params.Page-1)*params.Count)

	if err != nil {
		log.Println(err)
		return empty, err
	}
	defer row.Close()

	emails := make([]EmailEntry, 0, params.Count)

	for row.Next() {
		emailRow := new(EmailEntry)
		err := row.Scan(&emailRow.Id, &emailRow.Email, &emailRow.ConfirmedAt, &emailRow.OptOut)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		emails = append(emails, *emailRow)
	}

	return emails, nil
}
