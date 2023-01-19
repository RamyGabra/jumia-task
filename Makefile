run:
	go build
	sudo docker-compose -f docker_compose.yml build
	sudo docker-compose -f docker_compose.yml up

stop: 
	sudo docker-compose -f docker_compose.yml down