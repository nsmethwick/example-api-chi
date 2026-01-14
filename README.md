# API Boilerplate
written in Go using the Chi library

## Summary
This should have everything you need to start a RESTful API using Go.  There are a few features I'd like to add - ref [TODO.md](todo.md)

## Usage
You will need to install a few dependancies before you can run this project:
- [install go](https://go.dev/doc/install)
- install packages via homebrew:
```
brew install colima docker docker-compose
```
- Get package dependancies
```
# in project directory
go mod tidy
```
- Create an `.env.local` or `.env` file by copying and renaming the `.env.example` file
- Run `make run-with-deps` to start mysql and seed the database then start the app
- Run `make run` if mysql is already started and you don't want to overwrite any data