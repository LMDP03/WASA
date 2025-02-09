<script>
import Forwarding from '../components/ForwardMessage.vue'
import Commenting from '../components/Comments.vue'

export default {
    data() {
        return {
            errorMsg: "",
            convImg: localStorage.userImage,
            convId: parseInt(this.$route.params.convId),
            userToSend: localStorage.userName,
            userIdToSend: localStorage.userId,
            isGroup: localStorage.isGroup,
            text: "",
            image: "",
            some_data: [],
            showUserSearch: false,
            showComments: false,
            messageToFordward: null,

            // Messaggio da commentare
            messageToComment: null,

            // Commenti del messaggio selezionato
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
			localStorage.clear();
			const url = `users/${sessionStorage.userId}/conversation/${this.convId}`;
			this.$axios.get(url, { headers: { 'Authorization': `${sessionStorage.token}` } }).then(response => {
				this.some_data = response.data.Messages;
			}).catch(e => {
				this.errorMsg = e.toString();
			});
		},
        check() {
            if (isNaN(this.convId) || this.convId == undefined) {
                this.startConversation();
            }
            else {
                this.sendMessage();
            }
        },
        startConversation() {
            this.errorMsg = "";
            this.$axios.post(`/users/${sessionStorage.userId}/conversations/private?rcvId=${localStorage.userId}`, {
                text: this.text,
                image: this.image,
            }, { headers: { 'Authorization': `${sessionStorage.token}`, 'Content-Type': 'application/json'}}).then(response => {
                this.convId = response.data.Id;
                this.$router.push(`/conversation/${this.convId}`);
            }).catch( e => {
                this.errorMsg = e.toString();
            });
        },
        async deleteMessage(msgId) {
            this.errorMsg = "";
            this.$axios.delete(`/users/${sessionStorage.userId}/conversation/${this.convId}/messages/${msgId}`, { headers: { 'Authorization': `${sessionStorage.token}` } }).then(() => {
                this.getConversation();
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
            this.$axios.post(`/users/${sessionStorage.userId}/conversation/${this.convId}/messages`, {
                text: this.text,
                image: this.image,
            }, { headers: { 'Authorization': `${sessionStorage.token}`, 'Content-Type': 'application/json'}}).then(() => {
                this.text = "";
                this.image = null;
                this.getConversation();
            }).catch(e => {
                this.errorMsg = e.toString();
            });
        },
        handleSearchMod(msg) {
            this.messageToFordward = msg;
			this.showUserSearh = !this.showUserSearh;
		},
        handleCommentMod(cmt, msgComs) {
            this.messageToComment = cmt;
            this.comments = msgComs;
            this.showComments = !this.showComments;
        },
        goToGroupInfo() {
            this.$router.push(`conversation/${this.userIdToSend}`);
        }
    },
    mounted() {
        if (!sessionStorage.token) {
            this.$router.push("/");
            return;
        }
        if (this.convId != undefined && !isNaN(this.convId)) {
            this.getConversation()
        }
    },
    components: {Forwarding, Commenting}
}
</script>

<template>
    <div>
        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            <!-- User photo -->
            <div class="top-profile-container">
                <img :src="`data:image/jpg;base64,${this.convImg}`">
            </div>
            <!-- Controlla se la conversazione è con un gruppo o con un utente -->
            <div v-if="isGroup">
                <!-- Se è un gruppo mostra il nome del gruppo -->
                <h1 class="h1 clickable" @click="goToGroupInfo">{{ this.userToSend }}</h1>
            </div>
            <div v-else>
                <!-- Se è un utente mostra il nome dell'utente -->
                <h1 class="h1">{{ this.userToSend }}</h1>
            </div>
  
            <!-- Modali della pagina -->
    
            <!-- Modale utilizzato per lasciare un commento a un messaggio -->
            <Commenting :show="showComments" :comments="comments" :msg="messageToComment" @close="handleCommentMod" title="comments">
                <template>
                    <h3>Comments</h3>
                </template>
            </Commenting>
            <!-- Modale utilizzato per selezionare una conversazione in cui inoltrare un messaggio -->
            <Forwarding :show="showUserSearch" :msg="messageToFordward" @close="handleSearchMod" title="conversations">
                <template v-slot:header>
                    <h3>Conversations</h3>
                </template>
            </Forwarding>
   
        </div>
        <!-- Lista dei messaggi della conversazione -->
        <div class="messages" v-for="response in some_data" :key="response.message.messageId">
            <!-- Mostra il contenuto del messaggio, con chi lo ha mandato, il contenuto e il timeStamp -->
            <p v-if="response.Text !== '' || response.Image !== ''">
                {{ response.Sender.Name }}
            </p>
            <p v-if="response.Text !== ''">
                {{ response.Text }}
            </p>
            <p></p>
            <!-- Mostra la foto contenuta nel messaggio nel caso in cui il messaggio contiene una foto -->
            <img class="msg_photo" v-if="response.Image !== ''":src="`data:image/jpg;base64,${response.Image}`" alt="Message Photo">
            <p v-if="response.text !== '' || response.Image !== ''">
                {{ response.timeMsg }}
            </p>
            <div class="btn-group me-2">
                <!-- Pulsante per inoltrare il messaggio in un'altra conversazione -->
                <button type="button" class="btn btn-sm btn-outline-secondary" @click="handleSearchMod(response.message)">
                    Forward Message
                </button>
                <!-- Pulsante per eliminare il messaggio dalla conversazione -->
                <button type="button" class="btn btn-sm btn-outline-secondary" @click="deleteMessage(response.message.messageId)">
                    Delete message
                </button>
                <!-- Pulsante per commentare il messaggio -->
                <button type="button" class="btn btn-sm btn-outline-secondary" @click="handleCommentMod(response.message, response.comments)">
                    Comment message
                </button>
            </div>
            <hr v-if="response.message.text !== '' && response.message.photo !== ''">
        </div>
        <!-- Input per invaire un messaggio testuale -->
        <div class="input-group">
            <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
            <!-- Form per inviare un messaggio -->
            <form @submit.prevent="check">
                <input type="file" ref="file" accept=".jpg,.jpeg" @change="checkFile" />
                <input type="text" class="form-control" v-model="text" placeholder="Type your message here">
                <button type="submit" class="btn btn-outline-primary">Send</button>
            </form>            
        </div>
  
      <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
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