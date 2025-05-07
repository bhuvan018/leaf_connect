module github.com/plantexchange/app

go 1.23.0

toolchain go1.24.2

require (
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/sessions v1.4.0
	github.com/rs/cors v1.11.1
	golang.org/x/crypto v0.37.0
)

require (
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
)

replace github.com/gorilla/sessions => github.com/gorilla/sessions v1.2.1
