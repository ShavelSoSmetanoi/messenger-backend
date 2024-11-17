docker-up:
	docker-compose -f deployments/docker-compose.yml up -d
	goose -dir migrations postgres "user=myuser password=mypassword dbname=mydatabase sslmode=disable" up