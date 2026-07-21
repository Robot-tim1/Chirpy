## Chirpy
A hands-on Go backend learning project for building a simple social API with authentication, persistence, and HTTP routing. Courtesy of [boot.dev](https://www.boot.dev).

The guided project followed was mostly hands off with implementation details, at most providing pseudo code for what to do. To see what I learned from this project, scroll down past the installation section.

The api docs will be in the /docs folder(currently not made).

## Installation
First clone the project and go into the created folder.
```
git clone https://github.com/Robot-tim1/Chirpy.git
cd Chirpy
```
To set up the server you'll need a PostgreSQL database for it on the machine, which you'll need to migrate the schema too, which will be done with Goose. 

Install Goose.
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```
if you run into problems with goose, try exporting your Go binary path.
```
export PATH="$PATH:$(go env GOPATH)/bin"
```
Once PostgreSQL is installed, start the service.
```
sudo service postgresql start
```
Set your password if needed depending on your os, enter the psql shell. On linux it's.
```
sudo -u postgres psql
```
Then create and setup the database.
```
CREATE DATABASE chirpy;
\c chirpy
ALTER USER postgres PASSWORD 'postgres';
```
You can make the password whatever you want, but I made it postgres for simplicity.

Now to run the migrations. You'll need your connection string, which can be different depending on os. For Linux, this will typically look like:
```
postgres://postgres:postgres@localhost:5432/chirpy
```
The format is 'protocol://username:password@host:port/database'. 

In the root directory run the migrations with your connection string.
```
goose -dir ./sql/schema postgres postgres://postgres:postgres@localhost:5432/chirpy up
```
Now you will need to create a .env file in the root directory and it should look something like this.
```
DB_URL="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
SECRET="replace-with-a-random-string"
POLKA_KEY="f271c81ff7084ee5b99a5091b42d486e"
```
The DB_URL format is the connection string with ?sslmode=disable appended. The POLKA_KEY is used for the mock webhook endpoint and can be left unchanged. For the SECRET value, generate a random string in the terminal with:
```
openssl rand -base64 64
```
Build and run the project
```
go build ./...
```
Or, if you prefer not to build first:
```
go run .
```
