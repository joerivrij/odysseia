<template>
  <div>
    <sentence-area></sentence-area>
  </div>
</template>

<script>
import SentenceArea from "@/components/SentenceArea";

export default {
  name: 'Herodotos',
  data () {
    return {
    }
  },
  components: {
    SentenceArea,
    'sentence-area': SentenceArea
  },
  methods: {},
  computed: {
  },
}
</script>