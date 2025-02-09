<script setup>
import { RouterLink, RouterView } from 'vue-router'
import Search from './components/SearchUsers.vue';
import SearchUsers from './components/SearchUsers.vue';
</script>
<script>
export default {
	data() {
		return {
			showSearch: false,

			logged: sessionStorage.token > 0 ? true : false,

			userId: sessionStorage.userId,
			userName: sessionStorage.userName,
			userImage: sessionStorage.userImage,

			showNameUpdate: false,
			newName: "",

			showImageUpdate: false,
			newImage: null,

		}
	},
	methods: {
		// Funzione utilizzata per controllare se il file inserito dall'utente è del formato corretto
		async checkFile(event) {
			this.errorMsg = "";
			const file = event.target.files[0]; // Prende il file inserito dall'utente
			if (file.type != "image/jpeg") {
				this.errorMsg = "Unsupported image type: only jpg or jpeg images allowed.";
				return
			}
			if (file.size > 5242880) {
				this.errorMsg = "File size exceeded: MAX size 5MB."
				return
			}
			this.newImage = file;
		},
		handleImageUpdate() {
			sessionStorage.userImage = this.userImage;
			this.showImageUpdate = !this.showImageUpdate;
			this.newImage = "";
			this.errorMsg = "";
		},
		async setMyPhoto() {
			this.errorMsg = "";

			const formData = new FormData();
			formData.append('image', this.newImage);
			this.$axios.put(`/users/${this.userId}/image`, formData, {headers: { 'Authorization': `${sessionStorage.token}`}}).then(response => {
				this.userImage = response.data.Image;
				this.handleImageUpdate();
			}).catch(e => {
				this.errorMsg = e.toString();
			});
		},
		async handleNameUpdate() {
			sessionStorage.userName = this.userName;
			this.showNameUpdate = !this.showNameUpdate;
			this.newName = "";
			this.errorMsg = "";
		},
		async setMyUserName() {
			this.errorMsg = "";

			if (this.newName == this.userName) {
				this.errorMsg = "Choose a new name different from the current one.";
				return
			}
			if (this.newName.length < 1 || this.newName.length > 20) {
				this.errorMsg = "Invalid username, it must contains min 3 characters and max 16 characters";
				return
			}
			try {
				let _ = await this.$axios.put(`/users/${this.userId}/name`, this.newName,  {headers: { 'Authorization': `${sessionStorage.token}`, 'Content-Type': 'text/plain'}});
				this.userName = this.newName;
				this.updateName();
			} catch (e) {
				if (e.response.data == "Username already exists") {
					this.errorMsg = "Username already used, please choose another one."
				} else {
					this.errorMsg = e.toString();
				}
			}
		},
		handleSearch() {
			this.showSearch = !this.showSearch;
		},
		logout() {
			sessionStorage.clear();
			this.logged = false;
			this.$router.push("/");
		},
		handleLogin() {
			this.logged = true;
			this.userName = sessionStorage.userName;
			this.userImage = sessionStorage.userImage;
		}

	}
}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">Wasa Text</a>
	</header>

	<div class="container-fluid">
		<div class="row">
		<!-- Navigation bar -->
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse" v-show="logged">
				<div class="position-sticky pt-3 sidebar-sticky">

				<!-- Modale utilizzato per la ricerca di un utente con cui aprire una conversazione -->
				<Search :show="showSearch" @close="handleSearch" title="search">
					<template v-slot:header>
						<h3>Users</h3>
					</template>
				</Search>

				<!-- Modale utilizzato per l'aggiornamento dell'username dell'utente -->
				<Search :show="showNameUpdate" @close="handleNameUpdate" title="username">
					<template v-slot:header>
						<h3>Update Username</h3>
					</template>
					<template v-slot:body>
					<!-- Input in cui viene inserito il nuovo nome dell'utente -->
						<form class="username-form">
							<ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
							<input type="text" v-model="newName" placeholder="New username" />
							<button type="submit" @click.prevent="setMyUserName">Update</button>
						</form>
					</template>
				</Search>

				<!-- Modale utilizzato per l'aggiornamento della foto profilo dell'utente -->
				<Search :show="updateProPicIsVisible" @close="handleUpdateProPicToggle" title="photo">
					<template v-slot:header>
						<h3>Update Profile Picture</h3>
					</template>
					<template v-slot:body>
						<!-- Input in cui viene inserita la nuova immagine del profilo dell'utente -->
						<form class="username-form">
							<ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
							<input type="file" ref="file" accept=".jpg,.jpeg" @change="checkFile" />
							<button type="submit" @click.prevent="setMyPhoto">Update</button>
						</form>
					</template>
				</Search>

				<!-- Titolo della NavBar -->
				<h6
					class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
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
					<!-- Apre il modale (solo se l'utente è loggato) per la ricerca di un utente con cui aprire una conversazione -->
					<li class="nav-item m-2" v-if="logged">
						<a class="nav-link" @click="handleSearch">
							<!-- Icona Search -->
							<svg class="feather">
								<use href="/feather-sprite-v4.29.0.svg#search" />
							</svg>
							Search
						</a>
					</li>
					<!-- Esegue il logout (solo se l'utente è loggato) ritornando alla pagina di login -->
					<li class="nav-item m-2" v-if="logged">
						<a class="nav-link" @click="logout">
							<!-- Icona Logout -->
							<svg class="feather">
								<use href="/feather-sprite-v4.29.0.svg#log-out" />
							</svg>
							Logout
						</a>
					</li>
					<!-- Apre il modale per l'inserimento di un nuovo username (Mostato solo se l'utente è loggato) -->
					<li class="nav-item m-2" v-if="logged">
						<button @click="handleNameUpdate">
							<!-- Icona Edit -->
							<svg class="feather">
								<use href="/feather-sprite-v4.29.0.svg#edit" />
							</svg>
						</button>
						Set new username
					</li>
					<!-- Apre il modale per l'inserimento di una nuova immagine del profilo (Mostato solo se l'utente è loggato) -->
					<li class="nav-item m-2" v-if="logged">
					<button @click="handleImageUpdate">
						<!-- Icona Edit -->
						<svg class="feather">
							<use href="/feather-sprite-v4.29.0.svg#edit" />
						</svg>
					</button>
					Set new profile picture
					</li>
					<!-- Mostra l'immagine del profilo e l'username dell'utente loggato -->
					<li class="nav-item m-2" v-if="logged">
					<!-- Immagine del profilo che viene convertito da base64 -->
					<img :src="`data:image/jpg;base64,${userImage}`" alt="Profile Picture" class="profile-picture" />
					<span class="username">{{ userName }}</span>
					</li>
				</ul>
				</div>
			</nav>

			<!-- Contenuto principale della pagina -->
			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView @successful-login="handleLogin" />
			</main>
		</div>
	</div>
</template>

<!--Styles for username and profile picture-->
<style>
.profile-picture {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
  object-fit: cover;
}

.username {
  font-size: 14px;
  font-weight: bold;
}
</style>
