<script>
export default  {
    props: {
        show: Boolean,
        users: Array,
        title: String,
    },
    data() {
        return {
            errorMsg: "",
            searchText: "",
            groupName: "",
            filteredUsers: [],
            selectedUsers: [],
            owner: sessionStorage.userName,

        };
    },
    methods: {
        closeMod() {
            this.searchText = "";
            this.groupName = "";
            this.selectedUsers = [];
            this.$emit('close');
        },
        async filterUsers() {
            this.errorMsg = "";
            this.filteredUsers = this.users;

            if (this.searchText.length > 0) {
                if (this.searchName.length > 16) {
                    this.errorMsg = "Invalid username, it can contain at most 16 characters.";
                    this.filteredUsers = [];
                    return;
                }

                if (this.title === "Search") {
                    try {
                        const url =  `users/${sessionStorage.userId}/others?srcName=${this.searchText}`
                        let response = await this.$axios.get(url, { headers: { 'Authorization': `${sessionStorage.token}` } });
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
        async createGroup() {
            if (this.groupName.length < 1 || this.groupName.length > 20) {
                this.errorMsg = "Group name must be between 1 and 20 characters";
                return;
            }
            try {
                const url = `users/${sessionStorage.userId}/conversations/group`;
                let response = await this.$axios.post(url, {
                    name: this.groupName,
                    participants: this.selectedUsers,
                }, {headers: { 'Authorization': sessionStorage.token}});
                localStorage.clear();
                localStorage.convId = response.data.Id;
                localStorage.convName = response.data.Name;
                localStorage.convImage = response.data.Image;
                localStorage.isGroup = response.data.Group;
                this.closeMod();
                this.$router.push(`/conversation`)
            } catch (e) {
                this.errorMsg = e.toString();
            }
        },
        selectUser(user) {
            if (!this.selectedUsers.find(u => u.Name === user.Name)) {
                this.selectedUsers.push(user);
            }
        },
        removeUser(name) {
            this.selectedUsers = this.selectedUsers.filter(user => user.Name !== name);
        },
    },
    watch: {
        searchText() {
            this.filterUsers();
        },
        show() {
            this.filteredUsers = this.users;
        }
    }
}
</script>


<template>
    <Transition name="modal">
        <div v-if="show" class="modal-mask">
            <div class="modal-wrapper">
                <div class="modal-container">

                    <div class="modal-header">
                        <button class="like-btn" @click="closeMod">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#x" />
                            </svg>
                        </button>
                    </div>
                    
                    <div class="modal-body">
                        <slot name="body">
                            <!-- Selezione del nome del gruppo -->
                            <div class="search-input">
                                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                                <input type="text" v-model="groupName" placeholder="Select group name" />
                            </div>
                            <!-- Campo di ricerca -->
                            <div class="search-input">
                                <input type="text" v-model="searchText" placeholder="Search" />
                            </div>
                            <p></p>
                            <div class="btn-group me-2">
                                <button class="btn btn-sm btn-outline-primary" @click="createGroup">Create Group</button>
                            </div>
            
                            <!-- Risultati della ricerca -->
                            <div class="search-results">
                                <div v-for="user in filteredUsers" :key="user.Id" @click="selectUser(user)" class="user">
                                    <p> {{ user.Name }}</p>
                                </div>
                            </div>
            
                            <!-- Lista di utenti selezionati -->
                            <div class="selected-users">
                                <h4>Selected Users:</h4>
                                <span class="selected-user">{{ owner }}</span>
                                <div v-for="user in selectedUsers" :key="user.Id" class="selected-user">
                                    <span>{{ user.Name }}</span>
                                    <button v-if="user.Name !== owner" @click="removeUser(user.Name)">Remove</button>
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
    border-top: 1px solid #ccc;
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
