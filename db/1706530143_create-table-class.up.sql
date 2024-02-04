CREATE TABLE classes(
  id serial primary key ,
  name varchar(50) not null ,
  description varchar(255) not null ,
  key varchar(255) not null ,
  teacher_id varchar(255) references teachers NOT NULL,
  created_at TIMESTAMP NOT NULL ,
  updated_at TIMESTAMP NOT NULL ,
  deleted_at TIMESTAMP
);