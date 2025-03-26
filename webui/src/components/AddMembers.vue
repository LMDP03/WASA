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
            usernameValidate: new RegExp('^\\w{0,16}$'),
            filteredUsers: [],
            searchText: "",
            owner: localStorage.userName,
            convId: sessionStorage.convId,
            selectedUsers:[],
            groupMembers: JSON.parse(sessionStorage.members),
        };
    },
    methods: {
        closeModal() {
            this.searchText = "";
            this.selectedUsers = [];
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
                        let response = await this.$axios.get(url, { headers: {'Authorization': localStorage.token } });
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
            this.filteredUsers = this.filteredUsers.filter(user => !this.groupMembers.find(m => m.Name === user.Name));
        },
        async addToGroup() {
            try {
                let response = await this.$axios.post(`/users/${localStorage.userId}/conversation/${this.convId}`, this.selectedUsers, { headers: {'Authorization': `${localStorage.token}` } });
                sessionStorage.members = JSON.stringify(response.data);
                this.closeModal();
                window.location.reload();
                this.$router.push('/conversation/groupSettings');
            } catch (e) {
                this.errorMsg = e.toString();
            }
        },
        selectUser(user) {
            if (!this.selectedUsers.find(u => u.Name === user.Name)) {
                this.selectedUsers.push(user);
            }
        },
        removeUser(user) {
            this.selectedUsers = this.selectedUsers.filter(u => u.Name !== user.Name);
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
                        <slot name="header">default header</slot>
                        <button class="like-btn" @click="closeModal">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#x" />
                            </svg>
                        </button>
                    </div>

                    <div class="modal-body">
                        <slot name="body">
                            <div class="search-input">
                                <input type="text" v-model="searchText" placeholder="Search" />
                            </div>
                            <p></p>
                            <div class="btn-group me-2">
                                <button class="btn btn-sm btn-outline-primary" @click="addToGroup">Add To Group</button>
                            </div>
                            <div class="search-results">
                                <div v-for="user in filteredUsers" :key="user.Name" @click="selectUser(user)" class="user">
                                    <p>{{ user.Name }}</p>
                                </div>
                            </div>
                            <div class="selected-users">
                                <h4>Selected Users</h4>
                                <div v-for="user in selectedUsers" :key="user.Name" class="selected-user">
                                    <span>{{ user.Name }}</span>
                                    <button v-if="user.Name !== owner" @click="removeUser(user)">Remove</button>
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
.selected-users {
    margin-top: 20px;
    padding: 10px;
    border-top: 1px solid white;
}

.selected-user {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
}

.selected-user span {
    font-size: 14px;
    font-weight: bold;
}

.selected-user button {
    background: red;
    color: white;
    border: none;
    border-radius: 5px;
    padding: 5px 10px;
    cursor: pointer;
}
</style>