CREATE TABLE quote_tag(
    tag TEXT NOT NULL,
    quote_id INTEGER NOT NULL,
    PRIMARY KEY (tag, quote_id)
);
