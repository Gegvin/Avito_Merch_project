package repository

import (
	"database/sql"
	"errors"

	"Avito_Merch_project/internal/models"

	_ "github.com/lib/pq"
)

var ErrInvalidAmount = errors.New("invalid amount")

type Repository struct {
	db *sql.DB
}

func NewPostgresDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUserIfNotExists(username string) error {
	_, err := r.db.Exec("INSERT INTO users (username, coins) VALUES ($1, 1000) ON CONFLICT (username) DO NOTHING", username)
	return err
}

func (r *Repository) UserExists(username string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	return exists, err
}

func (r *Repository) CreateUser(username string) error {
	_, err := r.db.Exec("INSERT INTO users (username, coins) VALUES ($1, 1000)", username)
	return err
}

func (r *Repository) GetUserInfo(username string) (int, []models.InventoryItem, models.CoinHistory, error) {
	var coins int
	err := r.db.QueryRow("SELECT coins FROM users WHERE username = $1", username).Scan(&coins)
	if err != nil {
		return 0, nil, models.CoinHistory{}, err
	}

	rows, err := r.db.Query("SELECT item, quantity FROM inventory WHERE username = $1", username)
	if err != nil {
		return coins, nil, models.CoinHistory{}, err
	}
	defer rows.Close()

	var inventory []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return coins, nil, models.CoinHistory{}, err
		}
		inventory = append(inventory, item)
	}

	rRows, err := r.db.Query("SELECT from_user, amount FROM coin_transactions WHERE to_user = $1", username)
	if err != nil {
		return coins, inventory, models.CoinHistory{}, err
	}
	defer rRows.Close()

	var received []models.ReceivedCoin
	for rRows.Next() {
		var rc models.ReceivedCoin
		if err := rRows.Scan(&rc.FromUser, &rc.Amount); err != nil {
			return coins, inventory, models.CoinHistory{}, err
		}
		received = append(received, rc)
	}

	sRows, err := r.db.Query("SELECT to_user, amount FROM coin_transactions WHERE from_user = $1", username)
	if err != nil {
		return coins, inventory, models.CoinHistory{}, err
	}
	defer sRows.Close()

	var sent []models.SentCoin
	for sRows.Next() {
		var sc models.SentCoin
		if err := sRows.Scan(&sc.ToUser, &sc.Amount); err != nil {
			return coins, inventory, models.CoinHistory{}, err
		}
		sent = append(sent, sc)
	}

	coinHistory := models.CoinHistory{
		Received: received,
		Sent:     sent,
	}
	return coins, inventory, coinHistory, nil
}

func (r *Repository) TransferCoins(fromUser, toUser string, amount int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var fromCoins int
	err = tx.QueryRow("SELECT coins FROM users WHERE username = $1 FOR UPDATE", fromUser).Scan(&fromCoins)
	if err != nil {
		return err
	}
	if fromCoins < amount {
		return errors.New("недостаточно монет")
	}

	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE username = $2", amount, fromUser)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO users (username, coins) VALUES ($1, 1000) ON CONFLICT (username) DO NOTHING", toUser)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE users SET coins = coins + $1 WHERE username = $2", amount, toUser)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO coin_transactions (from_user, to_user, amount) VALUES ($1, $2, $3)", fromUser, toUser, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) PurchaseMerch(username, item string, price int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var coins int
	err = tx.QueryRow("SELECT coins FROM users WHERE username = $1 FOR UPDATE", username).Scan(&coins)
	if err != nil {
		return err
	}
	if coins < price {
		return errors.New("недостаточно монет для покупки")
	}

	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE username = $2", price, username)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO inventory (username, item, quantity)
		VALUES ($1, $2, 1)
		ON CONFLICT (username, item) DO UPDATE SET quantity = inventory.quantity + 1`, username, item)
	if err != nil {
		return err
	}

	return tx.Commit()
}
