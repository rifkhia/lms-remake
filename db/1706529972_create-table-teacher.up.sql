CREATE TABLE teachers(
    id varchar(255) primary key,
    name varchar(50) not null ,
    npm int not null UNIQUE ,
    email varchar(255) not null UNIQUE ,
    password varchar(255) not null,
    created_at TIMESTAMP NOT NULL ,
    updated_at TIMESTAMP NOT NULL ,
    deleted_at TIMESTAMP
);