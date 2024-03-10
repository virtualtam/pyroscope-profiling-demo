# Load test - Locust scenarii
COOK_HOST := "http://localhost:8080"
WAITER_HOST := "http://localhost:8081"

locust-cook-v1:
	API_VERSION=1 locust -f loadtest/locustfile.py --host $(COOK_HOST) --web-port 8090
.PHONY: locust-cook-v1

locust-waiter-v1:
	API_VERSION=1 locust -f loadtest/locustfile.py --host $(WAITER_HOST)
.PHONY: locust-waiter-v1

locust-cook-v2:
	API_VERSION=2 locust -f loadtest/locustfile.py --host $(COOK_HOST) --web-port 8090
.PHONY: locust-cook-v2

locust-waiter-v2:
	API_VERSION=2 locust -f loadtest/locustfile.py --host $(WAITER_HOST)
.PHONY: locust-waiter-v2

# Live development server - PostgreSQL console
psql:
	docker compose exec postgres psql -U cook restaurant
.PHONY: psql
