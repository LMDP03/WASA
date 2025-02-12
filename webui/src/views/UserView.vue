<script>
import Search from '../components/SearchUsers.vue';

export default {
    data() {
        return {
            errorMsg: "",

            filteredUsers: [],
            owner: sessionStorage.userName,
            selectedUsers: [],

            newName: "",
            newImg: null,
            showName: false,
            showImage: false,

            userName: sessionStorage.userName,
            userImage: sessionStorage.userImage,
            userId: sessionStorage.userId,
        };
    },
    emits: ['successful-login'],
    methods: {
        async checkFile(event) {
            this.errorMsg = "";
            const file = event.target.files[0];
            if (file.type !== "image/jpeg") {
                this.errorMsg = "File not supported: only jppg and jpeg images allowed";
                return;
            }
            if (file.size > 5242880) {
                this.errorMsg = "Image is too big, max size allowed is 5 MB"
                return;
            }
            this.newGroupImg = file;
        },
        handleNameUpdate() {
            sessionStorage.userName = this.userName;
            this.showName = !this.showName;
            this.newName = "";
            this.errorMsg = "";
        },
        handleImageUpdate() {
            sessionStorage.userImage = this.userImage;
            this.showImage = !this.showImage;
            this.newImg = "";
            this.errorMsg = "";
        },
        async setMyPhoto() {
            this.errorMsg = "";
            const formData = new FormData();
            formData.append('image', this.newImg);
            this.$axios.put(`/users/${sessionStorage.userId}/image`, formData, { headers: { 'Authorization': sessionStorage.token}}).then(response => {
                this.userImage = response.data.Image;
                this.handleImageUpdate();
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        async setMyUserName() {
            this.errorMsg = "";
            if (this.newName == this.userName) {
                this.errorMsg = "Please choose a new user name";
                return;
            }
            if (this.newName.length < 3 || this.newName.length > 16) {
                this.errorMsg = "The Group name must be between 1 and 20 characters"
                return;
            }
            this.$axios.put(`/users/${sessionStorage.userId}/name`, this.newName, { headers: { 'Authorization': sessionStorage.token}}).then(response => {
                this.userName = response.data.Name;
                this.handleNameUpdate();
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        logout() {
			sessionStorage.clear();
			this.logged = false;
			this.$router.push("/");
            window.location.reload();
		},
    },
    mounted() {
        if (!sessionStorage.token) {
            this.$router.push("/");
            return;
        }
    },
    components: {Search}
}
</script>

<template>
    <div>
        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            <!-- Group photo -->
            <div class="top-profile-container">
                <img :src="`data:image/jpg;base64,${userImage}`">
            </div>
            <h1 class="h1">{{ userName }}</h1>

            <!-- Modali della pagina -->

            <!-- Modale utlizzato per aggiornare il nome utente -->
            <Search :show="showName" @close="handleNameUpdate" title="username">
                <template v-slot:header>
                    <h3>Change Name</h3>
                </template>
                <template v-slot:body>
                    <!-- Input per l'inserimento del nuovo nome per il gruppo -->
                    <form class="username-form">
                        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                        <input type="text" v-model="newName" placeholder="New Name" />
                        <button type="submit" @click.prevent="setMyUserName">Change</button>
                    </form>
                </template>
            </Search>

            <!-- Modale utlizzato per aggiornare la foto profilo  -->
            <Search :show="showImage" @close="handleImageUpdate" title="photo">
                <template v-slot:header>
                    <h3>Change Photo</h3>
                </template>
                <template v-slot:body>
                <!-- Input per l'inserimento della nuova foto per il gruppo -->
                <form class="username-form">
                    <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                    <input type="file" ref="file" accept=".jpg,.jpeg" @change="checkFile" />
                    <button type="submit" @click.prevent="setMyPhoto">Update</button>
                </form>
                </template>
            </Search>


            <!-- Group name -->
            <div class="btn-toolbar mb-2 mb-md-0">
                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                <!-- Pulsante per aggiornare il nome del gruppo -->
                <button type="button" class="btn btn-sm btn-outline-primary" @click="handleNameUpdate">
                    Change Name
                </button>
                <!-- Pulsante per l'aggiornamento della foto del gruppo -->
                <button type="button" class="btn btn-sm btn-outline-primary" @click="handleImageUpdate">
                    Change Photo
                </button>
                <button type="button" class="btn btn-sm btn-outline-primary" @click="logout">
                    Logout
                </button>
            </div>
        </div>

        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
    </div>
</template>
  
  
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