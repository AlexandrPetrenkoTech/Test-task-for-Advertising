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

func TestE2E_GetAdvertByID(t *testing.T) {
	Convey("E2E: GET /api/adverts/:id", t, func() {
		// 1) Вставляем в базу advert и связанные фото
		var id int
		err := db.QueryRow(`
            INSERT INTO adverts (name, description, price)
            VALUES ($1, $2, $3)
            RETURNING id
        `, "Detail Ad", "Detailed description", 500).Scan(&id)
		So(err, ShouldBeNil)

		_, err = db.Exec(`
            INSERT INTO photos (advert_id, url, position) VALUES
            ($1, $2, 1),
            ($1, $3, 2)
        `, id, "http://main-photo", "http://gallery-photo")
		So(err, ShouldBeNil)

		// 2) Без параметра fields (должны вернуть только name, price и главную фотографию)
		Convey("без ?fields должно вернуть только id, name, price и main_photo", func() {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/adverts/%d", id), nil)
			rec := httptest.NewRecorder()
			srv.ServeHTTP(rec, req)

			So(rec.Code, ShouldEqual, http.StatusOK)

			var resp struct {
				ID        int     `json:"id"`
				Name      string  `json:"name"`
				Price     float64 `json:"price"`
				MainPhoto string  `json:"main_photo"`
			}
			So(json.Unmarshal(rec.Body.Bytes(), &resp), ShouldBeNil)
			So(resp.ID, ShouldEqual, id)
			So(resp.Name, ShouldEqual, "Detail Ad")
			So(resp.Price, ShouldEqual, 500.0)
			So(resp.MainPhoto, ShouldEqual, "http://main-photo")
		})

		// 3) С параметром fields=true (должны вернуть все поля + все фото)
		Convey("c ?fields=true должно вернуть все поля и список photos", func() {
			url := fmt.Sprintf("/api/adverts/%d?fields=true", id)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			srv.ServeHTTP(rec, req)

			So(rec.Code, ShouldEqual, http.StatusOK)

			var resp struct {
				ID          int      `json:"id"`
				Name        string   `json:"name"`
				Description string   `json:"description"`
				Price       float64  `json:"price"`
				Photos      []string `json:"photos"`
			}
			So(json.Unmarshal(rec.Body.Bytes(), &resp), ShouldBeNil)
			So(resp.ID, ShouldEqual, id)
			So(resp.Name, ShouldEqual, "Detail Ad")
			So(resp.Description, ShouldEqual, "Detailed description")
			So(resp.Price, ShouldEqual, 500.0)
			So(len(resp.Photos), ShouldEqual, 2)
			So(resp.Photos[0], ShouldEqual, "http://main-photo")
			So(resp.Photos[1], ShouldEqual, "http://gallery-photo")
		})
	})
}

func TestE2E_ListAdverts(t *testing.T) {
	Convey("E2E: GET /api/adverts?page=1&sort=price_desc", t, func() {
		// 1) Засеяем два объявления с разными ценами
		_, err := db.Exec(`
            INSERT INTO adverts (name, description, price) VALUES
            ('Cheap Ad', 'desc1', 100),
            ('Expensive Ad', 'desc2', 500)
        `)
		So(err, ShouldBeNil)

		// Свяжем к каждому объявлению по одному фото
		// Получаем их id в правильном порядке
		rows, err := db.Query(`SELECT id FROM adverts ORDER BY price DESC`)
		So(err, ShouldBeNil)
		defer rows.Close()

		var ids []int
		for rows.Next() {
			var id int
			So(rows.Scan(&id), ShouldBeNil)
			ids = append(ids, id)
		}
		So(len(ids), ShouldEqual, 2)

		for _, id := range ids {
			_, err := db.Exec(`
                INSERT INTO photos (advert_id, url, position) VALUES
                ($1, $2, 1)
            `, id, fmt.Sprintf("http://photo-%d", id))
			So(err, ShouldBeNil)
		}

		// 2) Выполняем GET /api/adverts?page=1&sort=price_desc
		req := httptest.NewRequest(
			http.MethodGet,
			"/api/adverts?page=1&sort=price_desc",
			nil,
		)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)

		Convey("Должен вернуть список из двух элементов в порядке убывания цены", func() {
			So(rec.Code, ShouldEqual, http.StatusOK)

			// Описываем структуру summary-ответа
			type summary struct {
				ID        int     `json:"id"`
				Name      string  `json:"name"`
				MainPhoto string  `json:"main_photo"`
				Price     float64 `json:"price"`
			}
			var resp []summary
			So(json.Unmarshal(rec.Body.Bytes(), &resp), ShouldBeNil)

			// Длина ответа
			So(len(resp), ShouldEqual, 2)

			// Первый элемент — Expensive Ad
			So(resp[0].Name, ShouldEqual, "Expensive Ad")
			So(resp[0].Price, ShouldEqual, 500.0)
			So(resp[0].MainPhoto, ShouldEqual, fmt.Sprintf("http://photo-%d", ids[0]))

			// Второй — Cheap Ad
			So(resp[1].Name, ShouldEqual, "Cheap Ad")
			So(resp[1].Price, ShouldEqual, 100.0)
			So(resp[1].MainPhoto, ShouldEqual, fmt.Sprintf("http://photo-%d", ids[1]))
		})
	})
}

func TestE2E_UpdateAdvert(t *testing.T) {}
func TestE2E_DeleteAdvert(t *testing.T) {}
