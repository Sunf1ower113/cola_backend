# Auth-API

## Description
This project is a REST API built in Go using clean architecture principles. It provides user registration, authentication, and authorization service.

## Features
- **User Registration:** Allows users to register by providing their email and password.
- **User Authentication:** Supports user login with email and password.
- **User Authorization:** Provides access control based on user roles and permissions.
- **Secure Password Storage:** Utilizes secure hashing algorithms to store user passwords safely.
- **JWT-based Authentication:** Uses JSON Web Tokens (JWT) for user authentication and authorization.
- **Middleware for accecing Protected Routes:** Includes middleware for accecing protected routes.

## Installation
1. Clone the repository: `git clone https://github.com/your/repository.git`
2. Navigate to the project directory: `cd project-directory`
3. Install dependencies: `go mod tidy`
4. Build the project: `go build`
5. Run the executable: `./project-name`

## Usage
1. **Register a New User:** Send a POST request to `/register` endspoint with user details (email and password) in the request body.
2. **Authenticate Uer:** Send a POST request to `/login` endpoint with user credentials (email and password) in the request body. Upon successful authentication, the server will respond with a JWT token.
3. **Access Protected Routes:** Include the JWT token in the Authorization header of subsequent requests to access protected routes.

## Dependencies
- [JWT-Go](https://github.com/dgrijalva/jwt-go): Library for JSON Web Tokens (JWT) in Go.
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt): Package for secure password hashing in Go.
