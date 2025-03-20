import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/home', component: HomeView},
		{path: '/conversation', component: HomeView},
		{path: '/conversation/groupSettings', component: HomeView},
	]
})

export default router
