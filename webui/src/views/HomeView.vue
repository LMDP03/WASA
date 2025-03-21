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
            this.$router.push('/conversations')
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
        <div class="d-flex justify content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            <h1 class="h2">Home</h1>

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
                <button v-if="preview.LastMessage.Image == ''" type="button" class="btn btn-sm btn-outline-primary" @click="goToConversation(preview)">
                    {{ preview.Name }} <br> {{ preview.LastMessage.Sender.Name }}: {{ preview.LastMessage.Text }}
                </button>
                <button v-else type="button" class="btn btn-sm btn-outline-primary" @click="goToConversation(preview)">
                    {{ preview.Name }} <br> {{ preview.LastMessage.Sender.Name }}: 
                    <svg class="feather">
                        <use href="/feather-sprite-v4.29.0.svg#user" />
                    </svg> {{ preview.LastMessage.Text }}
                </button>
            </div>
        </div>
        <div v-else>
            <p>Start Chatting!</p>
        </div>
        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
    </div>
</template>

<style></style>