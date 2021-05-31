<template>
  <div>
    <headingImage></headingImage>
    <search-area></search-area>
  </div>
</template>

<script>
import HeadingImage from '@/components/HeadingImage'
import SearchArea from "@/components/SearchArea";

export default {
  name: 'SearchArea',
  data () {
    return {
    }
  },
  components: {
    SearchArea,
    'headingImage': HeadingImage,
    'search-area': SearchArea
  },
  methods: {},
  computed: {
  },
}
</script>