<template>
    <div>
        <div class="topArea">
        <header>
            <h1>Tijd voor de Quiz</h1>
        </header>
        </div>

        <div id="quizarea" class="quizarea">
        <h4> Je hebt {{correctlyAnswered}} van de {{answered}} goed beantwoord ({{percentage}}%)</h4>
        <h3>Je hebt gekozen voor {{categorie}}</h3>

        <button class="categorieButton"
                id="nominaButton"
                @click="categorie = 'nomina'"
                v-on:click="getNextQuestion">
            <strong>Nomina</strong></button>

        <button class="categorieButton"
                id="verbaButton"
                @click="categorie = 'verba'"
                v-on:click="getNextQuestion">
            <strong>Verba</strong></button>

        <button class="categorieButton"
                id="miscButton"
                @click="categorie = 'misc'"
                v-on:click="getNextQuestion">
            <strong>Misc</strong></button>

        <h3>Vertaal: {{quizWord}}</h3><br>
        <div v-for="answer in answers" :key="answer">
            <button class="answerButton"
                    @click="selectedAnswer = answer"
                    v-on:click="postAnswer">{{answer}}</button>
        </div>

        <div v-if="showAnswer">
            <h2 v-if="correct" class="rightText">Correct</h2>
            <h2 v-else class="wrongText">Verkeerd</h2>
            <h3 v-if="!correct"> Het juiste antwoord was: {{correctAnswer}} </h3>

            <button class="categorieButton"
                    id="nextButton"
                    v-on:click="getNextQuestion">
                <strong>Volgende</strong></button>
        </div>
    </div>
    </div>
</template>

<script>
    export default {
        name: "Quiz",
        data() {
            return {
                showAnswer: false,
                correctAnswer: "",
                quizWord: [],
                answers: [],
                selectedAnswer : "",
                correct: false,
                categorie: "nomina",
                answered: 0,
                correctlyAnswered: 0,
                percentage: 100,
            }
        },
        created() {
        },
        methods: {
            getQuestion: function () {
                this.$apiClient.get(`${this.$proxyUrl}/${this.categorie}`)
                    .then((response) => {
                        let shuffeledArray = response.data.slice(1, 5)
                        for (let i = shuffeledArray.length - 1; i > 0; i--) {
                            const j = Math.floor(Math.random() * (i + 1));
                            [shuffeledArray[i], shuffeledArray[j]] = [shuffeledArray[j], shuffeledArray[i]];
                        }
                        console.log(shuffeledArray)
                        this.answers = shuffeledArray;
                        this.quizWord = response.data[0];
                        this.correctAnswer = response.data[1]
                    })
                    .catch(e => {
                        this.errors.push(e)
                    })
            },
            getNextQuestion: function () {
                this.showAnswer = false;
                this.getQuestion()
            },
            postAnswer: function () {
                this.showAnswer = true;
                this.answered++
                this.$apiClient({
                    method: 'post',
                    url: `${this.$proxyUrl}/answer`,
                    data: {
                        "answer" : this.selectedAnswer,
                        "quizWord": this.quizWord,
                        "categorie": this.categorie
                    }
                })
                    .then((response) => {
                        this.correct = response.data['correctAnswer']
                        if (this.correct) {
                            this.correctlyAnswered++
                        }

                        this.percentage = Math.round(this.correctlyAnswered / this.answered * 100)
                    })
                    .catch(e => {
                        this.errors.push(e)
                    })
            }
        },
        mounted() {
            this.getQuestion()
            this.correctlyAnswered = 0
            this.answered = 0
            this.percentage = 100
        }

    }
</script>

<style scoped>
    .quizarea {
        background: lightblue;
        text-align: center;
        margin-left: 25em;
        /*margin-right: 20em;*/
        width: 50%
    }

    .categorieButton {
        background-color: #f96332;
        color: #fff;
        border-width: 2px;
        font-weight: 400;
        font-size: .8571em;
        line-height: 1.35em;
        margin: 10px 1px;
        border: none;
        border-radius: .1875rem;
        padding: 11px 22px;
        cursor: pointer;
        text-align: center;
        vertical-align: middle;
    }

    .categorieButton:hover {
        background: rgba(249, 99, 49, 0.71);
    }

    .answerButton {
        background-color: cadetblue;
        color: #fff;
        border-width: 2px;
        font-weight: 400;
        font-size: .8571em;
        line-height: 1.35em;
        margin: 10px 1px;
        border: none;
        border-radius: .1875rem;
        padding: 11px 22px;
        cursor: pointer;
        text-align: center;
        vertical-align: middle;
        width: 40em;
    }

    .answerButton:hover {
        background: darkseagreen;
    }

    .wrongText {
        color: red
    }

    .rightText {
        color: green;
    }

    ol {
        padding: 0;
        text-align: center;
        list-style-type: decimal;
        margin-left: 15em;
        width: 200px;
    }


    .topArea {
        padding: 8em 8em 8em 8em;
        text-align: center;
        position: relative;
        height: 5em;
    }
    .topArea::after {
        background-image: url("../assets/images/parthenon.jpg");
        opacity: 0.3;
        position: absolute;
        content: "";
        top: 0;
        left: 25em;
        width: 50%;
        height: 50em;
        z-index: -1;
    }

</style>