import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import Home from './views/Home.vue'
import Create from './views/Create.vue'
import Poll from './views/Poll.vue'
import Results from './views/Results.vue'
import './style.css'
const router = createRouter({ history:createWebHistory(), routes:[{path:'/',component:Home},{path:'/create',component:Create},{path:'/poll/:id',component:Poll},{path:'/poll/:id/results',component:Results},{path:'/polls/:id',component:Poll},{path:'/polls/:id/results',component:Results}] })
createApp(App).use(router).mount('#app')
