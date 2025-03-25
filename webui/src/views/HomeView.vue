<script>
import SearchUsers from '../components/SearchUsers.vue'
import CreateGroup from '../components/CreateGroup.vue'

export default {
    data() {
        return {
            errorMsg: "",
            conversations: [],
            searchModalIsVisible: false,
            createGroupModalIsVisible: false,
            users: [],
            intervalId: null,
        }
    },
    methods: {
        async getMyConversations() {
            this.errorMsg = "";
            try {
                let response = await this.$axios.get(`/users/${localStorage.userId}/conversations`, { headers: { 'Authorization': localStorage.token } });
                this.conversations = response.data;
            } catch (e) {
                this.errorMsg = e.toString();
            }
        },
        goToConversation(preview) {
            sessionStorage.clear();
            sessionStorage.convId = preview.Id;
            sessionStorage.convName = preview.Name;
            sessionStorage.convImg = preview.Image;
            sessionStorage.isGroup = preview.Group;
            this.$router.push('/conversation')
        },
        handleSearchModal() {
            this.searchModalIsVisible = !this.searchModalIsVisible;
        },
        handleCreateGroupModal() {
            this.createGroupModalIsVisible = !this.createGroupModalIsVisible;
        },
    },
    emits: ['login-success', 'username-changed'],
    mounted() {
        if (!localStorage.token) {
            this.$router.push('/')
        }
        this.getMyConversations();
        this.intervalId = setInterval(async () => {
            clearInterval(this.intervalId);
            await this.getMyConversations();
            this.intervalId = setInterval(this.getMyConversations, 1000);
        }, 1000);
    },
    beforeUnmount() {
        if (this.intervalId) {
            clearInterval(this.intervalId);
        }
    },
    components: {SearchUsers, CreateGroup}
}
</script>

<template>
    <div>
        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            <h1>Home</h1>

            <CreateGroup :show="createGroupModalIsVisible" @close="handleCreateGroupModal" title="search">
                <template v-slot:header>
                    <h3>New Group</h3>
                </template>
            </CreateGroup>
            <SearchUsers :show="searchModalIsVisible" @close="handleSearchModal" title="search">
                <template v-slot:header>
                    <h3>Start Conversation</h3>
                </template>
            </SearchUsers>

            <div class="btn-toolbar mb-2 mb-md-0">
                <div class="btn-group me-2">
                    <button type="button" class="btn btn-sm btn-outline-secondary" @click="getMyConversations">
                        <svg class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#refresh-ccw" />
                            Refresh
                        </svg>
                    </button>
                </div>
                <div class="btn-group me-2">
                    <button type="button" class="btn btn-sm btn-outline-primary" @click="handleCreateGroupModal">
                        <svg class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#users" />
                            New Group
                        </svg>
                    </button>
                </div>
                <div class="btn-group me-2">
                    <button type="button" class="btn btn-sm btn-outline-primary" @click="handleSearchModal">
                        <svg class="feather">
                            <use href="/feather-sprite-v4.29.0.svg#user" />
                            Start New Chat
                        </svg>
                    </button>
                </div>
            </div>
        </div>

        <div v-if="conversations.length !== 0">
            <div class="conversations" v-for="preview in conversations" :key="preview.Id">
                <div class= "preview" @click="goToConversation(preview)">
                    <h3>
                        <img :src="`data:image/jpg;base64,${preview.Image}`" alt="Profile Picture" class="profile-picture" />
                        {{ preview.Name }}
                    </h3>
                    <div class="preview-snippet" v-if="preview.LastMessage.MsgId !== 0 && preview.LastMessage.Image === ''">
                        <p>{{ preview.LastMessage.Sender.Name }}: {{ preview.LastMessage.Text }}</p>
                    </div>
                    <div class="preview-snippet" v-if="preview.LastMessage.MsgId !== 0 && preview.LastMessage.Image !== ''">
                        <p>{{ preview.LastMessage.Sender.Name }}: 
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#image" />
                            </svg> {{ preview.LastMessage.Text }}
                        </p>
                    </div>
                </div>
            </div>
        </div>
        <div v-else>
            <p>Start Chatting!</p>
        </div>
        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
    </div>
</template>

<style>
.preview {
    padding: 10px;
    border-bottom: 1px solid lightgray;
    cursor: pointer;
}
.profile-picture {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    margin-right: 10px;
}
.preview-snippet {
    font-size: large;
    margin-left: 60px;
}
</style>