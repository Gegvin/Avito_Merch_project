package main

import (
	"Avito_Merch_project/config"
	"Avito_Merch_project/internal/handlers"
	"Avito_Merch_project/internal/middleware"
	"Avito_Merch_project/internal/repository"
	"Avito_Merch_project/internal/services"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// передача данных в шаблоны
type PageData struct {
	Username  string
	Error     string
	Coins     int
	Inventory interface{}
	Received  interface{}
	Sent      interface{}
}

var templates = template.Must(template.ParseGlob("./web/*.html"))

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключение к базе данных
	db, err := repository.NewPostgresDB(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Инициализация репозитория и сервисов
	repo := repository.NewRepository(db)
	authService := services.NewAuthService(repo, cfg.JWTSecret)
	coinService := services.NewCoinService(repo)
	merchService := services.NewMerchService(repo)

	handlers.SetAuthService(authService)
	handlers.SetCoinService(coinService)
	handlers.SetMerchService(merchService)
	handlers.SetRepository(repo)

	r := mux.NewRouter()

	// (без JWT проверки)
	apiPublic := r.PathPrefix("/api").Subrouter()
	apiPublic.HandleFunc("/auth", handlers.AuthHandler).Methods("POST")
	apiPublic.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")

	// (с JWT проверкой)
	apiPrivate := r.PathPrefix("/api").Subrouter()
	apiPrivate.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	apiPrivate.HandleFunc("/info", handlers.InfoHandler).Methods("GET")
	apiPrivate.HandleFunc("/sendCoin", handlers.SendCoinHandler).Methods("POST")
	apiPrivate.HandleFunc("/buy/{item}", handlers.BuyHandler).Methods("GET")

	// для веба

	// редирект на /dashboard, если есть токен, иначе на /login
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if token, err := r.Cookie("token"); err != nil || token.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}).Methods("GET")

	// (GET)
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "login.html", nil)
	}).Methods("GET")

	// (POST)
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		token, err := authService.Authenticate(username, password)
		if err != nil {
			templates.ExecuteTemplate(w, "login.html", PageData{Error: "Неверный логин или пароль"})
			return
		}
		// Сохраняем JWT в куки
		cookie := &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}).Methods("POST")

	// (GET)
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "register.html", nil)
	}).Methods("GET")

	//  (POST)
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		token, err := authService.Register(username, password)
		if err != nil {
			templates.ExecuteTemplate(w, "register.html", PageData{Error: "Ошибка регистрации: " + err.Error()})
			return
		}
		cookie := &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}).Methods("POST")

	//dashboard
	r.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		data, ok := getUserDataFromToken(r, cfg.JWTSecret)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		templates.ExecuteTemplate(w, "dashboard.html", data)
	}).Methods("GET")

	// Страница "Мой кошелёк и история"
	r.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		data, ok := getUserDataFromToken(r, cfg.JWTSecret)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		coins, inventory, coinHistory, err := repo.GetUserInfo(data.Username)
		if err != nil {
			http.Error(w, "Ошибка получения информации", http.StatusInternalServerError)
			return
		}
		data.Coins = coins
		data.Inventory = inventory
		data.Received = coinHistory.Received
		data.Sent = coinHistory.Sent
		templates.ExecuteTemplate(w, "info.html", data)
	}).Methods("GET")

	// Страница "Перевести монеты"
	r.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := getUserDataFromToken(r, cfg.JWTSecret); !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		templates.ExecuteTemplate(w, "send.html", nil)
	}).Methods("GET")

	// (POST)
	r.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		data, ok := getUserDataFromToken(r, cfg.JWTSecret)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
			return
		}
		toUser := r.FormValue("toUser")
		amountStr := r.FormValue("amount")
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			templates.ExecuteTemplate(w, "send.html", PageData{Error: "Неверное количество монет"})
			return
		}
		if err := coinService.SendCoins(data.Username, toUser, amount); err != nil {
			templates.ExecuteTemplate(w, "send.html", PageData{Error: "Ошибка перевода: " + err.Error()})
			return
		}
		http.Redirect(w, r, "/info", http.StatusSeeOther)
	}).Methods("POST")

	// Страница "Купить мерч"
	r.HandleFunc("/buy", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := getUserDataFromToken(r, cfg.JWTSecret); !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		templates.ExecuteTemplate(w, "buy.html", nil)
	}).Methods("GET")

	// (POST)
	r.HandleFunc("/buy", func(w http.ResponseWriter, r *http.Request) {
		data, ok := getUserDataFromToken(r, cfg.JWTSecret)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
			return
		}
		item := r.FormValue("item")
		if err := merchService.BuyMerch(data.Username, item); err != nil {
			templates.ExecuteTemplate(w, "buy.html", PageData{Error: "Ошибка покупки: " + err.Error()})
			return
		}
		http.Redirect(w, r, "/info", http.StatusSeeOther)
	}).Methods("POST")

	// удаляет куку и редирект на страницу входа
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:    "token",
			Value:   "",
			Expires: time.Unix(0, 0),
			Path:    "/",
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}).Methods("GET")

	// Отдача статических файлов
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Сервис запущен на порту 8080")
	log.Fatal(srv.ListenAndServe())
}

// извлекает данные пользователя из куки с JWT и возвращает структуру PageData.
func getUserDataFromToken(r *http.Request, secret string) (PageData, bool) {
	cookie, err := r.Cookie("token")
	if err != nil || cookie.Value == "" {
		return PageData{}, false
	}
	tokenStr := cookie.Value
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return PageData{}, false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return PageData{}, false
	}
	username, ok := claims["username"].(string)
	if !ok {
		username = "Неизвестный"
	}
	return PageData{Username: username}, true
}
