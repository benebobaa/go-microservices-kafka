.PHONY: sqlc migratecreate migratedown migrateup

createdb:
	sudo docker exec -it postgres createdb --username=root --owner=root orchestra_svc

dropdb:
	sudo docker exec -it postgres dropdb orchestra_svc

sqlc:
	sqlc generate

mockdb:
	mockgen -source=internal/repository/sqlc/store.go -package mockdb -destination=internal/repository/mock/store.go


migratecreate:
	migrate create -ext sql -dir db/migration -seq $(name)

migratedown:
	migrate -path db/migration -database "postgresql://root:root@127.0.0.1:5432/orchestra_svc?sslmode=disable" -verbose down

migrateup: 
	migrate -path db/migration -database "postgresql://root:root@127.0.0.1:5432/orchestra_svc?sslmode=disable" -verbose up
