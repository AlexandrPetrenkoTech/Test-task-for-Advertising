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
	// 1. Read DSN for the test database from env
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		panic("TEST_DATABASE_DSN is not set")
	}

	// 2. Recreate and migrate the test DB
	exec.Command("make", "migrate-test-up").Run()

	// 3. Connect to the database
	var err error
	db, err = sqlx.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	// 4. Initialize Echo server exactly as in cmd/main.go
	srv = cmd.NewServer(db) // or however it's named

	// 5. Run all E2E tests
	code := m.Run()

	// 6. Rollback migrations and exit
	exec.Command("make", "migrate-test-down").Run()
	os.Exit(code)
}

func TestE2E_CreateAdvert(t *testing.T) {
	Convey("When POST /api/adverts with valid JSON", t, func() {
		payload := map[string]interface{}{
			"name":        "E2E Ad",
			"description": "Test advert",
			"photos":      []string{"http://img1"},
			"price":       1000,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/adverts", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Pass through the echo server view
		srv.ServeHTTP(rec, req)

		Convey("Should return 201 and body with ID > 0", func() {
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
		// 1) Insert advert and related photos into the DB
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

		// 2) Without fields param (should return only name, price, and main photo)
		Convey("without ?fields should return only id, name, price, and main_photo", func() {
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

		// 3) With fields=true param (should return all fields and all photos)
		Convey("with ?fields=true should return all fields and photos", func() {
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
		// 1) Seed two adverts with different prices
		_, err := db.Exec(`
            INSERT INTO adverts (name, description, price) VALUES
            ('Cheap Ad', 'desc1', 100),
            ('Expensive Ad', 'desc2', 500)
        `)
		So(err, ShouldBeNil)

		// Link one photo to each advert
		// Get their ids in correct order
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

		// 2) Perform GET /api/adverts?page=1&sort=price_desc
		req := httptest.NewRequest(
			http.MethodGet,
			"/api/adverts?page=1&sort=price_desc",
			nil,
		)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)

		Convey("Should return a list of two items sorted by price descending", func() {
			So(rec.Code, ShouldEqual, http.StatusOK)

			// Define summary response structure
			type summary struct {
				ID        int     `json:"id"`
				Name      string  `json:"name"`
				MainPhoto string  `json:"main_photo"`
				Price     float64 `json:"price"`
			}
			var resp []summary
			So(json.Unmarshal(rec.Body.Bytes(), &resp), ShouldBeNil)

			// Check length of response
			So(len(resp), ShouldEqual, 2)

			// First item should be Expensive Ad
			So(resp[0].Name, ShouldEqual, "Expensive Ad")
			So(resp[0].Price, ShouldEqual, 500.0)
			So(resp[0].MainPhoto, ShouldEqual, fmt.Sprintf("http://photo-%d", ids[0]))

			// Second — Cheap Ad
			So(resp[1].Name, ShouldEqual, "Cheap Ad")
			So(resp[1].Price, ShouldEqual, 100.0)
			So(resp[1].MainPhoto, ShouldEqual, fmt.Sprintf("http://photo-%d", ids[1]))
		})
	})
}

func TestE2E_UpdateAdvert(t *testing.T) {
	Convey("E2E: PUT /api/adverts/:id", t, func() {
		// 1) Insert initial advert
		var id int
		So(db.QueryRow(`
            INSERT INTO adverts (name, description, price)
            VALUES ($1, $2, $3)
            RETURNING id
        `, "Old Name", "Old Desc", 150).Scan(&id), ShouldBeNil)

		// And main photo
		So(db.Exec(`
            INSERT INTO photos (advert_id, url, position)
            VALUES ($1, $2, 1)
        `, id, "http://old-photo"), ShouldBeNil)

		// 2) Create update payload
		updatePayload := map[string]interface{}{
			"title":       "New Name",
			"description": "New Desc",
			"photos":      []string{"http://new-main", "http://new-gallery"},
			"price":       300,
		}
		body, _ := json.Marshal(updatePayload)

		// 3) Perform PUT /api/adverts/:id
		url := fmt.Sprintf("/api/adverts/%d", id)
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)

		// 4) Expect 204 No Content
		So(rec.Code, ShouldEqual, http.StatusNoContent)

		// 5) Immediately check updated data with GET detail
		req2 := httptest.NewRequest(http.MethodGet, url+"?fields=true", nil)
		rec2 := httptest.NewRecorder()
		srv.ServeHTTP(rec2, req2)

		So(rec2.Code, ShouldEqual, http.StatusOK)
		var resp struct {
			ID          int      `json:"id"`
			Name        string   `json:"name"`
			Description string   `json:"description"`
			Price       float64  `json:"price"`
			Photos      []string `json:"photos"`
		}
		So(json.Unmarshal(rec2.Body.Bytes(), &resp), ShouldBeNil)

		So(resp.ID, ShouldEqual, id)
		So(resp.Name, ShouldEqual, "New Name")
		So(resp.Description, ShouldEqual, "New Desc")
		So(resp.Price, ShouldEqual, 300.0)
		So(len(resp.Photos), ShouldEqual, 2)
		So(resp.Photos[0], ShouldEqual, "http://new-main")
		So(resp.Photos[1], ShouldEqual, "http://new-gallery")
	})
}

func TestE2E_DeleteAdvert(t *testing.T) {
	Convey("E2E: DELETE /api/adverts/:id", t, func() {
		// 1) Insert test advert and photo
		var id int
		So(db.QueryRow(`
            INSERT INTO adverts (name, description, price)
            VALUES ($1, $2, $3)
            RETURNING id
        `, "ToDelete", "Will be deleted", 123).Scan(&id), ShouldBeNil)

		So(db.Exec(`
            INSERT INTO photos (advert_id, url, position)
            VALUES ($1, $2, 1)
        `, id, "http://tobedeleted"), ShouldBeNil)

		// 2) Perform DELETE /api/adverts/:id
		url := fmt.Sprintf("/api/adverts/%d", id)
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)

		// 3) Expect 204 No Content
		So(rec.Code, ShouldEqual, http.StatusNoContent)
		So(rec.Body.Len(), ShouldEqual, 0)

		// 4) Check that GET by the same ID returns 404
		req2 := httptest.NewRequest(http.MethodGet, url, nil)
		rec2 := httptest.NewRecorder()
		srv.ServeHTTP(rec2, req2)

		So(rec2.Code, ShouldEqual, http.StatusNotFound)
	})
}
