import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'

Vue.use(Vuetify)

// const opts = {}

const vuetify = new Vuetify({
    theme: {
        themes: {
            light: {
                primary: '#3384FF',
                secondary: '#33EAFF',
                triadic: '#1CBCD1',
                accent: '#8c9eff',
                error: '#b71c1c',
                background: '#e0f7fa',
            },
        },
    },
})

export default vuetify