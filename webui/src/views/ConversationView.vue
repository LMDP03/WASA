<script>
import Forward from '../components/ForwardMessage.vue'
import Comment from '../components/Comments.vue'

export default {
    data() {
        return {
            errorMsg: "",

            convImg: localStorage.convImage,
            convId: localStorage.convId,
            convName: localStorage.convName,
            isGroup: localStorage.isGroup,

            receiverId: localStorage.receiverId,

            userId: sessionStorage.userId,
            
            text: "",
            image: null,

            messages: [],

            showForward: false,
            showComments: false,

            messageToFordward: null,
            messageToComment: null,

            comments: [],
        }
    },
    emits: ['successful-login'],
    methods: {
        async checkFile(event) {
			this.errorMsg = "";
			const file = event.target.files[0]; // Prende il file inserito dall'utente
			if (file.type !== "image/jpeg") {
				this.errorMsg = "Unsupported image type: only jpg or jpeg images allowed.";
				return
			}
			if (file.size > 5242880) {
				this.errorMsg = "File size exceeded: MAX size 5MB."
				return
			}
			this.image = file;
		},
        async getConversation() {
            this.errorMsg = "";
			localStorage.clear();
			const url = `users/${sessionStorage.userId}/conversation/${this.convId}`;
			await this.$axios.get(url, { headers: { 'Authorization': `${sessionStorage.token}` } }).then(response => {
				this.messages = response.data.Messages;
                localStorage.members = response.data.Participants;
			}).catch(e => {
				this.errorMsg = e.toString();
			});
		},
        checkType() {
            if (this.isGroup || this.receiverId == 0) {
                this.sendMessage();
            }
            else {
                this.startConversation();
            }
        },
        async startConversation() {
            this.errorMsg = "";
            const formData = new FormData();
            formData.append('text', this.text);
            if (this.image != null) {
                formData.append('image', this.image);
            }
            const url = `/users/${this.userId}/conversations/private?rcvId=${this.receiverId}`;
            await this.$axios.post(url, formData, { headers: { 'Authorization': `${sessionStorage.token}`}}).then(response => {
                this.convId = response.data.Id;
                this.receiverId = 0;
                window.location.reload();
            }).catch( e => {
                this.errorMsg = e.toString();
            });
        },
        async deleteMessage(msgId) {
            this.errorMsg = "";
            const url = `/users/${this.userId}/conversation/${this.convId}/messages/${msgId}`;
            await this.$axios.delete(url, { headers: { 'Authorization': `${sessionStorage.token}` } }).then(() => {
            }).catch(e => {
                this.errorMsg = e.toString()
            });
        },
        async sendMessage() {
            this.errorMsg = "";
            const formData = new FormData();
            formData.append('text', this.text);
            if (this.image != null) {
                formData.append('image', this.image);
            }
            const url = `/users/${this.userId}/conversation/${this.convId}/messages`;
            await this.$axios.post(url, formData, { headers: { 'Authorization': `${sessionStorage.token}`}}).then(() => {
                this.text = "";
                this.image = null;
                this.getConversation();
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        async uncommentMessage(msgId) {
            this.errorMsg = "";
            const url = `/users/${this.userId}/conversation/${this.convId}/messages/${msgId}/reactions`;
            await this.$axios.delete(url, { headers: { 'Authorization': `${sessionStorage.token}`}}).then(() => {}).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        handleForward(msg) {
            this.messageToFordward = msg;
			this.showForward = !this.showForward;
		},
        handleComment(msg) {
            this.messageToComment = msg;
            this.showComments = !this.showComments;
        },
        goToGroupInfo() {
            this.$router.push(`/conversation/groupSettings`);
        }
    },
    mounted() {
        if (!sessionStorage.token) {
            this.$router.push("/");
            return;
        }
        if (this.convId != undefined && !isNaN(this.convId)) {
            this.getConversation();
        }
    },
    components: {Forward, Comment}
}
</script>

<template>
    <div>
        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            
            <!-- Controlla se la conversazione è con un gruppo o con un utente -->
            <div v-if="isGroup">
                <!-- Se è un gruppo mostra il nome del gruppo -->
                <img :src="`data:image/jpg;base64,${this.convImg}`">
                <h1 class="h1 clickable" @click="goToGroupInfo">{{ this.convName}}</h1>
            </div>
            <div v-else>
                <!-- Se è un utente mostra il nome dell'utente -->
                <img :src="`data:image/jpg;base64,${this.convImg}`">
                <h1 class="h1">{{ this.convName }}</h1>
            </div>
  
            <!-- Modali della pagina -->
    
            <!-- Modale utilizzato per lasciare un commento a un messaggio -->
            <Comment :show="showComments" :comments="comments" :msg="messageToComment" @close="handleComment" title="comments">
                <template>
                    <h3>Comments</h3>
                </template>
            </Comment>
            <!-- Modale utilizzato per selezionare una conversazione in cui inoltrare un messaggio -->
            <Forward :show="showForward" :msg="messageToFordward" @close="handleForward" title="forward">
                <template v-slot:header>
                    <h3>Forward To</h3>
                </template>
            </Forward>

        </div>

        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>

        <!-- Lista dei messaggi della conversazione -->
        <div class="messages" v-for="message in messages" :key="message.MsgId">
            <!-- Mostra il contenuto del messaggio, con chi lo ha mandato, il contenuto e il timeStamp -->
            <p v-if="message.Text !== '' || message.Image !== ''">
                {{ message.Sender.Name }}
            </p>
            <br>
            <div v-if="message.ResponseTo.MsgId != 0">
                <p v-if="message.ResponseTo.Image == ''" >
					{{ message.ResponseTo.Sender.Name }}: {{ message.ResponseTo.Text }}
				</p>
				<p v-else >
					{{ message.ResponseTo.Sender.Name }}:
                    <svg class="feather">
                        <use href="/feather-sprite-v4.29.0.svg#image" />
                    </svg>
                    {{ message.ResponseTo.Text }}
                </p>
            </div>
            <br>
            <img class="msg_photo" v-if="message.Image !== ''" :src="`data:image/jpg;base64,${message.Image}`" alt="Message Photo">
            {{ message.Text }}
            <br>
            <p v-if="message.text !== '' || message.Image !== ''">
                {{ message.TimeStamp }}
            </p>
            <div class="btn-group me-2">
                <!-- Pulsante per commentare il messaggio -->
                <div v-for="cmt in message.Reactions" :key="cmt.Sender.Id">
                    <p>{{ cmt.Sender.Name }}: {{ cmt.Emoji }}</p>
                    <button v-if="cmt.Sender.Id == userId" type="button" class="btn btn-sm btn-outline-secondary" @click="uncommentMessage(message.MsgId)">
                        Uncomment
                    </button>
                </div>
                <button type="button" class="btn btn-sm btn-outline-secondary" @click="handleComment(message)">
                    Comment
                </button>
                <!-- Pulsante per inoltrare il messaggio in un'altra conversazione -->
                <button type="button" class="btn btn-sm btn-outline-secondary" @click="handleForward(message)">
                    Forward
                </button>
                <!-- Pulsante per eliminare il messaggio dalla conversazione -->
                <button v-if="message.Sender.Id == userId" type="button" class="btn btn-sm btn-outline-secondary" @click="deleteMessage(message.MsgId)">
                    Delete
                </button>
                
            </div>
            <hr v-if="message.text !== '' && message.photo !== ''">
        </div>


            <!-- Body della pagina -->
        <div class="btn-toolbar mb-2 mb-md-0">
            <!-- Form per inviare una foto -->
            <div class="btn-group me-2">
                <form @submit.prevent="checkType">
                    <input type="file" ref="file" accept=".jpg,.jpeg" @change="checkFile" />
                    <input type="text" class="form-control" v-model="text" placeholder="New Message">
                    <!-- Pulsante per inviare il messaggio -->
                    <button type="submit" class="btn btn-sm btn-outline-primary">
                        Send Message
                    </button>
                </form>
            </div>
        </div>
    </div>
</template>
  
<style>
/* Stile utilizzato per visualizzare l'immagine profilo dell'utente o del gruppo con cui si sta conversando */
.top-profile-container {
    width: auto;
    height: auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-between;
}

/* Stile utilizzato per visualizzare la foto di un messaggio */
.msg_photo {
    width: 25%;
    height: 25%;
}

/* Stile utilizzato nel caso in cui la conversazione è con un gruppo */
.clickable {
    cursor: pointer;
    color: white;
    background-color: black;
    text-decoration: underline;
}

.clickable:hover {
    background-color: white;
    color: black;
}
</style>