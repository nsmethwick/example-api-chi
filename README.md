# {Project Name}

This application has everything you need to start building a RESTful API using Go.  

This project includes:
- mysql support
- support for route organization + middleware injection
- Handler Layer <> Service Layer <> Repository Layer Seperation

Copy the project use the example routes to get a quick start on building an API without have to write all of the boilerplate.

## Setup and Local Development

This project structure is designed to have everything you need for local development.  To get started:

```zsh
# install required dependancies via homebrew:
brew install colima docker docker-compose

# in the project directory, install go package dependancies
go mod tidy
```
This project depends on an `.env.local` file for configuration definitions.  Copy the `.env.example` file in the project repository and rename it.  Within that file you can change and of the environment variables to suit your config needs as well as add any new environment variables you might want to use

> [Note]: I might update this to read from a json config in the future if managing env files is too annoying

## Running the Project

Make targets have been defined to make running this application easy.  In the project directory run:
```zsh
make {make-target}
```
### Target: `make run-with-deps`
This make target will start a mysql instance in a docker container as well as run the sql statements defined in `./local/mysql/initdb.d` to seed the database.  This does require `colima` to be running.

### Target: `make run`
This will start the application without starting a mysql instance.  This can be used when you have a mysql instance started that you don't want to overwrite, or if you are connecting to an external database.

## Documentation
Specific documentation for project structure/usage/functionality will be added to the docs folder.  Package overviews will be defined in a README file in each respective folder.

Link: [Docs](./docs/)

## Support
- Author Info
- Contact Info
