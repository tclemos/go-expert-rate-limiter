.PHONY: run
run:
	docker compose up -d --build

.PHONY: stop
stop:
	docker compose down --remove-orphans

.PHONY: restart
restart: stop run

