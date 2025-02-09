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
            groupId: localStorage.userId,

        };
    },
    methods: {
        closeMod() {
            this.searchText = "";
            this.selectedUsers = [];
            this.$emit('close');
        },
        async filterUsers() {
            this.errorMsg = "";
            this.filteredUsers = this.users;
            
            if (this.searchText.length > 0) {
                if (this.searchName.length > 16) {
                    this.errorMsg = "Invalid username, it can contain only letters and numbers for a maximum of 16 characters.";
                    this.filteredUsers = [];
                    return;
                }
                
                if (this.title === "search") {
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
        async addToGroup() {
            try {
                const url = `users/${sessionStorage.userId}/conversation/${this.groupId}`;
                var newMembers = [];
                for (let i in this.selectedUsers) {
                    newMembers[i] = this.selectedUsers[i].Name;
                }
                let response = await this.$axios.post(url, {participants: newMembers,}, {headers: { 'Authorization': `${sessionStorage.token}`, 'Content-Type': 'application/json'}});
                localStorage.clear();
                localStorage.users = JSON.stringify(response.data)
                this.closeMod();
                window.location.reload();
                this.$router.push(`/conversations/${this.groupId}`)
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
                        <h3>default header</h3>
                        <button class="like-btn" @click="closeMod">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#x" />
                            </svg>
                        </button>
                    </div>
        
                    <div class="modal-body">
                        <body>
                            <!-- Selezione del nome del gruppo -->
                            <div class="search-input">
                                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                                <input type="text" v-model="groupName" placeholder="Select group name" />
                            </div>
                            <!-- Campo di ricerca -->
                            <div class="search-input">
                                <input type="text" v-model="searchText" placeholder="search" />
                            </div>
                            
                            <div class="btn-group me-2">
                                <button class="btn btn-sm btn-outline-primary" @click="addToGroup">Add to Group</button>
                            </div>
            
                            <!-- Risultati della ricerca -->
                            <div class="search-results">
                                <div v-for="user in filteredUsers" :key="user.Id" @click="selectUser(user)" class="user">
                                    <p v-if="user.Name !== owner">{{ user.Name }}</p>
                                </div>
                            </div>
            
                            <!-- Lista di utenti selezionati -->
                            <div class="selected-users">
                                <h4>Selected Users:</h4>
                                <div v-for="user in selectedUsers" :key="user.Id" class="selected-user">
                                    <span>{{ user.Name }}</span>
                                    <button v-if="user.name !== owner" @click="removeUser(user.name)">
                                        <svg class="feather">
                                            <use href="/feather-sprite-v4.29.0.svg#x" />
                                        </svg>
                                    </button>
                                </div>
                            </div>
                        </body>
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
    border-top: 1px solid gray;
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
    background: white;
    color: grey;
    border: none;
    border-radius: 5px;
    padding: 5px 10px;
    cursor: pointer;
}
.modal-header button svg {
    width: 20px;
    height: 20px;
}
.selected-user button:hover {
    color: red;
}
</style>
  