CREATE TABLE quote_tag(
    tag VARCHAR(32) NOT NULL,
    quote_id INTEGER NOT NULL,
    PRIMARY KEY (tag, quote_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
