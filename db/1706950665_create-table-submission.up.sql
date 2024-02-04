CREATE TABLE submissions(
    id serial PRIMARY KEY ,
    title VARCHAR(50) NOT NULL ,
    file VARCHAR(255) ,
    class_section_id int REFERENCES class_sections,
    created_at TIMESTAMP NOT NULL ,
    updated_at TIMESTAMP NOT NULL ,
    deleted_at TIMESTAMP
)