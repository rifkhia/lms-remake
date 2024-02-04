CREATE TABLE students(
    id  VARCHAR(255) PRIMARY KEY ,
    name VARCHAR(255) NOT NULL ,
    nim VARCHAR(255) NOT NULL UNIQUE ,
    email VARCHAR(255) NOT NULL UNIQUE ,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL ,
    updated_at TIMESTAMP NOT NULL ,
    deleted_at TIMESTAMP
);