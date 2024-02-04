CREATE TABLE class_sections (
    id serial PRIMARY KEY ,
    title VARCHAR(50) NOT NULL ,
    description VARCHAR(255) ,
    "order" int NOT NULL ,
    class_id int REFERENCES classes,
    created_at TIMESTAMP NOT NULL ,
    updated_at TIMESTAMP NOT NULL ,
    deleted_at TIMESTAMP
)