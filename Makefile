server:		
			go run server.go
run: 		stop up
up:
			docker compose up -d
stop:
			docker compose stop
down:
			docker compose down
build: 
			docker compose build
rebuild: 	stop remove build up
remove: 
			docker compose rm