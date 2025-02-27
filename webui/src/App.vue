<!-- 

Pagina principale dell'applicazione, contiene la navbar con i link alle pagine principali 
dell'applicazione e i modali per l'aggiornamento dell'username e dell'immagine del profilo dell'utente loggato.

L'utente può:
- visualizzare le conversazioni con altri utenti o gruppi andando alla home 
- cercare un utente con cui aprire una nuova conversazione
- aggiornare l'username
- aggiornare l'immagine del profilo
- effettuare il logout
- visualizzare la sua immagine del profilo e il suo username

-->

<script setup>
import { RouterLink, RouterView } from 'vue-router'
import Modal from './components/Modal.vue'
</script>
<script>
export default {
    data() {
        return {
            errorMsg: null,

            showSearch: false,
            isLoggedIn: sessionStorage.token ? true : false,

            userID: sessionStorage.userId,

            username: sessionStorage.userName,
            photo: sessionStorage.userImage,

            showNameUpdate: false,
            newUsername: "", 

            showImageUpdate: false,
            newUserImage: null,

            showSearch: false,
            showGroup: false,

        }
    },
    methods: {
        async handleFileChange(event) {
            this.errorMsg = "";
            const file = event.target.files[0];
            if (file.type !== "image/jpeg") {
                this.errorMsg = "File type not supported, only jpg and jpeg are allowed";
                return
            }
            if (file.size > 5242880) {
                this.errorMsg = "File size is too big. Max size is 5MB";
                return
            }
            this.newUserImage = file;
        },
        handleImageUpdate() {
        sessionStorage.userImage = this.photo;  
        this.showImageUpdate = !this.showImageUpdate;
        this.newUserImage = "";
        this.errorMsg = "";
        },
        handleNameUpdate() {
        sessionStorage.userName = this.username;  
        this.showNameUpdate = !this.showNameUpdate;
        this.newUsername = "";
        this.errorMsg = "";
        },
        async setMyPhoto() {
            this.errorMsg = "";

            const formData = new FormData();
            formData.append('image', this.newUserImage);

            this.$axios.put(`/users/${sessionStorage.userID}/Image`, formData, { headers: { 'Authorization': `${sessionStorage.token}` } }).then(response => {
                this.photo = response.data.Image; 
                this.handleImageUpdate();
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        async setMyUserName() {
            // COntrolla se l'username inserito dall'utente è uguale a quello attuale
            if (this.newUsername == this.username) {
                this.errorMsg = "You must enter a new username";
                return
            }
            // Controlla se l'username inserito dall'utente ha una lunghezza valida
            if (this.newUsername.length < 3 || this.newUsername.length > 16) {
                this.errorMsg = "Invalid username, it must contains min 3 characters and max 16 characters";
                return
            }
            try {
                let _ = await this.$axios.put(`/users/${sessionStorage.userId}/name`, this.newUsername, { headers: { 'Authorization': `${sessionStorage.token}` } })
                this.username = this.newUsername;
                this.handleNameUpdate();
            } catch (e) {
                if (e.response.data == "Username already exist\n") {
                    this.errorMsg = "This username is already taken. Please try another one.";
                } else {
                    this.errorMsg = e.toString();
                }
            }
        },
        handleSearch() {
            this.showSearch = !this.showSearch;
        },
        handleGroup() {
            this.showGroup = !this.showGroup;
        },
        logout() {
            sessionStorage.clear();
            this.isLoggedIn = false;
            this.$router.push("/");
            },
        handleLoginSuccess() {
        this.isLoggedIn = true;
        this.userID = sessionStorage.userId;
        this.username = sessionStorage.userame;
        this.photo = sessionStorage.photo;
        },

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
      <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse" v-show="isLoggedIn">
        <div class="position-sticky pt-3 sidebar-sticky">

            <!-- Modale utilizzato per la ricerca di un utente con cui aprire una conversazione -->
            <Modal :show="showSearch" @close="handleSearch" title="search">
                <template v-slot:header>
                <h3>Users</h3>
                </template>
            </Modal>

            <!-- Modale utilizzato per l'aggiornamento dell'username dell'utente -->
            <Modal :show="showNameUpdate" @close="handleNameUpdate" title="username">
                <template v-slot:header>
                <h3>Update Username</h3>
                </template>
                <template v-slot:body>
                <!-- Input in cui viene inserito il nuovo nome dell'utente -->
                <form class="username-form">
                    <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                    <input type="text" v-model="newUsername" placeholder="New username" />
                    <button type="submit" @click.prevent="setMyUserName">Update</button>
                </form>
                </template>
            </Modal>

            <!-- Modale utilizzato per l'aggiornamento della foto profilo dell'utente -->
            <Modal :show="showImageUpdate" @close="handleImageUpdate" title="photo">
                <template v-slot:header>
                    <h3>Update Profile Picture</h3>
                    </template>
                    <template v-slot:body>
                    <!-- Input in cui viene inserita la nuova immagine del profilo dell'utente -->
                    <form class="username-form">
                        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                        <input type="file" ref="file" accept=".jpg,.jpeg" @change="handleFileChange" />
                        <button type="submit" @click.prevent="setMyPhoto">Update</button>
                    </form>
                </template>
            </Modal>

                <!-- Modale utilizzato per la creazione di un nuovo gruppo -->
            <Group :show="showGroup" @close="handleGroup" title="search">
                <template v-slot:header>
                <h3>Select users</h3>
                </template>
            </Group>
            <!-- Modale utilzzato per la ricerca degli utenti con cui aprire una nuova conversazione -->
            <Modal :show="showSearch" @close="handleSearch" title="search">
                <template v-slot:header>
                <h3>Users</h3>
                </template>
            </Modal>


          <!-- Titolo della NavBar -->
          <h6
            class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
            <span>General</span>
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
            <li class="nav-item m-2" v-if="isLoggedIn">
              <a class="nav-link" @click="handleSearch">
                <!-- Icona Search -->
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#search" />
                </svg>
                Search
              </a>
            </li>

            
            <!-- Esegue il logout (solo se l'utente è loggato) ritornando alla pagina di login -->
            <li class="nav-item m-2" v-if="isLoggedIn">
              <a class="nav-link" @click="logout">
                <!-- Icona Logout -->
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#log-out" />
                </svg>
                Logout
              </a>
            </li>
            <!-- Apre la pagina di login (solo se l'utente non è loggato) -->
            <li class="nav-item" v-else>
              <RouterLink to="/session" class="nav-link m-2">
                <!-- Icona Login -->
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#key" />
                </svg>
                Login
              </RouterLink>
            </li>
            <!-- Apre il modale per l'inserimento di un nuovo username (Mostato solo se l'utente è loggato) -->
            <li class="nav-item m-2" v-if="isLoggedIn">
              <button @click="handleNameUpdate">
                <!-- Icona Edit -->
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#edit" />
                </svg>
              </button>
              Set new username
            </li>
            <!-- Apre il modale per l'inserimento di una nuova immagine del profilo (Mostato solo se l'utente è loggato) -->
            <li class="nav-item m-2" v-if="isLoggedIn">
              <button @click="handleImageUpdate">
                <!-- Icona Edit -->
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#edit" />
                </svg>
              </button>
              Set new profile picture
            </li>
            <!-- Mostra l'immagine del profilo e l'username dell'utente loggato -->
            <li class="nav-item m-2" v-if="isLoggedIn">
              <!-- Immagine del profilo che viene convertito da base64 -->
              <img :src="`data:image/jpg;base64,${photo}`" alt="Profile Picture" class="profile-picture" />
              <span class="username">{{ username }}</span>
            </li>
          </ul>
        </div>
      </nav>

      <!-- Contenuto principale della pagina -->
      <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
        <RouterView @login-success="handleLoginSuccess"/>
      </main>
    </div>
  </div>
</template>


<!-- Stili per l'immagine del profilo e dell'username dell'utente loggato -->
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
