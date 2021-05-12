<template>
  <div>
    <headingImage></headingImage>
    <sentence-area></sentence-area>
  </div>
</template>

<script>
import HeadingImage from '@/components/HeadingImage'
import SentenceArea from "@/components/SentenceArea";

export default {
  name: 'herodotos',
  data () {
    return {
    }
  },
  components: {
    SentenceArea,
    'headingImage': HeadingImage,
    'sentence-area': SentenceArea
  },
  methods: {},
  computed: {
  },
}
</script>