<script setup>
import { RouterLink, RouterView } from 'vue-router'
import Search from './components/SearchUsers.vue';
import Group from './components/CreateGroup.vue'
</script>

<script>
export default {
	data() {
		return {
			showSearch: false,
			showGroup: false,

			logged: sessionStorage.token ? true : false,

			userId: sessionStorage.userId,
			userName: sessionStorage.userName,
			userImage: sessionStorage.userImage,
		}
	},
	methods: {
		// Funzione utilizzata per controllare se il file inserito dall'utente è del formato corretto
		handleGroup() {
			this.showGroup = !this.showGroup;
		},
		handleSearch() {
			this.showSearch = !this.showSearch;
		},
		handleLogin() {
			this.logged = true;
			this.userId = sessionStorage.userId;
			this.userName = sessionStorage.userName;
			this.userImage = sessionStorage.userImage;
		}
	}
}
</script>

<template>

	<div class="container-fluid">
		<div class="row">
			<!-- Navigation bar -->
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse" v-if="logged">
				<div class="position-sticky pt-3 sidebar-sticky">

					<!-- Modale utilizzato per la ricerca di un utente con cui aprire una conversazione -->
					<Search :show="showSearch" @close="handleSearch" title="Search">
						<template v-slot:header>
							<h3>Users</h3>
						</template>
					</Search>

					<Group :show="showGroup" @close="handleGroup" title="Search">
						<template v-slot:header>
							<h3>New Group</h3>
						</template>
					</Group>

					<!-- Titolo della NavBar -->
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<img :src="`data:image/jpg;base64,${userImage}`" alt="Profile Picture" class="profile-picture" />
						<span class="username">{{ userName }}</span>
					</h6>

					<!-- Lista dei link della NavBar -->
					<ul class="nav flex-column">

						<li class="nav-item">
							<!-- Link alla pagina Home -->
							<RouterLink to="/home" class="nav-link m-2">
								<!-- Icona Home -->
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#home" />
								</svg>
								Home
							</RouterLink>
						</li>

						<li class="nav-item">
							<!-- Link alla pagina User -->
							<RouterLink to="/userSettings" class="nav-link m-2">
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#settings" />
								</svg>
								User
							</RouterLink>
						</li>

						<!-- Apre il modale (solo se l'utente è loggato) per la ricerca di un utente con cui aprire una conversazione -->
						<li class="nav-item" v-if="logged">
							<a class="nav-link m-2" @click="handleSearch">
								<!-- Icona Search -->
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#user" />
								</svg>
								New Conversation
							</a>
						</li>

						<!-- Apre il modale (solo se l'utente è loggato) per la creazione di un nuovo gruppo -->
						<li class="nav-item" v-if="logged">
							<a class="nav-link m-2" @click="handleGroup">
								<!-- Icona Search -->
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#users" />
								</svg>
								New Group
							</a>
						</li>

					</ul>
					
				</div>
			</nav>
		</div>
		<!-- Contenuto principale della pagina -->
		<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
			<RouterView @successful-login="handleLogin" />
		</main>
	</div>
</template>

<!--Styles for username and profile picture-->
<style>
.profile-picture {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  margin-right: 10px;
  object-fit: cover;
}

.username {
  font-size: 16px;
  font-weight: bold;
}
</style>
