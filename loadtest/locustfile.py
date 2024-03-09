import os
import random

from locust import HttpUser, between, task


API_VERSION = os.getenv("API_VERSION", "1")


class WaiterAPIClient(HttpUser):
    wait_time = between(1, 5)

    @task(5)
    def get_menu_10(self):
        self.client.get(f"/api/v{API_VERSION}/restaurant/10/menu")

    @task(5)
    def get_menu_19(self):
        self.client.get(f"/api/v{API_VERSION}/restaurant/19/menu")

    @task
    def get_menu_rand(self):
        restaurant_id = random.randint(20, 100)
        self.client.get(f"/api/v{API_VERSION}/restaurant/{restaurant_id}/menu")
