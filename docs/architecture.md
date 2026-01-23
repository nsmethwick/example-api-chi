# Architecture
This project uses a Handler Layer <> Service Layer <> Repository Layer structure for organizing and seperating functionality.

## Handler Layer
> Managed in `/internal/routes/`

This layer is responsible for defining route handlers and all things related to HTTP Request and Response.  This includes:
- Serializers
- Request input validation
- Mapping errors to HTTP Status codes
- Checking content headers

## Service Layer
> Managed in `/internal/service/`

This layer is responsible for business logic.  This layer handles responses/errors from the database but does not interact with the database directly.  It is also responsible for passing results back to the Handler layer though it does not worry about formatting/response codes.

## Repository Layer
> Managed in `/internal/repository/`

This layer is responsible for all interactions with our database.  All database queries are defined here for use in the service layer.  This layer has no knowledge of business logic.