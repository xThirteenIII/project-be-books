CREATE TABLE reviews (
    id          CHAR(36) PRIMARY KEY,
    book_id     INT NOT NULL,
    score       INT NOT NULL,
    review_text TEXT NOT NULL,
    status      ENUM('pending', 'ready', 'failed') NOT NULL DEFAULT 'pending',
    book_title  VARCHAR(500),
    authors     TEXT,
    cover_url   VARCHAR(500),
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
