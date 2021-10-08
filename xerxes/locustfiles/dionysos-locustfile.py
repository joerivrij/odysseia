from helpers import generate_random_declensions_word
from locust import HttpUser, task, between


class Dionysos(HttpUser):
    """runs the loadtests for Dionysos api"""
    wait_time = between(0.5, 2.5)

    @task(1)
    def health(self):
        self.client.get("/dionysos/v1/health")

    @task(1)
    def ping(self):
        self.client.get("/dionysos/v1/ping")

    @task(8)
    def ping(self):
        word = generate_random_declensions_word()
        self.client.get(f"/dionysos/v1/checkGrammar?word={word}")
