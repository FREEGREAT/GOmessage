docker-up:
	cd user && docker-compose up -d
	cd messager && docker-compose up -d
	cd media && docker-compose up -d
docker-down:
	cd user && docker-compose down
	cd messager && docker-compose down
	cd media && docker-compose down