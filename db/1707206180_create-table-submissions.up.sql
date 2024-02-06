CREATE TABLE submissions(
    id serial primary key ,
    title varchar not null ,
    description varchar ,
    file varchar ,
    class_section_id int references class_sections
)