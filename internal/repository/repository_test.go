package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateUserIfNotExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	// Ожидаем выполнение запроса
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (username, coins) VALUES ($1, 1000) ON CONFLICT (username) DO NOTHING")).
		WithArgs("testuser").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUserIfNotExists("testuser")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	// Ожидаем запрос
	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)")).
		WithArgs("existinguser").
		WillReturnRows(rows)

	exists, err := repo.UserExists("existinguser")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if !exists {
		t.Errorf("expected user to exist")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTransferCoins(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	mock.ExpectBegin()
	// Ожидаем запрос SELECT coins FROM users
	rows := sqlmock.NewRows([]string{"coins"}).AddRow(100)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT coins FROM users WHERE username = $1 FOR UPDATE")).
		WithArgs("user1").
		WillReturnRows(rows)
	// Обновляем баланс отправителя
	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET coins = coins - $1 WHERE username = $2")).
		WithArgs(50, "user1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	// Добавляем получателя
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (username, coins) VALUES ($1, 1000) ON CONFLICT (username) DO NOTHING")).
		WithArgs("user2").
		WillReturnResult(sqlmock.NewResult(1, 1))
	// Обновляем баланс получателя
	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET coins = coins + $1 WHERE username = $2")).
		WithArgs(50, "user2").
		WillReturnResult(sqlmock.NewResult(1, 1))
	// Вставляем транзакцию
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO coin_transactions (from_user, to_user, amount) VALUES ($1, $2, $3)")).
		WithArgs("user1", "user2", 50).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.TransferCoins("user1", "user2", 50)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPurchaseMerch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	mock.ExpectBegin()
	// Ожидаем запрос
	rows := sqlmock.NewRows([]string{"coins"}).AddRow(200)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT coins FROM users WHERE username = $1 FOR UPDATE")).
		WithArgs("user1").
		WillReturnRows(rows)
	// Обновляем баланс пользователя
	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET coins = coins - $1 WHERE username = $2")).
		WithArgs(80, "user1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	// Обновляем инвентарь вставляем или обновляем запись
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO inventory (username, item, quantity)
		VALUES ($1, $2, 1)
		ON CONFLICT (username, item) DO UPDATE SET quantity = inventory.quantity + 1`)).
		WithArgs("user1", "t-shirt").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.PurchaseMerch("user1", "t-shirt", 80)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
