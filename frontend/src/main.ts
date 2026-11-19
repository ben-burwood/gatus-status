import { createApp } from 'vue';
import App from './App.vue';
import './style.css';

export const SERVER_URL = `${window.location.origin}/api`;

createApp(App).mount('#app')
