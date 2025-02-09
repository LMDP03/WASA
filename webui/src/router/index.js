import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import GroupView from '../views/GroupView.vue'
import ConversationView from '../views/ConversationView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/home', component: HomeView},
		{path: '/conversation', component: ConversationView},
		{path: '/conversation/:convId', component: GroupView},
	]
})

export default router
