build:
	docker build -t rendezvous-server .
run:
	docker run -p 9000:9000 rendezvous-server