from helpers import generate_random_word
from locust import HttpUser, task, between


class Alexandros(HttpUser):
    """runs the loadtests for Alexandros api"""
    wait_time = between(0.5, 2.5)

    @task(1)
    def health(self):
        self.client.get("/alexandros/v1/health")

    @task(1)
    def ping(self):
        self.client.get("/alexandros/v1/ping")

    @task(8)
    def ping(self):
        word = generate_random_word()
        self.client.get(f"/alexandros/v1/search?word={word}")
