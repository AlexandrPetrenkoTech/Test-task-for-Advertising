CREATE TABLE adverts (
                         id          SERIAL PRIMARY KEY,
                         name        VARCHAR(200) NOT NULL,
                         description TEXT NOT NULL,
                         price       NUMERIC(12, 2) NOT NULL,
                         created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE photos (
                        id         SERIAL PRIMARY KEY,
                        advert_id  INTEGER NOT NULL REFERENCES adverts(id) ON DELETE CASCADE,
                        url        TEXT NOT NULL,
                        position   INTEGER NOT NULL    -- 1 = main photo, 2,3… = gallery order
);

CREATE INDEX idx_photos_advert_id ON photos(advert_id);
