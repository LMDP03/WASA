<script>
import CommentMessage from '../components/CommentMessage.vue';
import ForwardMessage from '../components/ForwardMessage.vue';

export default {
    data() {
        return {
            errorMsg: "",
            convId: sessionStorage.convId,
            convName: sessionStorage.convName,
            convImg: sessionStorage.convImg,
            isGroup: sessionStorage.isGroup,
            userId: sessionStorage.userId,
            text: "",
            image: null,
            conversations: [],
            users: [],
            messages: [],
            comments: [],
            messageToForward: null,
            messageToComment: null,
            messageToRespond: null,
            intervalId: null,
            commentModalIsVisible: false,
            forwardModalIsVisible: false,
            emojis: ["😀", "😂", "😍", "😎", "😭", "😡", "🎉", "❤️", "👍", "🔥"],
            ownerId: localStorage.userId,
        }
    },
    methods: {
        async handleFileChange(event) {
            this.errorMsg = "";
            const file = event.target.files[0];
            if (file.type !== "image/jpeg") {
                this.errorMsg = "Image can only be jpg or jpeg format";
                return;
            }
            if (file.size > 10485760) {
                this.errorMsg = "Image can be at most 10 MB big.";
                return;
            }
            this.image = file;
        },
        async getConversation() {
            this.errorMsg = "";
            this.$axios.get(`/users/${localStorage.userId}/conversation/${this.convId}`, { headers: {'Authorization': localStorage.token } }).then(response => {
                this.messages = response.data.Messages;
                this.isGroup = response.data.Group;
                sessionStorage.members = JSON.stringify(response.data.Participants);
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        check() {
            if (!this.convId) {
                this.startConversation();
            } else {
                this.sendMessage();
            }
        },
        async startConversation() {
            this.errorMsg = "";
            const formData = new FormData();
            formData.append('text', this.text);
            if (this.image != null) {
                formData.append('image', this.image);
            }
            await this.$axios.post(`/users/${this.ownerId}/conversations/private?rcvId=${this.userId}`, formData, { headers: { 'Authorization': localStorage.token } }).then(response => {
                sessionStorage.convId = response.data.Id;
                this.messages = response.data.Messages;
                window.location.reload();
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        async sendMessage() {
            this.errorMsg = "";
            const formData = new FormData();
            formData.append('text', this.text);
            if (this.image != null) {
                formData.append('image', this.image);
            }
            if (this.messageToRespond) {
                await this.$axios.post(`/users/${localStorage.userId}/conversation/${this.convId}/messages?responseTo=${this.messageToRespond.MsgId}`, formData, { headers: { 'Authorization': localStorage.token } }).then(() => {
                    this.text = "";
                    this.image = null;
                    this.messageToRespond = null;
                }).catch(e => {
                    this.errorMsg = e.toString();
                });
            } else {
                await this.$axios.post(`/users/${localStorage.userId}/conversation/${this.convId}/messages`, formData, { headers: { 'Authorization': localStorage.token } }).then(() => {
                    this.text = "";
                    this.image = null;
                }).catch(e => {
                    this.errorMsg = e.toString();
                });
            }
            
        },
        async deleteMessage(msg) {
            this.errorMsg = "";
            await this.$axios.delete(`/users/${localStorage.userId}/conversation/${this.convId}/messages/${msg.MsgId}`, { headers: { 'Authorization': localStorage.token } }).then(() => {
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        handleCommentModal(msg) {
            this.messageToComment = msg;
            this.commentModalIsVisible = !this.commentModalIsVisible;
        },
        handleForwardModal(msg) {
            this.messageToForward = msg;
            this.forwardModalIsVisible = !this.forwardModalIsVisible;
        },
        replyToMessage(msg) {
            this.messageToRespond = msg;
        },
        goToInfo() {
            this.$router.push(`/conversation/groupSettings`);
        },
        async uncommentMessage(msgId) {
            this.errorMsg = "";
            const url = `/users/${localStorage.userId}/conversation/${this.convId}/messages/${msgId}/reactions`;
            await this.$axios.delete(url, { headers: { 'Authorization': localStorage.token } }).then(() => {}).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        countEmojis(msg, emoji) {
            return msg.Reactions.filter(c => c.Emoji === emoji).length;
        },
        checkReactions(msg) {
            return msg.Reactions.filter(c => c.Sender.Id == this.ownerId).length;
        }
    },
    emits: ['login-succes'],
    mounted() {
        if (!localStorage.token) {
            this.$router.push("/");
            return;
        }
        if (this.convId) {
            this.getConversation();
            this.intervalId = setInterval(async () => {
                clearInterval(this.intervalId);
                await this.getConversation();
                this.intervalId = setInterval(this.getConversation, 1000)
            }, 1000);
        }
    },
    beforeUnmount(){
        if (this.intervalId) {
            clearInterval(this.intervalId);
        }
    },
    components: {CommentMessage, ForwardMessage},
}
</script>

<template>
    <div>
        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            <div class="top-profile-container">                
                <img class="conv-photo" :src="`data:image/jpg;base64,${this.convImg}`">
                <h1 class="h1">{{ this.convName }}</h1>
                <button v-show="isGroup" type="button" class="btn btn-sm btn-outline-primary" @click="goToInfo">Group Settings</button>
            </div>          

            <CommentMessage :show="commentModalIsVisible" :msg="messageToComment" @close="handleCommentModal" title="comments">
                <template>
                    <h3>Comments</h3>
                </template>
            </CommentMessage>
            <ForwardMessage :show="forwardModalIsVisible" :msg="messageToForward" @close="handleForwardModal" title="conversations">
                <template v-slot:header>
                    <h3>Conversations</h3>
                </template>
            </ForwardMessage>
        </div>

        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>

        <div class="messages" v-for="message in messages" :key="message.MsgId">
            <div v-if="message.ResponseTo.Sender.Id !== 0">
                <p>Response to: {{ message.ResponseTo.Sender.Name }}</p>
                <span v-if="message.ResponseTo.Image !== ''"> Photo; </span>
                <p>{{ message.ResponseTo.Text }}</p>
            </div>
            <p>{{ message.Sender.Name }}</p>
            <img class="msg_photo" v-if="message.Image !== ''" :src="`data:image/jpg;base64,${message.Image}`" alt="Message Photo">
            <p v-if="message.Text !== ''">{{ message.Text }}</p>
            <p>
                {{ message.Timestamp }}
                <span v-if="message.Checkmark === 'received'">✔️</span>
                <span v-else>✔️✔️</span>
            </p>
            <div v-for="cmt in emojis" :key="cmt">
                <p v-if="countEmojis(message, cmt) > 0">{{ cmt }}: {{ countEmojis(message, cmt) }}</p>
            </div>
            <div class="btn-group-me-2">
                <button type="button" class="btn btn-sm btn-outline-secondary" @click="replyToMessage(message)">Reply</button>
                <button type="button" class="btn btn-sm btn-outline-secondary" @click="handleCommentModal(message)">Comment</button>
                <button type="button" class="btn btn-sm btn-outline-secondary" @click="handleForwardModal(message)">Forward</button>
                <button type="button" class="btn btn-sm btn-outline-secondary" @click="deleteMessage(message)">Delete</button>
                <button v-if="checkReactions(message) !== 0" type="button" class="btn btn-sm btn-outline-secondary" @click="uncommentMessage(message.MsgId)">Uncomment</button>
            </div>
        </div>

        <div class="new-message">
            <svg v-if="messageToRespond" class="feather" @click="messageToRespond = null">
                    <use href="/feather-sprite-v4.29.0.svg#x" />
                </svg>
            <div v-if="messageToRespond" class="reply-snippet">
                <h6>Replying to: {{ messageToRespond.Sender.Name }}</h6>
                <p v-if="messageToRespond.Image === ''">{{ messageToRespond.Text }}</p>
                <p v-else>
                    <svg class="feather">
                        <use href="/feather-sprite-v4.29.0.svg#image" />
                    </svg> {{ messageToRespond.Text }}
                </p>
            </div>
            <div class="input-group">
                <input type="file" class="form-control" ref="file" accept=".jpg,.jpeg" @change="handleFileChange" />
                <input type="text" class="form-control" v-model="text" placeholder="Write a message">
                <button class="btn btn-outline-primary" @click="check">Send</button>
            </div>
        </div>
    </div>
</template>

<style>
.top-profile-container {
    width: 100vw;
    height: auto;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
}

.msg_photo {
    width: 15%;
    height: 15%;
}

.h1 {
    font-weight: bold;
}

.conv-photo {
    width: 120px;
    height: 120px;
    border-radius: 50%;
    margin-right: 10px;
}

.messages {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    justify-content: space-evenly;
    padding: 10px;
    border-bottom: 1px solid lightgray;
    overflow-y: scroll;
}

.new-message {
    position: fixed;
    bottom: 0;
    right: 0;
    width: 83vw;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-evenly;
    padding: 10px;
    border-top: 1px solid lightgray;
    z-index: 999;
    background-color: white;
}

.reply-snippet {
    display: flex;
    flex-direction: column;
    align-items: start;
    justify-content: space-between;
    padding: 10px;
    border-bottom: 1px solid lightgray;
    background-color: inherit;
    margin-left: 60px;
    width: inherit;
    left: 0;
}

.new-message svg {
    align-self: flex-end;
    cursor: pointer;
    height: 20px;
    width: 20px;
    margin-right: 10px;
}


.input-group {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    padding: 10px;
    background-color: inherit;
    width: 100%;
}

.input-group input[type="text"] {
    width: 50vw;
}


</style>