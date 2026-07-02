package repos

import (
	"database/sql"
	"errors"
	realtimeforum "real-time-forum"
)

type DB struct {
	Conn *sql.DB
}

func (db *DB) DoesEmailExists(email string) error {
	var existingEmail string
	err := db.Conn.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&existingEmail)
	if err == nil {
		return realtimeforum.ErrEmailExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return realtimeforum.ErrInternal
	}
	return nil
}

func (db *DB) DoesNicknameExists(nickname string) error {
	var existingNick string
	err := db.Conn.QueryRow("SELECT nickName FROM users WHERE nickName = ?", nickname).Scan(&existingNick)
	if err == nil {
		return realtimeforum.ErrNickName
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return realtimeforum.ErrInternal
	}
	return nil
}

func (db *DB) InsertUser(nickName, firstName, lastName, email, hashedPassword string, yearOfBirth int, gender string) (int64, error) {
	result, err := db.Conn.Exec(
		`INSERT INTO users (nickName, firstName, lastName, email, hashedPassword, yearOfBirth, gender)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		nickName, firstName, lastName, email, hashedPassword, yearOfBirth, gender,
	)
	if err != nil {
		println(err.Error())
		return 0, realtimeforum.ErrInternal
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, realtimeforum.ErrInternal
	}

	return userID, nil
}

func (db *DB) GetUserCredentials(identifier string) (int64, string, error) {
	var userID int64
	var hashedPassword string

	err := db.Conn.QueryRow(
		"SELECT userId, hashedPassword FROM users WHERE email = ? OR nickName = ?",
		identifier, identifier,
	).Scan(&userID, &hashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, "", realtimeforum.ErrInvalidCredentials
	}
	if err != nil {
		return 0, "", realtimeforum.ErrInternal
	}

	return userID, hashedPassword, nil
}
