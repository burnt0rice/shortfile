import * as Vue from 'vue'
import './style.css'
import App from './App.vue'
import axios from 'axios'
import VueAxios from 'vue-axios'

const app = Vue.createApp(App)
const client = axios.create({
    baseURL: "/api",
});

app.use(VueAxios, client);
app.mount('#app')
