package storage

const (
	createTableQuery = `
    CREATE TABLE IF NOT EXISTS Links (
        ID SERIAL PRIMARY KEY,
        CorrelationID VARCHAR(128) NULL,
        ShortURL VARCHAR(128) NOT NULL,
        OriginalURL VARCHAR(512) NOT NULL UNIQUE,
        AddedDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`

	insertLinkQuery = `
    INSERT INTO Links (ShortURL, OriginalURL)
    VALUES ($1, $2)
    ON CONFLICT (OriginalURL)
    DO NOTHING`

	selectOriginalURLQuery = `
    SELECT OriginalURL FROM Links WHERE ShortURL = $1`
)
