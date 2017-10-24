CREATE TABLE tempaltes(
    lang     VARCHAR(4) NOT NULL,
    name     VARCHAR(32) NOT NULL,
    value    TEXT NOT NULL,
    PRIMARY KEY(lang, name)
);