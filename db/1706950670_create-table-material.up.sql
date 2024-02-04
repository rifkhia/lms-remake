CREATE TABLE materials (
    id int PRIMARY KEY ,
    title VARCHAR(50) NOT NULL ,
    file VARCHAR(50) NOT NULL ,
    class_section_id int REFERENCES class_sections,
    created_at TIMESTAMP NOT NULL ,
    updated_at TIMESTAMP NOT NULL ,
    deleted_at TIMESTAMP
)