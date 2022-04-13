CREATE TABLE users
(
    id        serial PRIMARY KEY,
    firstname VARCHAR(100) NOT NULL,
    lastname  VARCHAR(100) NOT NULL,
    email     VARCHAR(100) NOT NULL,
    password  VARCHAR(100) NOT NULL
);

