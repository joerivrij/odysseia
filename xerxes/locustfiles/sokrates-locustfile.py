from locust import HttpUser, task, between


class Sokrates(HttpUser):
    """runs the loadtests for Sokrates api"""
    wait_time = between(1, 5)

    @task(1)
    def health(self):
        self.client.get("/sokrates/v1/health")

    @task(1)
    def ping(self):
        self.client.get("/sokrates/v1/ping")

    @task(6)
    def create_nomina_question(self):
        with self.client.get("/sokrates/v1/chapters/nomina", catch_response=True) as response:
            chapters = response.json()['lastChapter']
            for i in range(chapters):
                self.client.get(f"/sokrates/v1/createQuestion?category=nomina&chapter={i+1}")

    @task(6)
    def create_verba_question(self):
        with self.client.get("/sokrates/v1/chapters/verba", catch_response=True) as response:
            chapters = response.json()['lastChapter']
            for i in range(chapters):
                self.client.get(f"/sokrates/v1/createQuestion?category=verba&chapter={i+1}")

    @task(6)
    def create_misc_question(self):
        with self.client.get("/sokrates/v1/chapters/misc", catch_response=True) as response:
            chapters = response.json()['lastChapter']
            for i in range(chapters):
                self.client.get(f"/sokrates/v1/createQuestion?category=misc&chapter={i+1}")

    @task(8)
    def check_answer(self):
        body = {"answerProvided": "godin", "quizWord": "θεός", "category": "nomina"}
        self.client.post("/sokrates/v1/answer", json=body)
