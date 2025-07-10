CREATE TABLE scheduler (
                           id SERIAL PRIMARY KEY,
                           date DATE,
                           title TEXT,
                           comment TEXT,
                           repeat TEXT
);