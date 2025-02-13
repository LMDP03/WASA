<script>
import Modal from '../components/Modal.vue';
import Group from '../components/ModalGroup.vue';

export default {
    data: function () {
        return {
            errormsg: null,

            conversations: [],

            

            users: [],
        }
    },
    emits: ['login-success', 'username-changed'],
    methods: {
        async getMyConversations() {
            this.errormsg = null;
            try {
                let response = await this.$axios.get(`/users/${sessionStorage.userId}/conversations`, { headers: { 'Authorization': sessionStorage.token } });
                this.conversations = response.data;
            } catch (e) {
                this.errormsg = e.toString();
            }
        },
        async getConversation(prev) {
            this.errormsg = null;
            this.$axios.get(`/users/${sessionStorage.userId}/conversation/${prev.Id}`, { headers: { 'Authorization': sessionStorage.token } }).then(response => {
                localStorage.convName = response.data.Name;
                localStorage.convId = response.data.Id;
                localStorage.convImage = response.data.Image;
                localStorage.isGroup = response.data.Group;
                localStorage.isNew = false;
                localStorage.members = JSON.stringify(response.Participants);
                this.$router.push(`/conversation`);
            }).catch(e => {
                this.errormsg = e.toString();
            });
        },
    },
    mounted() {
        if (!sessionStorage.token) {
            this.$router.push("/");
            return;
        }
        this.getMyConversations();
    },
    components: { Modal, Group }
}
</script>

<template>
    <div>
        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
        <h1 class="h2">Home page</h1>

        
        <div class="btn-toolbar mb-2 mb-md-0">
            <div class="btn-group me-2">
                <button type="button" class="btn btn-sm btn-outline-secondary" @click="getMyConversations">
                    Refresh
                </button>
            </div>
        </div>
    </div>

        <div v-if="conversations.length !== 0">
            <!-- Mostra le conversazioni dell'utente iterando all'interno di conversations dove sono salvate tutte le conversazioni -->
            <div class="conversations" v-for="prev in conversations" :key="prev.Id">
                <button v-if="prev.Image == ''" type="button" class="btn btn-sm btn-outline-primary" @click="getConversation(prev)">
                    {{ prev.Name }} <br> {{ prev.Sender.Name }}: {{ prev.Text }}
                </button>
                <button type="button" class="btn btn-sm btn-outline-primary" @click="getConversation(prev)" v-else>
                    {{ prev.Name }} <br> {{ prev.Sender.Name }}: Photo
                </button>
                
                <hr>
            </div>
        </div>

        <div v-else>
            <p>Start a conversation</p>
        </div>

        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    </div>
</template>

<style></style>
