//go:build e2e
// +build e2e

package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var (
	db  *sqlx.DB
	srv *echo.Echo
)

func TestMain(m *testing.M) {
	// 1. Читаем DSN для тестовой базы из env
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		panic("TEST_DATABASE_DSN is not set")
	}

	// 2. Пересоздаём и мигрируем тестовую БД
	exec.Command("make", "migrate-test-up").Run()

	// 3. Подключаемся к БД
	var err error
	db, err = sqlx.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	// 4. Инициализируем Echo‑сервер точно так же, как в cmd/main.go
	srv = cmd.NewServer(db) // или как у вас называется

	// 5. Запускаем все E2E‑тесты
	code := m.Run()

	// 6. Откатываем миграции и выхлопаем
	exec.Command("make", "migrate-test-down").Run()
	os.Exit(code)
}

func TestE2E_CreateAdvert(t *testing.T) {
	Convey("Когда POST /api/adverts с валидным JSON", t, func() {
		payload := map[string]interface{}{
			"name":        "E2E Ad",
			"description": "Тестовое объявление",
			"photos":      []string{"http://img1"},
			"price":       1000,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/adverts", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Прогоняем через вью вашего echo‑сервера
		srv.ServeHTTP(rec, req)

		Convey("Должен вернуть 201 и тело с ID > 0", func() {
			So(rec.Code, ShouldEqual, http.StatusCreated)
			var resp struct {
				ID int `json:"id"`
			}
			So(json.Unmarshal(rec.Body.Bytes(), &resp), ShouldBeNil)
			So(resp.ID, ShouldBeGreaterThan, 0)
		})
	})
}
