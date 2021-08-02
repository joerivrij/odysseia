from helpers import generate_random_number
from locust import HttpUser, task, between


class Herodotos(HttpUser):
    """runs the loadtests for Herodotos api"""
    wait_time = between(1, 5)

    @task(1)
    def health(self):
        self.client.get("/herodotos/v1/health")

    @task(1)
    def ping(self):
        self.client.get("/herodotos/v1/ping")

    @task(8)
    def check_creation_and_answering_question(self):
        with self.client.get("/herodotos/v1/authors", catch_response=True) as resp:
            random_number = generate_random_number(len(resp.json()['authors']) - 1)
            random_author = resp.json()['authors'][random_number]['author']
            with self.client.get(f"/herodotos/v1/createQuestion?author={random_author}", catch_response=True) as response:
                sentence_id = response.json()['sentenceId']
                body = {
                    "answerSentence": "this is just some answer",
                    "sentenceId": sentence_id,
                    "author": random_author
                }
                self.client.post("/herodotos/v1/checkSentence", json=body)
