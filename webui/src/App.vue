<script setup>
import { RouterLink, RouterView } from 'vue-router'
import SearchUsers from './components/SearchUsers.vue'
</script>

<script>
export default {
    data() {
        return {
            errorMsg: "",
            searchModalIsVisible: false,
            isLoggedIn: localStorage.token ? true : false,
            useId: localStorage.userId,
            userName: localStorage.userName,
            userImage: localStorage.userImage,
            updateNameIsVisible: false,
            newName: "",
            updateImageIsVisible: false,
            newImage: null,
            usernameValidate: new RegExp('^\\w{0,16}$'),
        };
    },
    methods: {
        handleSearchModal() {
            this.searchModalIsVisible = !this.searchModalIsVisible;
        },
        logout() {
            localStorage.clear();
            this.isLoggedIn = false;
            this.$router.push('/');
        },
        handleLoginSucces() {
            this.isLoggedIn = true;
            this.userId = localStorage.userId;
            this.userName = localStorage.userName;
            this.userImage = localStorage.userImage;
        },
        handleFileChange(event) {
            this.errorMsg = "";
            const file = event.target.files[0];
            if (file.type !== "image/jpeg") {
                this.errorMsg = "Images can only be jpg or jpeg";
                return;
            }
            if (file.size > 10485760) {
                this.errorMsg = "Image can be at most 10 MB large";
                return;
            }
            this.newImage = file;
        },
        handleUpdateImage() {
            localStorage.userImage = this.userImage;
            this.updateImageIsVisible = !this.updateImageIsVisible;
            this.newImage = null;
            this.errorMsg = "";
        },
        handleUpdateName() {
            localStorage.userName = this.userName;
            this.updateNameIsVisible = !this.updateNameIsVisible;
            this.newName = "";
            this.errorMsg = "";
        },
        async setMyPhoto() {
            this.errorMsg = "";
            const formData = new FormData();
            formData.append('image', this.newImage);

            this.$axios.put(`/users/${localStorage.userId}/image`, formData, { headers: { 'Authorization': `${localStorage.token}` } }).then(response => {
                this.userImage = response.data.Image;
                this.handleUpdateImage();
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        async setMyUserName() {
            if (this.userName == this.newName) {
                this.errorMsg = "Please choose a username different from your actual one";
                return;
            }
            if (!this.usernameValidate.test(this.newName)) {
                this.errorMsg = "Username must be between 3 and 16 alphanumerical characters.";
                return;
            }
            try {
                let _ = await this.$axios.put(`/users/${localStorage.userId}/name`, this.newName, { headers: { 'Authorization': `${localStorage.token}` } });
                this.userName = this.newName;
                this.handleUpdateName();
            } catch (e) {
                if (e.response.data == "Username already exists\n") {
                    this.errorMsg = "Username already taken; please choose another one.";
                } else {
                    this.errorMsg = e.toString();
                }
            }
        }
    }
}
</script>

<template>
    <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
        <a class="navbar-bran col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">Wasa Text</a>
    </header>

    <div class="container-fluid">
        <div class="row">
            <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse" v-show="isLoggedIn">
                <div class="position-sticky pt-3 sidebar-sticky">

                    <SearchUsers :show="searchModalIsVisible" @close="handleSearchModal" title="search">
                        <template v-slot:header>
                            <h3>Users</h3>
                        </template>
                    </SearchUsers>
                    <SearchUsers :show="updateNameIsVisible" @close="handleUpdateName" title="username">
                        <template v-slot:header>
                            <h3>Update UserName</h3>
                        </template>
                        <template v-slot:body>
                            <form class="username-form">
                                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                                <input type="text" v-model="newName" placeholder="Choose a new name" />
                                <button type="submit" @click.prevent="setMyUserName">Update</button>
                            </form>
                        </template>
                    </SearchUsers>
                    <SearchUsers :show="updateImageIsVisible" @close="handleUpdateImage" title="photo">
                        <template v-slot:header>
                            <h3>Updated Profile Picture</h3>
                        </template>
                        <template v-slot:body>
                            <form class="username-form">
                                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                                <input type="file" ref="file" accept=".jpg,.jpeg" @change="handleFileChange" />
                                <button type="submit" @click.prevent="setMyPhoto">Update</button>
                            </form>
                        </template>
                    </SearchUsers>

                    <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
                        <span>Menu</span>
                    </h6>

                    <ul class="nav flex-column">
                        <li class="nav-item m-2" v-if="isLoggedIn">
                            <img :src="`data:image/jpg;base64,${userImage}`" alt="Profile Picture" class="profile-picture" />
                            <span class="username">{{ userName }}</span>
                        </li>
                        <li class="nav-item" v-else>
                            <RouterLink to="/session" class="nal-link m-2">
                                <svg class="feather">
                                    <use href="/feather-sprite-v4.29.0.svg#log-in" />
                                </svg>
                                Login
                            </RouterLink>
                        </li>
                        <li class="nav-item">
                            <RouterLink to="/home" class="nav-link m-2">
                                <svg class="feather">
                                    <use href="/feather-sprite-v4.29.0.svg#home" />
                                </svg>
                                Home
                            </RouterLink>
                        </li>
                        <li class="nav-item m-2" v-if="isLoggedIn">
                            <a class="nav-link" @click="handleSearchModal">
                                <svg class="feather">
                                    <use href="/feather-sprite-v4.29.0.svg#search" />
                                </svg>
                                Search Users
                            </a>
                        </li>
                        <li class="nav-item m-2" v-if="isLoggedIn">
                            <a class="nav-link" @click="logout">
                                <svg class="feather">
                                    <use href="/feather-sprite-v4.29.0.svg#log-out" />
                                </svg>
                                Log-Out
                            </a>
                        </li>
                    </ul>

                </div>
            </nav>

            <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
                <RouterView @login-success="handleLoginSucces" />
            </main>

        </div>
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