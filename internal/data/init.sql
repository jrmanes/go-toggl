CREATE TABLE questions (
                           id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
                           body TEXT NOT NULL,
                           created_at TEXT NOT NULL,
                           updated_at TEXT NOT NULL
);
CREATE TABLE "options" (
                           id INTEGER NOT NULL,
                           body TEXT NOT NULL,
                           correct BOOLEAN NOT NULL,
                           foreign key (id) references questions (id) on delete cascade on update cascade
);