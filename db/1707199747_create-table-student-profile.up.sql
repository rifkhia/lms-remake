CREATE TABLE student_profile(
    id varchar primary key references students,
    dateofbirth date not null ,
    gender varchar(1) not null ,
    address varchar not null,
    phone varchar(15) not null ,
    created_at timestamp not null ,
    updated_at timestamp not null ,
    deleted_at timestamp
)