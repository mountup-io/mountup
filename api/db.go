package api

import (
	"database/sql"
	"os"
	"path/filepath"
	"strconv"
)

type DataStore interface {
	InitTables() error
	PutVM(vm *VM) error
	GetHostnameForClientname(clientname string) (string, error)
	PutAuthTokens(tokenDetails TokenDetails) error
	GetAuthTokens() (accessToken string, refreshToken string, accessTokenExp string, refreshTokenExp string, err error)
	DeleteAuthTokens() error
}

type sqlDB struct {
	*sql.DB
}

func NewDB() DataStore {
	homeDir, _ := os.UserHomeDir()
	dbSourcePath := filepath.Join(homeDir, ".mountup/data.db")

	db, _ := sql.Open("sqlite3", dbSourcePath)
	return &sqlDB{db}
}

func (db *sqlDB) InitTables() error {
	stmt, _ := db.Prepare(`
							CREATE TABLE IF NOT EXISTS vm (
							id INTEGER PRIMARY KEY,
							owner_id integer NOT NULL,
							name varchar(32) NOT NULL,
							instance_id varchar(64) NOT NULL,
							public_dns varchar(128) NOT NULL,
							public_ip varchar(64) NOT NULL,
							key_name varchar(64) NOT NULL,
							zone varchar(32) NOT NULL)
							`)
	_, err := stmt.Exec()

	stmt, _ = db.Prepare(`
							CREATE TABLE IF NOT EXISTS auth (
							access_token string NOT NULL,
							refresh_token string NOT NULL,
							access_token_exp string NOT NULL,
							refresh_token_exp string NOT NULL
							)
							`)
	_, err = stmt.Exec()
	return err
}

func (db *sqlDB) PutVM(vm *VM) error {
	stmt, _ := db.Prepare(`
							INSERT INTO vm (
							owner_id, name,
							instance_id, public_dns,
							public_ip, key_name, zone)
							VALUES (?, ?, ?, ?, ?, ?, ?)
							`)

	ownerID, _ := strconv.Atoi(vm.OwnerID)
	_, err := stmt.Exec(ownerID, vm.Name, vm.InstanceID, vm.PublicDNS, vm.PublicIP, vm.KeyName, vm.Zone)
	return err
}

func (db *sqlDB) GetHostnameForClientname(clientname string) (string, error) {
	row := db.QueryRow(`
				SELECT public_dns FROM vm WHERE name=?
				`, clientname)
	var host string
	err := row.Scan(&host)
	if err != nil {
		return "", err
	}

	return host, nil
}

func (db *sqlDB) PutAuthTokens(tokenDetails TokenDetails) error {
	stmt, _ := db.Prepare(`
							INSERT INTO auth (
							access_token, refresh_token,
							access_token_exp, refresh_token_exp)
							VALUES (?, ?, ?, ?)
							`)

	_, err := stmt.Exec(tokenDetails.AccessToken, tokenDetails.RefreshToken,
		tokenDetails.AccessTokenExp.String(), tokenDetails.RefreshTokenExp.String())
	return err
}

func (db *sqlDB) GetAuthTokens() (accessToken string, refreshToken string, accessTokenExp string, refreshTokenExp string, err error) {
	row := db.QueryRow(`
							SELECT * FROM auth
							LIMIT 1
							`)

	err = row.Scan(&accessToken, &refreshToken, &accessTokenExp, &refreshTokenExp)
	if err != nil {
		return "", "", "", "", err
	}
	return
}

func (db *sqlDB) DeleteAuthTokens() error {
	stmt, _ := db.Prepare(`
							DELETE FROM auth
							`)

	_, err := stmt.Exec()
	return err
}
