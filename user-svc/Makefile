
migratecreate:
	migrate create -ext sql -dir db/migration -seq $(name)

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/user-svc?sslmode=disable" -verbose down

migrateup: 
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/user-svc?sslmode=disable" -verbose up
