CREATE TABLE quote(
    id INTEGER PRIMARY KEY,
    text  TEXT NOT NULL,
    score  INTEGER NOT NULL ,
    time_created INTEGER NOT NULL, -- seconds since 1970-01-01 00:00:00 UTC
    is_offensive INTEGER NOT NULL, -- a boolean
    is_nishbot INTEGER NOT NULL -- a boolean
);
