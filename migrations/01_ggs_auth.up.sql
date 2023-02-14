CREATE TABLE users(
    id serial primary key,
    email varchar(35),
    password_hash varchar(255)
);

CREATE TABLE refresh_tokens(
    refresh_token varchar(255),
    access_code varchar(255)
);