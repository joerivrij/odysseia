<template>
    <div>
        <quiz></quiz>
    </div>
</template>

<script>
    import Quiz from '../components/Quiz.vue'

    export default {
        name: 'homepage',
        data () {
            return {}
        },
        components: {
            'quiz': Quiz,
        },
        computed: {
        },
    }
</script>