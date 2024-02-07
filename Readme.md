# LMS-Remake

---

This project is a remake of my old [project](https://github.com/rifkeh/alterra-mini-project). The aim of
this project is to develop a LMS or Learning Management System that includes the CMS or Content Management System.  

The project itself develop using Golang and Fiber framework. For the database, i'm using PostgreSql as database driver
and sqlx as library in Golang. For the storage itself using Supabase which to store some file.  

For documentation of this API, you can check in this [link](https://documenter.getpostman.com/view/32763424/2s9Yyy7xdZ)

## Migrate Database

---

Before running the server, we have to migrate the database first.  

To migrate database, we will use command that already added in Makefile, here is list of the command in Makefile:

```
make create-migration name={YOUR_FILE_NAME}
```
First we have the ``create-migration`` command, which will make two file according to name that you assign in the command.
The file will most likely have the same name excep for the suffix, one is ``up`` and the other is ``down``. 

The purpose of the ``up`` file is to create a migration that most likely to create and insert. For the ``down``
file is to rollback any creation or insertion that we previously made.  

Example
```
make create-migration name=create-student-table
```

Next,

```
make up-migration DATABASE_URL={YOUR_DATABASE_URL}
```

The ``up-migration`` command use to apply all of the sql file that having suffix `up` into our database. Or in other word,
this command will execute the creation or insertion of all .sql file.

Example
```
make up-migration DATABASE_URL=postgres://postgres:postgres@localhost:5432/lms_remake
```

Last,
```
make rollback-migration DATABASE_URL={YOUR_DATABASE_URL}
```

The `rollback-migration` command use to execute all of the sql file that having suffix `down` into our database. Or in other word,
this command will rollback all of the creation or insertion we made.

## Running The Server

---

To start running server, the first step will be installing Go.  

Once you install Go, open the terminal and run following command:
```
go mod download
```

Which will install all of the module that we need in order to run the server. After the progress complete, run this
following command:
```
go build -o main /cmd/
```

This command will compile this source code into executable binary. After its complete, there will be a file named `main`.
That is the executable binary that we will run, to run that use command:
```
./main
```

## Next Feature

---

These are my next plan on improving this project :
* Adding OTP by email or by phone number on registering
* Adding Teacher feedback on student submissions
* Adding exam using Websocket to get handshake from all student on the exam
