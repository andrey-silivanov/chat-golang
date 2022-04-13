import {createApp} from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import http from './http';
import router from './router'
import auth from './plugins/auth.js';

const app = createApp(App)

app.use(http)
    .use(router)
    .use(auth)
    .mount('#app')