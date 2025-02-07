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
            searchName: "",
            groupName: "",
            filteredUsers: [],
            selectedUsers: [],
            username: sessionStorage.userName,
            groupId: localStorage.userId,

        };
    },
    methods: {
        close() {
            this.searchName = "";
            this.selectedUsers = [];
            this.$emit('close');
        },
        async filterUsers() {
            this.errorMsg = "";
            this.filteredUsers = this.users;

            if (this.searchName.length > 16) {
                this.errorMsg = "Invalid username, it can contain only letters and numbers for a maximum of 16 characters.";
                this.filteredUsers = [];
                return;
            }

            if (this.title == "search") {
                try {
                    const url =  `users/${sessionStorage.userId}/others?srcName=${this.searchName}`
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
                this.filteredUsers = this.users.filter(user => user.name.toLowerCase().includes(this.searchName.toLowerCase()));
            }
        },
        async addToGroup() {
            if (this.groupName.length == 0 || this.groupName.length > 20) {
                this.errorMsg = "Group name must be between 1 and 20 characters";
                return;
            }
            try {
                const url = `users/${sessionStorage.userId}/conversation/${this.groupId}`;
                var newMembers = [];
                for (let i in this.selectedUsers) {
                    newMembers[i] = this.selectedUsers[i].name;
                }
                let response = await this.$axios.post(url, {
                    participants: newMembers,
                }, {headers: { 'Authorization': `${sessionStorage.token}`}});
                localStorage.clear();
                _ = localStorage.users.push(JSON.stringify(response.data.participants));
                this.close();
                this.$router.push(`/conversations/${response.data.id}`)
            } catch (e) {
                this.errorMsg = e.toString();
            }
        },
        selectUser(user) {
            if (!this.selectedUsers.find(u => u.name == user.name)) {
                this.selectedUsers.push(user);
            }
        },
        removeUser(name) {
            this.selectedUsers = this.selectedUsers.filter(user => user.name != name);
        },
    },
    watch: {
        searchName() {
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
                        <slot name="header">default header</slot>
                        <button class="like-btn" @click="close">
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
                                <input type="text" v-model="searchName" placeholder="Search" />
                            </div>
                            <p></p>
                            <div class="btn-group me-2">
                                <button class="btn btn-sm btn-outline-primary" @click="addToGroup">Add to Group</button>
                            </div>
            
                            <!-- Risultati della ricerca -->
                            <div class="search-results">
                                <div v-for="user in filteredUsers" :key="user.id" @click="selectUser(user)" class="user">
                                    <p v-if="user.name != username">{{ user.name }}</p>
                                </div>
                            </div>
            
                            <!-- Lista di utenti selezionati -->
                            <div class="selected-users">
                                <h4>Selected Users:</h4>
                                <span class="selected-user">{{ username }}</span>
                                <div v-for="user in selectedUsers" :key="user.id" class="selected-user">
                                    <span>{{ user.name }}</span>
                                    <button v-if="user.name != username" @click="removeUser(user.name)">Remove</button>
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
  