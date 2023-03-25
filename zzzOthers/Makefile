# A Makefile is a plain text file that contains a set of rules for building, compiling, and executing code. It is a tool that automates the building of executable files and libraries from source code by executing a series of commands. Makefiles are often used in software development projects to simplify the process of building and managing complex systems. The Makefile consists of a set of targets and dependencies, which are used to specify the commands needed to build a specific output file from a set of source files. It is widely used in Linux and other Unix-based systems to manage software builds and deployments.

tes:
	echo "hello world"

opendb:
	podman exec -it ps15 psql -U root

# Makefile digunakan untuk memudahkan kita menjalankan kode dari CLI dan juga memudahkan bekerja secara team agar team tidak perlu susah memperlajari kode asing yang kita buat
postgres:
	podman run --name ps15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=21204444 -d docker.io/library/postgres:15.2-alpine

createdb:
	podman exec -it ps15 createdb --username=root --owner=root bank_mandiri

dropdb:
	podman exec -it ps15 dropdb bank_mandiri

migup:
	migrate -path db/migration -database "postgresql://root:21204444@localhost:5432/bank_mandiri?sslmode=disable" -verbose up

migdown:
	migrate -path db/migration -database "postgresql://root:21204444@localhost:5432/bank_mandiri?sslmode=disable" -verbose down

sqlc:
	sqlc generate
# jika terdapat lebih dari satu versi schema di migration itu akan membuat error proses generatenya

#  -v membuat log verbose, -cover membuat test coverage, ./... run unit test in all of them
test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go backendmaster/db/sqlc Store

.PHONY: postgres createdb dropdb tes opendb migup sqlc mock