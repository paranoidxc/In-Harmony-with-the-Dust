import {createApp} from 'vue'
import App from './App.vue'
import './style.css';

import VueDiff from 'vue-diff'
import 'vue-diff/dist/index.css'

const app = createApp(App);
app.use(VueDiff);
app.mount('#app')
//createApp(App).mount('#app')
