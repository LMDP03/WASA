<script>
export default {
    props: {
        show: Boolean,
        title: String,
        msg: Object,
    },
    data() {
        return {
            errorMsg: "",
            convs: [],
            filteredConvs: [],
            convId: parseInt(this.$route.params.convId),
        };
    },
    methods: {
        closeMod() {
            this.convs = [];
            window.location.reload();
            this.$emit('close')
        },
        async forwardMessage(destId, user) {
            this.errorMsg = "";
            const url = `users/${sessionStorage.userId}/conversation/${this.convId}/messages/${this.msg.Id}?destId=${destId}`;
            this.$axios.post(url, {}, { headers: { 'Authorization': `${sessionStorage.token}` } }).then(() => {
                localStorage.clear();
                localStorage.userId = user.Id;
                localStorage.userName = user.Name;
                localStorage.userImage = user.Image;
                localStorage.isGroup = response.data.Group;
                this.closeMod();
            }).catch(e => {
                this.errorMsg = e.toString();
            });       
        },
        async getMyConversations() {
            this.errorMsg = "";
            try {
                const url = `users/${sessionStorage.userId}/conversation`
                let response = await this.$axios.get(url, { headers: { 'Authorization': `${sessionStorage.token}` } });
                this.convs = response.data;
            } catch (e) {
                this.errorMsg = e.toString();
            }
        }
        
    },
    watch: {
        show() {
            this.getMyConversations();
        }
    },
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

                    <body>
                        <div class="search-results">
                            <div v-for="response in convs" :key="response.Id" @click="forwardMessage(response.Id, response.LastMessage.Sender)">
                                <RouterLink :to="'/conversation/' + response.Id" class="custom-link" replace force>
                                    <div class="user" v-if="response.Id !== convId">
                                        <p>{{ response.LastMessage.Sender.Name }}</p>
                                    </div>
                                </RouterLink>
                            </div>
                        </div>
                    </body>

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