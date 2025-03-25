<script>
import AddMembers from '../components/AddMembers.vue';
import SearchUsers from '../components/SearchUsers.vue'

export default {
    data() {
        return{
            errorMsg: "",
            usernameValidate: new RegExp('^\\w{0,16}$'),
            users: [],
            owner: localStorage.userName,
            addMembersIsVisible: false,
            newName: "",
            newImage: null,
            updateNameModalisVisible: false,
            updateImgModalIsVisible: false,
            groupId: sessionStorage.convId,
            groupName: sessionStorage.convName,
            groupImg: sessionStorage.convImg,
            groupMembers: JSON.parse(sessionStorage.members),
        };
    },
    methods: {
        async handleFileChange(event) {
            this.errorMsg = "";
            const file = event.target.files[0];
            if (file.type !== "image/jpeg") {
                this.errorMsg = "Image can only be jpg or jpeg format";
                return;
            }
            if (file.size > 10485760) {
                this.errorMsg = "Image can be at most 10 MB big.";
                return;
            }
            this.newImage = file;
        },
        handleUpdateImage() {
            sessionStorage.convImg = this.newImage;
            this.updateImgModalIsVisible = !this.updateImgModalIsVisible;
            this.newImage = null;
            this.errorMsg = "";
        },
        handleUpdateName() {
            sessionStorage.convName = this.newName;
            this.updateNameModalisVisible = !this.updateNameModalisVisible;
            this.newName = "";
            this.errorMsg = "";
        },
        async setGroupPhoto() {
            this.errorMsg = "";
            const formData = new FormData();
            formData.append('image', this.newImage);
            await this.$axios.put(`/users/${localStorage.userId}/conversation/${this.groupId}/image`, formData, { headers: {'Authorization': localStorage.token } }).then(response => {
                this.groupImg = response.data.Image;
                this.handleUpdateImage();
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        async setGroupName() {
            this.errorMsg = "";
            if (!this.usernameValidate.test(this.newName)) {
                this.errorMsg = "username must be bewteen 3 and 16 alphanumeric characters";
                return;
            }
            try {
                let _ = await this.$axios.put(`users/${localStorage.userId}/conversation/${this.groupId}/name`, this.newName, { headers: {'Authorization': localStorage.token } });
                this.groupName = this.newName;
                this.handleUpdateName();
            } catch (e) {
                this.errorMsg = e.toString();
            }
        },
        async leaveGroup() {
            this.errorMsg = "";
            await this.$axios.delete(`/users/${localStorage.userId}/conversation/${this.groupId}`, { headers: {'Authorization': localStorage.token } }).then(() => {
                this.$router.push('/home');
            }).catch(e => {
                this.errorMsg = e.toString();
            });

        },
        handleAddMembers() {
            this.addMembersIsVisible = !this.addMembersIsVisible;
        },
    },
    emits: ['login-success'],
    components: {AddMembers, SearchUsers},
    }
</script>

<template>
    <div>
        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            <h1 class="h1">
                <img class="conv-photo" :src="`data:image/jpg;base64,${this.groupImg}`" />
                {{ this.groupName }}
            </h1>
            
            <div class="btn-toolbar mb-2 mb-md-0">
                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                <div class="btn-group me-2">
                    <button type="button" class="btn btn-sm btn-outline-primary" @click="handleAddMembers">
                        <svg class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#user-plus" />
                        </svg>
                        Add Members
                    </button>
                </div>
                <div class="btn-group me-2">
                    <button type="button" class="btn btn-sm btn-outline-primary" @click="handleUpdateName">
                        <svg class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#edit" />
                        </svg>
                        Change Group Name

                    </button>
                </div>
                <div class="btn-group me-2">
                    <button type="button" class="btn btn-sm btn-outline-primary" @click="handleUpdateImage">
                        <svg class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#edit" />
                        </svg>
                        Change Group Photo
                    </button>
                </div>
            </div>            
            

            <SearchUsers :show="updateNameModalisVisible" @close="handleUpdateName" title="username">
                <template v-slot:header>
                    <h3>Change Name</h3>
                </template>
                <template v-slot:body>
                    <form class="username-form">
                        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                        <input type="text" v-model="this.newName" placeholder="New Group Name" />
                        <button type="submit" @click.prevent="setGroupName">Update</button>
                    </form>
                </template>
            </SearchUsers>
            <SearchUsers :show="updateImgModalIsVisible" @close="handleUpdateImage" title="photo">
                <template v-slot:header>
                    <h3>Change Picture</h3>
                </template>
                <template v-slot:body>
                    <form class="username-form">
                        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                        <input type="file" ref="file" accept=".jpg,.jpeg" @change="handleFileChange" />
                        <button type="submit" @click.prevent="setGroupPhoto">Update</button>
                    </form>
                </template>
            </SearchUsers>
            <AddMembers :show="addMembersIsVisible" @close="handleAddMembers" title="search">
                <template v-slot:header>
                    <h3>Add Members</h3>
                </template>
            </AddMembers>

        </div>
        <div v-for="user in groupMembers" :key="user.Id">
            <p class="username">
                <img :src="`data:image/jpg;base64,${user.Image}`" class="profile-picture">
                {{ user.Name}}
                <button type="button" v-if="user.Name == owner" class="btn btn-sm btn-outline-primary" @click="leaveGroup">Leave</button>
            </p>
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

.username button {
    background: red;
    color: white;
    border: none;
    border-radius: 5px;
    margin-left: 10px;
}

.username button:hover {
    background: white;
    color: red;
    border: 1px solid red;
}
</style>