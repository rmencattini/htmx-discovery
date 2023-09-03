module htmx

go 1.20

replace htmx/templates => ./templates

require htmx/templates v0.0.0-00010101000000-000000000000

replace htmx/star => ./star

require (
	github.com/google/uuid v1.3.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	htmx/star v0.0.0-00010101000000-000000000000 // indirect
)
