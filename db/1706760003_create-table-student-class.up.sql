CREATE TABLE student_class(
    id serial primary key ,
    student_id varchar references students ,
    class_id int references classes,
    created_at TIMESTAMP NOT NULL ,
    deleted_at TIMESTAMP NOT NULL
)