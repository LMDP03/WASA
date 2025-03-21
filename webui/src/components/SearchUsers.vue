<script>
export default {
    props: {
        show: Boolean,
        users: Array,
        title: String,
    },
    data() {
        return {
            errorMsg: "",
            searchText: "",
            usernameValidate: new RegExp('^\\w{0,16}$'),
            filteredUsers: [],
        };
    },
    methods: {
        closeModal() {
            this.searchText = "";
            this.$emit('close');
        },
        async filterUsers() {
            this.errorMsg = "";
            this.filteredUsers = this.users;
            if (this.searchText.length > 0) {
                if (this.searchText.length > 16 || !this.usernameValidate.test(this.searchText)) {
                    this.errorMsg = "Invalid username, it can contain at most 16 alphanumerical characters.";
                    this.filteredUsers = [];
                    return;
                }

                if (this.title === "search") {
                    try {
                        const url = `/users/${localStorage.userId}/others?srcName=${this.searchText}`;
                        let response = await this.$axios.get(url, {headers: {'Authorization': `${localStorage.token}`}});
                        if (response.data == null) {
                            this.filteredUsers = [];
                            return;
                        }
                        this.filteredUsers = response.data;
                    } catch (e) {
                        this.errorMsg = e.toString();
                        this.filteredUsers = [];
                    }
                } else {
                    this.filteredUsers = this.users.filter(user => user.Name.toLowerCase().includes(this.searchText.toLowerCase()));
                }
            }
        },
        async selectUser(user) {
            sessionStorage.clear();
            sessionStorage.userId = user.Id;
            sessionStorage.convImg = user.Image;
            sessionStorage.convName = user.Name;
            sessionStorage.isGroup = false;
            await this.$router.push('/conversation');
            this.closeModal();
        },
    },
    watch: {
        searchText() {
            this.filterUsers();
        },
        show() {
            this.filteredUsers = this.users;
        }
    },
};
</script>

<template>
    <Transition name="modal">
        <div v-if="show" class="modal-mask">
            <div class="modal-wrapper">
                <div class="modal-container">
                    <div class="modal-header">
                        <slot name="header">deafult header</slot>
                        <button class="like-btn" @click="closeModal">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#x" />
                            </svg>
                        </button>
                    </div>

                    <div class="modal-body">
                        <slot name="body">
                            <div class="search-input">
                                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                                <input type="text" v-model="searchText" placeholder="Search" />
                            </div>
                            <div class="search-results">
                                <div v-for="user in filteredUsers" :key="user.Name" @click="selectUser(user)">
                                    <div class="user">
                                        <p>{{ user.Name }}</p>
                                    </div>
                                </div>
                            </div>
                        </slot>
                    </div>
                </div>
            </div>
        </div>
    </Transition>
</template>

<style>
.custom-link{
    color: inherit;
    text-decoration: none;
}
.modal-mask {
    position: fixed;
    z-index: 9998;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: black;
    display: table;
    transition: opacity 0.3 ease;
}
.modal-wrapper {
    display: table-cell;
    vertical-align: middle;
}

.modal-container {
    width: 350px;
    margin: 0px auto;
    background-color: white;
    border-radius: 2px;
    box-shadow: 0 2px 8px black;
    transition: all 0.3s ease;
}

.modal-header {
    height: 70px;
    padding: 20px 15px 10px 15px;
}

.modal-header h3 {
    margin-top: 0;
    font-size: 25px;
    color: black;
}

.modal-header button {
    color: gray;
    background: none;
    border: none;
    padding: 5px;
    line-height: 12px;
    font-size: 15px;
}

.modal-header button svg {
    width: 20px;
    height: 20px;
}

.search-input {
    padding: 0 15px;
}

.search-input input {
    height: 30px;
    width: 100%;
    outline: none;
    border-radius: 3px;
    border: 1px solid gray;
}

.search-results {
    font-size: 15px;
    padding: 10px 15px;
    border-bottom: 1px solid white;
    cursor: pointer;
    max-height: 200px;
    overflow-y: scroll;
}

.modal-deafult-button {
    float: right;
}

.username-form {
    display: flex;
    flex-direction: column;
    padding: 0 15px;
}

.username-form input {
    margin-bottom: 10px;
    margin-top: 5px;
    outline: none;
    border-radius: 3px;
    border: 1px solid gray;
}

.username-form button {
    margin-bottom: 15px;
}

.user {
    padding: 10px 0;
    border-bottom: 1px solid gray;
    cursor: pointer;
}
</style>