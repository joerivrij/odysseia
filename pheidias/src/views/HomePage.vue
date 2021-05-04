<template>
    <div>
        <headingImage></headingImage>
        <navbar v-on:update-chapter="setChapter"></navbar>
        <quiz-area v-on:update-categorie="setCategorie"></quiz-area>
    </div>
</template>

<script>
    import HeadingImage from '@/components/HeadingImage'
    import Navbar from "@/components/Navbar.vue";
    import QuizArea from "@/components/QuizArea";

    export default {
        name: 'homepage',
        data () {
          return {
            selectedChapter: 1,
            selectedCategorie: 'nomina'
          }
      },
        components: {
          QuizArea,
            'headingImage': HeadingImage,
            'navbar': Navbar,
            'quiz-area': QuizArea
        },
        methods: {
          setChapter(e) {
            this.selectedChapter = e;
          },
          setCategorie(e) {
            this.selectedCategorie = e;
          }
        },
        computed: {
        },
      }
</script>