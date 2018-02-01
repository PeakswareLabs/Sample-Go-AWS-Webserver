package oauth

import (
	"database/sql"
	"fmt"
)

// Oauth contains AccountID and Token
type Oauth struct {
	AccountID string
	Token     string
}

func createOauth(db *sql.DB, record *Oauth) error {
	_, err := db.Exec("create table IF NOT EXISTS oauth ( AccountId varchar(255) NOT NULL PRIMARY KEY, token varchar(255) NOT NULL)")
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into oauth(AccountId,token) values($1,$2) ON CONFLICT(AccountId) DO UPDATE SET token= EXCLUDED.token", record.AccountID, record.Token)
	if err != nil {
		return err
	}
	return nil
}

func deleteOauth(db *sql.DB, accountID string) error {
	command := fmt.Sprintf(`delete from oauth where accountId='%s'`, accountID)
	_, err := db.Exec(command)
	if err != nil {
		return err
	}
	return nil
}

// RetrieveOauth is used by others to get the token. Note: There should be an easy way to access this
func RetrieveOauth(db *sql.DB, accountID string) ([]byte, error) {
	var token []byte
	query := fmt.Sprintf(`select token from oauth where accountId='%s'`, accountID)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&token)
		if err != nil {
			return nil, err
		}
	}
	return token, nil
}
