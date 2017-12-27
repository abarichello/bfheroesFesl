run:
	go build -o main cmd/backend/main.go && sudo -H ./main

docker-up:
	sudo docker-compose up

docker-start:
	sudo docker-compose start

docker-stop:
	sudo docker-compose stop

docker-down:
	sudo docker-compose down
