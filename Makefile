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

# Network proxy to simulate latency/jitter
TOXIPROXY_CLI := docker compose exec toxiproxy /toxiproxy-cli
TOXIPROXY_COOK_PROXY := "cook_1"
TOXIPROXY_COOK_TOXIC := cook
TOXIPROXY_POSTGRES_PROXY := "postgres_1"
TOXIPROXY_POSTGRES_TOXIC := pg

toxiproxy-list:
	$(TOXIPROXY_CLI) list
.PHONY: toxiproxy-list

toxiproxy-inspect:
	$(TOXIPROXY_CLI) inspect $(TOXIPROXY_POSTGRES_PROXY)
.PHONY: toxiproxy-inspect

toxiproxy-add:
	$(TOXIPROXY_CLI) toxic add -n $(TOXIPROXY_COOK_TOXIC) -t latency -a latency=10 -a jitter=2 $(TOXIPROXY_COOK_PROXY)
	$(TOXIPROXY_CLI) toxic add -n $(TOXIPROXY_POSTGRES_TOXIC) -t latency -a latency=10 -a jitter=2 $(TOXIPROXY_POSTGRES_PROXY)
.PHONY: toxiproxy-add

toxiproxy-delete:
	$(TOXIPROXY_CLI) toxic delete -n $(TOXIPROXY_COOK_TOXIC) $(TOXIPROXY_COOK_PROXY)
	$(TOXIPROXY_CLI) toxic delete -n $(TOXIPROXY_POSTGRES_TOXIC) $(TOXIPROXY_POSTGRES_PROXY)
.PHONY: toxiproxy-delete

toxiproxy-fiber:
	$(TOXIPROXY_CLI) toxic update -n $(TOXIPROXY_COOK_TOXIC) -a latency=10 -a jitter=2 $(TOXIPROXY_COOK_PROXY)
	$(TOXIPROXY_CLI) toxic update -n $(TOXIPROXY_POSTGRES_TOXIC) -a latency=10 -a jitter=2 $(TOXIPROXY_POSTGRES_PROXY)
.PHONY: toxiproxy-fiber

toxiproxy-dsl:
	$(TOXIPROXY_CLI) toxic update -n $(TOXIPROXY_COOK_TOXIC) -a latency=25 -a jitter=10 $(TOXIPROXY_COOK_PROXY)
	$(TOXIPROXY_CLI) toxic update -n $(TOXIPROXY_POSTGRES_TOXIC) -a latency=25 -a jitter=10 $(TOXIPROXY_POSTGRES_PROXY)
.PHONY: toxiproxy-dsl
