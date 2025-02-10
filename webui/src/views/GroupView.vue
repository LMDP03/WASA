<script>
import AddGroup from '../components/AddMembers.vue'
import SearchUsers from '../components/SearchUsers.vue';

export default {
    data() {
        return {
            errorMsg: "",

            filteredUsers: [],
            owner: sessionStorage.userName,
            selectedUsers: [],
            showAddGroup: false,

            newGroupName: "",
            newGroupImg: null,
            showNameUpdate: false,
            showImageUpdate: false,

            groupName: localStorage.userName,
            groupImage: localStorage.userImage,
            groupId: localStorage.userId,
            groupMembers: JSON.parse(localStorage.users),
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
        handleMemebersUPdate() {
            this.showAddGroup = !this.showAddGroup
        },
        handleNameUpdate() {
            localStorage.userName = this.groupName;
            this.showNameUpdate = !this.showNameUpdate;
            this.newGroupName = "";
            this.errorMsg = "";
        },
        handleImageUpdate() {
            localStorage.userImage = this.newGroupImg;
            this.showImageUpdate = !this.showImageUpdate;
            this.newGroupImg = "";
            this.errorMsg = "";
        },
        async setGroupPhoto() {
            this.errorMsg = "";
            const formData = new FormData();
            formData.append('image', this.newGroupImg);
            this.$axios.put(`/users/${sessionStorage.userId}/conversation/${this.groupId}/image`, formData, { headers: { 'Authorization': `${sessionStorage.token}`}}).then(response => {
                this.groupImage = response.data.Image;
                this.handleImageUpdate();
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        async setGroupName() {
            if (this.newGroupName == this.groupName) {
                this.errorMsg = "Please choose a new username";
                return;
            }
            if (this.newGroupName.length == 0 || this.newGroupName.length > 20) {
                this.errorMsg = "The Group name must be between 1 and 20 characters"
                return;
            }
            this.$axios.put(`/users/${sessionStorage.userId}/conversation/${this.groupId}/name`, formData, { headers: { 'Authorization': `${sessionStorage.token}`}}).then(() => {
                this.groupName = this.newGroupName;
                this.handleNameUpdate();
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        leaveGroup() {
            this.errorMsg = "";
            this.$axios.delete(`/profiles/${sessionStorage.userID}/groups/${this.groupId}`, { headers: { 'Authorization': `${sessionStorage.token}` } }).then(() => {
            // Fa ritornare l'utente alla home dopo l'uscita dal gruppo
                this.$router.push("/home");
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        }
    },
    components: {AddGroup, SearchUsers},
}
</script>

<template>
    <div>
        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            <!-- Group photo -->
            <div class="top-profile-container">
                <img :src="`data:image/jpg;base64,${groupImage}`">
            </div>

            <!-- Modali della pagina -->

            <!-- Modale utlizzato per aggiornare il nome del gruppo -->
            <SearchUsers :show="showNameUpdate" @close="handleNameUpdate" title="username">
                <template v-slot:header>
                    <h3>Update Group Name</h3>
                </template>
                <template v-slot:body>
                <!-- Input per l'inserimento del nuovo nome per il gruppo -->
                <form class="username-form">
                    <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                    <input type="text" v-model="this.newGroupname" placeholder="New group name" />
                    <button type="submit" @click.prevent="setGroupName">Update</button>
                </form>
                </template>
            </SearchUsers>
            <!-- Modale utlizzato per aggiornare la foto del gruppo  -->
            <SearchUsers :show="showImageUpdate" @close="handleImageUpdate" title="photo">
                <template v-slot:header>
                    <h3>Update Group Picture</h3>
                </template>
                <template v-slot:body>
                <!-- Input per l'inserimento della nuova foto per il gruppo -->
                <form class="username-form">
                    <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                    <input type="file" ref="file" accept=".jpg,.jpeg" @change="handleFileChange" />
                    <button type="submit" @click.prevent="setGroupPhoto">Update</button>
                </form>
                </template>
            </SearchUsers>
            <!-- Modale utilizzato per aggiungere utenti al gruppo -->
            <AddGroup :show="showAddGroup" @close="handleMemebersUPdate" title="search">
                <template v-slot:header>
                    <h3>Add to group</h3>
                </template>
            </AddGroup>

            <!-- Body della pagina -->

            <!-- Group name -->
            <h1 class="h1">{{ this.groupName }}</h1>
            <div class="btn-toolbar mb-2 mb-md-0">
                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                <!-- Pulsante per aggiornare il nome del gruppo -->
                <div class="btn-group me-2">
                    <button type="button" class="btn btn-sm btn-outline-secondary" @click="handleNameUpdate">
                        Change group name
                    </button>
                <!-- Pulsante per l'aggiornamento della foto del gruppo -->
                    <button type="button" class="btn btn-sm btn-outline-secondary" @click="handleImageUpdate">
                        Change group photo
                    </button>
                </div>
                <!-- Pulsante per aggiungere utenti al gruppo -->
                <button type="button" class="btn btn-sm btn-outline-primary" @click="handleAddGroupModalToggle">
                    Add to group
                </button>
            </div>
        </div>
        
        <!-- Lista dei membri del gruppo -->
        <div v-for="user in groupMembers">
            <p class="username">
                <!-- Foto profilo dell'utente del gruppo -->
                <img :src="`data:image/jpg;base64,${user.Image}`" class="profile-picture">
                {{ user.Name }}
                <!-- Se l'utente è quello loggato allora mostra il pulstante per uscire dal gruppo -->
                <button type="button" v-if="user.Name == owner" class="btn btn-sm btn-outline-primary" @click="leaveGroup">
                Leave Group
                </button>
            </p>
        </div>

        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    </div>
</template>
  
  <!-- Stili utilizzati per mostarer il nome e la foto degli utenti membri del gruppo -->
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