<template>
    <div>
        <div class="topArea">
            <div class="sidenav">
                <h2>Beschikbare hoofdstukken</h2>
                <button class="categorieButton"
                        @click="chosenChapter = ''">
                    <strong>Alle hoofdstukken</strong></button>
                <div v-for="chapter in chapters" :key="chapter">
                    <a @click="chosenChapter = chapter" v-on:click="getNextQuestion">
                {{chapter}}</a>
                </div>
            </div>
        <header>
            <h1>Tijd voor de Quiz</h1>
            <h2 v-if="chosenChapter">Hoofdstuk {{chosenChapter}}</h2>
        </header>
        </div>

        <div id="quizarea" class="quizarea">
        <h4> Je hebt {{correctlyAnswered}} van de {{answered}} goed beantwoord ({{percentage}}%)</h4>
            <button class="categorieButton"
                    v-on:click="resetProgress">
                <strong>Reset voortgang</strong></button>
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
        <div v-if="showAnswer">
            <h2 v-if="correct" class="rightText">Correct</h2>
            <h2 v-else class="wrongText">Verkeerd</h2>
            <h3 v-if="!correct"> Het juiste antwoord was: {{correctAnswer}} </h3>

            <button class="categorieButton"
                    id="nextButton"
                    v-on:click="getNextQuestion">
                <strong>Volgende</strong></button>
        </div>
        <div v-for="answer in answers" :key="answer">
            <button class="answerButton"
                    @click="selectedAnswer = answer"
                    v-on:click="postAnswer">{{answer}}</button>
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
                chosenChapter : "",
                chapters : []
            }
        },
        created() {
            this.getChapters()
        },
        methods: {
            getQuestion: function () {
                let url = `${this.$proxyUrl}/${this.categorie}`
                if (this.chosenChapter !== "") {
                    url = url + `/${this.chosenChapter}`
                }
                this.$apiClient.get(url)
                    .then((response) => {
                        let shuffeledArray = response.data.slice(1, 5)
                        for (let i = shuffeledArray.length - 1; i > 0; i--) {
                            const j = Math.floor(Math.random() * (i + 1));
                            [shuffeledArray[i], shuffeledArray[j]] = [shuffeledArray[j], shuffeledArray[i]];
                        }
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
            },
            getChapters : function () {
                let url = `${this.$proxyUrl}/chapters`
                this.$apiClient.get(url)
                    .then((response) => {
                        let endOfArray = response.data['chapters']
                        for (let i = 1; i <= endOfArray; i++) {
                            this.chapters.push(i)
                        }
                    })
                    .catch(e => {
                        this.errors.push(e)
                    })
            },
            resetProgress : function () {
                this.correctlyAnswered = 0
                this.answered = 0
                this.percentage = 100
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
        margin: 0.1em;
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
        height: 1em;
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

    /* The sidebar menu */
    .sidenav {
        height: 100%; /* Full-height: remove this if you want "auto" height */
        width: 160px; /* Set the width of the sidebar */
        position: fixed; /* Fixed Sidebar (stay in place on scroll) */
        z-index: 1; /* Stay on top */
        top: 0; /* Stay at the top */
        left: 0;
        background-color: cadetblue;
        overflow-x: hidden; /* Disable horizontal scroll */
        padding-top: 20px;
    }

    /* The navigation menu links */
    .sidenav a {
        padding: 6px 8px 6px 16px;
        text-decoration: none;
        font-size: 25px;
        color: #ffffff;
        display: block;
        cursor: pointer;
    }

    /* When you mouse over the navigation links, change their color */
    .sidenav a:hover {
        color: darkseagreen;
    }

    /* Style page content */
    .main {
        margin-left: 160px; /* Same as the width of the sidebar */
        padding: 0px 10px;
    }

    /* On smaller screens, where height is less than 450px, change the style of the sidebar (less padding and a smaller font size) */
    @media screen and (max-height: 450px) {
        .sidenav {padding-top: 15px;}
        .sidenav a {font-size: 18px;}
    }

</style>