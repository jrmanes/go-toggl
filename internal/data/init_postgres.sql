CREATE TABLE IF NOT EXISTS questions (
     id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
     body VARCHAR(256) NOT NULL,
     created_at timestamp DEFAULT now(),
     updated_at timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS options (
     id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
     body VARCHAR(256) NOT NULL,
     correct BOOLEAN NOT NULL,
     created_at timestamp DEFAULT now(),
     updated_at timestamp NOT NULL,
     foreign key (id) references questions (id) on delete cascade on update cascade
);