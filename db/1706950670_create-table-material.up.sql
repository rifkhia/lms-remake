CREATE TABLE materials(
    id serial primary key ,
    title varchar not null ,
    description varchar ,
    file varchar ,
    class_section_id int references class_sections,
    created_at timestamp not null ,
    updated_at timestamp not null ,
    deleted_at timestamp
)