CREATE TABLE IF NOT EXISTS questions (
        id SERIAL NOT NULL PRIMARY KEY,
        body VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS options (
        id SERIAL NOT NULL,
        body VARCHAR(256) NOT NULL,
        correct BOOLEAN NOT NULL,
        foreign key (id) references questions (id) on delete cascade on update cascade
);