<script>
export default {
    props: {
        show: Boolean,
        title: String,
        msg: Object,
        conversations: Array,
    },
    data() {
        return {
            errorMsg: "",
            convId: sessionStorage.convId,
            searchText: "",
            filteredConvs: [],
            selectedConvs: [],
            filteredUsers: [],
            usernameValidate: new RegExp('^\\w{0,16}$'),
        };
    },
    methods: {
        closeModal() {
            this.searchText = "";
            this.$emit('close');
        },
        async filterConvs() {
            this.errorMsg = "";
            this.filteredConvs = this.conversations;
            if (this.searchText.length > 0) {
                if (this.searchText.length > 16 || !this.usernameValidate.test(this.searchText)) {
                    this.errorMsg = "Invalid name, it can contain at most 16 alphanumerical characters.";
                    this.filteredConvs = [];
                    return;
                }

                try {
                    const url = `/users/${localStorage.userId}/conversations?srcName=${this.searchText}`;
                    let response = await this.$axios.get(url, { headers: {'Authorization': localStorage.token } });
                    if (response.data == null) {
                        this.filteredConvs = [];
                        return;
                    }
                    this.filteredConvs = response.data;
                    let response2 = await this.$axios.get(`/users/${localStorage.userId}/others?srcName=${this.searchText}`, { headers: {'Authorization': localStorage.token } });
                    if (response2.data == null) {
                        this.filteredUsers = [];
                        return;
                    }
                    this.filteredUsers = response2.data;
                    this.filteredUsers = this.filteredUsers.filter(user => !this.filteredConvs.some(conv => conv.Name === user.Name && conv.Group === false));
                    this.filteredConvs = this.filteredConvs.concat(this.filteredUsers);
                } catch (e) {
                    this.errorMsg = e.toString();
                    this.filteredUsers = [];
                    this.filteredConvs = [];
                }

            }
        },
        selectConv(conv) {
            if (conv.hasOwnProperty('Group') && !this.selectedConvs.some(c => c.Id === conv.Id && c.hasOwnProperty('Group'))) {
                this.selectedConvs.push(conv);
            } else if (!conv.hasOwnProperty('Group') && !this.selectedConvs.some(c => c.Id === conv.Id && !c.hasOwnProperty('Group'))) {
                this.selectedConvs.push(conv);
            }
        },
        removeConv(conv) {
            if (conv.hasOwnProperty('Group')) {
                this.selectedConvs = this.selectedConvs.filter(c => (c.Id !== conv.Id && c.hasOwnProperty('Group')) || !c.hasOwnProperty('Group'));
            } else {
                this.selectedConvs = this.selectedConvs.filter(c => (c.Id !== conv.Id && !c.hasOwnProperty('Group')) || c.hasOwnProperty('Group'));
            }
        },
        async forwardMessage() {
            this.errorMsg = "";
            const destinations = [];
            for (let conv of this.selectedConvs) {
                var dest = {
                    id: conv.Id,
                    isConv: conv.hasOwnProperty('Group'),
                };
                destinations.push(dest);
            }
            try {
                const url = `/users/${localStorage.userId}/conversation/${this.convId}/messages/${this.msg.MsgId}`;
                let response = await this.$axios.post(url, destinations, { headers: { 'Authorization': localStorage.token } });
                sessionStorage.clear();
                sessionStorage.convId = response.data.Id;
                sessionStorage.convName = response.data.Name;
                sessionStorage.convImg = response.data.Image;
                sessionStorage.isGroup = response.data.Group;
                sessionStorage.members = JSON.stringify(response.data.Participants);
                sessionStorage.messages = response.data.Messages;
                window.location.reload();
                this.closeModal();
            } catch (e) {
                this.errorMsg = e.toString();
            }
        },
    },
    watch : {
        searchText() {
            this.filterConvs();
        },
        show() {
            this.filteredConvs = this.conversations;
        }
    }
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
                                <button class="btn btn-sm btn-outline-primary" @click="forwardMessage">Forward Message</button>
                            </div>
                            <div class="search-results">
                                <div v-for="conv in filteredConvs" :key="conv.Id" @click="selectConv(conv)" class="user">
                                    <p>{{ conv.Name }}</p>
                                </div>
                            </div>
                            <div class="selected-users">
                                <h4>Selected Destinations</h4>
                                <div v-for="conv in selectedConvs" :key="conv.Id" class="selected-user">
                                    <span>{{ conv.Name }}</span>
                                    <button @click="removeConv(conv)">Remove</button>
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
.custom-link {
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
    background-color: rgba(0, 0, 0, 0.5);
    display: table;
    transition: 0.3s ease;
}

.modal-wrapper {
    display: table-cell;
    vertical-align: middle;
}

.modal-container {
    width: 300px;
    margin: 0px auto;
    background-color: white;
    border-radius: 2px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
    transition: all 0.3 ease;
}

.modal-header {
    height: 70x;
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
    border: 1px solid lightgray;
}

.search-results {
    font-size: 15px;
    padding: 10px 15px;
    border-bottom: 1px solid white;
    cursor: pointer;
    max-height: 200px;
    overflow-y: scroll;
}

.modal-default-button {
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
    border-style: 3px;
    border: 1px solid lightgray;
}

.username-form button {
    margin-bottom: 15px;
}
</style>