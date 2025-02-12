<script>
import Add from '../components/AddMembers.vue'
import Search from '../components/SearchUsers.vue';

export default {
    data() {
        return {
            errorMsg: "",

            filteredUsers: [],
            owner: sessionStorage.userName,
            selectedUsers: [],
            showAdd: false,

            newGroupName: "",
            newGroupImg: null,
            showName: false,
            showImage: false,

            groupName: localStorage.convName,
            groupImage: localStorage.convImage,
            groupId: localStorage.convId,
            groupMembers: JSON.parse(localStorage.members),
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
        handleAddMembers() {
            this.showAdd = !this.showAdd;
        },
        handleNameUpdate() {
            localStorage.convName = this.groupName;
            this.showName = !this.showName;
            this.newGroupName = "";
            this.errorMsg = "";
        },
        handleImageUpdate() {
            localStorage.convImage = this.groupImage;
            this.showImage = !this.showImage;
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
            this.errorMsg = "";
            if (this.newGroupName == this.groupName) {
                this.errorMsg = "Please choose a new group name";
                return;
            }
            if (this.newGroupName.length == 0 || this.newGroupName.length > 20) {
                this.errorMsg = "The Group name must be between 1 and 20 characters"
                return;
            }
            this.$axios.put(`/users/${sessionStorage.userId}/conversation/${this.groupId}/name`, this.newGroupName, { headers: { 'Authorization': `${sessionStorage.token}`}}).then(response => {
                this.groupName = response.data.Name;
                this.handleNameUpdate();
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        leaveGroup() {
            this.errorMsg = "";
            this.$axios.delete(`/profiles/${sessionStorage.userId}/groups/${this.groupId}`, { headers: { 'Authorization': `${sessionStorage.token}` } }).then(() => {
            // Fa ritornare l'utente alla home dopo l'uscita dal gruppo
                this.$router.push("/home");
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        }
    },
    components: {Add, Search},
}
</script>

<template>
    <div>
        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            <!-- Group photo -->
            <div class="top-profile-container">
                <img :src="`data:image/jpg;base64,${groupImage}`">
                <h1>{{ groupName }}</h1>
            </div>

            <!-- Modali della pagina -->

            <!-- Modale utlizzato per aggiornare il nome del gruppo -->
            <Search :show="showName" @close="handleNameUpdate" title="username">
                <template v-slot:header>
                    <h3>Change Group Name</h3>
                </template>
                <template v-slot:body>
                    <!-- Input per l'inserimento del nuovo nome per il gruppo -->
                    <form class="username-form">
                        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                        <input type="text" v-model="newGroupName" placeholder="New Group Name" />
                        <button type="submit" @click.prevent="setGroupName">Change</button>
                    </form>
                </template>
            </Search>

            <!-- Modale utlizzato per aggiornare la foto del gruppo  -->
            <Search :show="showImage" @close="handleImageUpdate" title="photo">
                <template v-slot:header>
                    <h3>Change Group Photo</h3>
                </template>
                <template v-slot:body>
                <!-- Input per l'inserimento della nuova foto per il gruppo -->
                <form class="username-form">
                    <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                    <input type="file" ref="file" accept=".jpg,.jpeg" @change="checkFile" />
                    <button type="submit" @click.prevent="setGroupPhoto">Update</button>
                </form>
                </template>
            </Search>

            <!-- Modale utilizzato per aggiungere utenti al gruppo -->
            <Add :show="showAddGroup" @close="handleAddMembers" title="Search">
                <template v-slot:header>
                    <h3>Add New Members</h3>
                </template>
            </Add>

            <!-- Body della pagina -->

            <!-- Group name -->
            <div class="btn-toolbar mb-2 mb-md-0">
                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                <!-- Pulsante per aggiornare il nome del gruppo -->
                <button type="button" class="btn btn-sm btn-outline-primary" @click="handleNameUpdate">
                    Change group name
                </button>
                <!-- Pulsante per l'aggiornamento della foto del gruppo -->
                <button type="button" class="btn btn-sm btn-outline-primary" @click="handleImageUpdate">
                    Change group photo
                </button>
                <!-- Pulsante per aggiungere utenti al gruppo -->
                <button type="button" class="btn btn-sm btn-outline-primary" @click="handleAddMembers">
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

        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
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