<script setup>
import { RouterLink, RouterView } from 'vue-router'
import ModalUsers from './components/ModalUsers.vue';
</script>
<script>
export default {
	data() {
		return {
			searchModalVisible: false,

			logged: sessionStorage.token ? true : false,

			userId: sessionStorage.userId,
			userName: sessionStorage.userName,
			userImage: sessionStorage.userImage,

			nameModalVisible: false,
			newName: "",

			imageModalVisible: false,
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
		updateImage() {
			sessionStorage.userImage = this.userImage;
			this.imageModalVisible = !this.imageModalVisible;
			this.newImage = "";
			this.errorMsg = "";
		},
		async setMyPhoto() {
			this.errorMsg = "";

			const formData = new FormData();
			formData.append('image', this.newImage);
			try {
				let response = await this.$axios.put(`/users/${this.userId}/image`, formData, {headers: { 'Authorization': `${sessionStorage.token}`, 'Content-Type': 'application/json'}});
				this.userImage = response.data.image;
				this.updateImage();
			} catch (e) {
				this.errorMsg = e.toString();
			}
		},
		async updateName() {
			sessionStorage.userName = this.userName;
			this.nameModalVisible = !this.nameModalVisible;
			this.newName = "";
			this.errorMsg = "";
		},
		async setMyUserName() {
			this.errorMsg = "";

			if (this.newName == this.userName) {
				this.errorMsg = "Choose a new name different from the current one.";
				return
			}
			if (this.newUsername.length < 3 || this.newUsername.length > 16) {
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
		updateSearchModal() {
			this.searchModalVisible = !this.searchModalVisible;
		},
		logout() {
			sessionStorage.clear();
			this.logged = false;
			this.$router.push("/");
		},
		login() {
			this.logged = true;
			this.userName = sessionStorage.userName;
			this.userImage = sessionStorage.userImage;
		}

	}
}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASA-Text</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">

			<!-- Navigation bar -->
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse" v-show="logged">
				<div class="position-sticky pt-3 sidebar-sticky">

					<!--User search Modal-->
					<Modal v-show="searchModalVisible" @close="updateSearchModal" title="Search Users">
						<template v-slot:header>
							<h3>Search Users</h3>
						</template>
					</Modal>

					<!--Username update Modal-->
					<ModalUsers v-show="nameModalVisible" @close="updateName", title="Change Name">
						<template v-slot:header>
							<h3>Change Username</h3>
						</template>
						<template v-slot:body>
							<form class="username-form">
								<ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
								<input type="text" v-model="newName" placeholder="Insert new username" />
								<button type="submit" @click.prevent="setMyUserName">Update</button>
							</form>
						</template>
					</ModalUsers>

					<!--User picture update Modal-->
					<ModalUsers v-show="imageModalVisible" @close="updateImage", title="Change Picture">
						<template v-slot:header>
							<h3>Change Profile Picture</h3>
						</template>
						<template v-slot:body>
							<form class="username-form">
								<ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
								<input type="file" ref="file" accept=".jpg, .jpeg" @change="checkFile" placeholder="Upload new picture" />
								<button type="submit" @click.prevent="setMyPhoto">Update</button>
							</form>
						</template>
					</ModalUsers>


					<!--Navigation Bar Title-->
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>

					<!--Navigation Bar Links-->
					<ul class="nav flex-column">
						<li class="nav-item">
							<!-- Link alla pagina Home -->
							<RouterLink to="/home" class="nav-link m-2">
								<!--Home Icon-->
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
								Home
							</RouterLink>
						</li>
						<!--Opens user search modal (only if user is logged in)-->
						<li class="nav-item m-2" v-if="logged">
								<a class="nav-link" @click="updateSearchModal">
									<!-- Search Icon-->
									<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#search" />
									</svg>
									Search
								</a>
							</li>
							<!-- If the user is logged in, logs them out-->
						<li class="nav-item m-2" v-if="logged">
							<a class="nav-link" @click="logout">
								<!-- Icona Logout -->
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#log-out" />
								</svg>
								Logout
							</a>
						</li>
						<!--If the user isn't logged opens the login page-->
						<li class="nav-item" v-else>
							<RouterLink to="/session" class="nav-link m-2">
								<!-- Icona Login -->
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#key" />
								</svg>
								Login
							</RouterLink>
						</li>
						<!--Opens the Modal for the username change-->
						<li class="nav-item m-2" v-if="logged">
							<button @click="updateName">
								<!--Edit Icon-->
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#edit" />
								</svg>
							</button>
							Set new username
						</li>
						<!--Opens the Modal for the profile picture change-->
						<li class="nav-item m-2" v-if="logged">
							<button @click="updateImage">
								<!--Edit Icon-->
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

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView @successful-login="login" />
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
