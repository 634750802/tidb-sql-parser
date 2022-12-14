import { createApp } from 'vue'
import App from './App.vue'
import BootstrapVue3 from 'bootstrap-vue-3'

import 'bootstrap/dist/css/bootstrap-reboot.css'
import 'bootstrap-vue-3/dist/bootstrap-vue-3.css'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-dark-5/dist/css/bootstrap-dark-plugin.min.css'

createApp(App)
    .use(BootstrapVue3)
    .mount('#app')
