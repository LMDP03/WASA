<script>
import ModalUsers from '../components/ModalUsers.vue';
import Group from '../components/GroupSettings.vue';
export default {
	data: function() {
		return {
			errorMsg: "",
			searchModalIsVisible: false,
			createGroupModalIsVisible: false,
			some_data: [],
			users: [],
		}
	},
	emits: ['successful-login', 'username-changed'],
	methods: {
		async getMyConversations() {
            this.errorMsg = "";
            try {
                const url = `users/${sessionStorage.userId}/conversations?srcName=${this.searchName}`
                let response = await this.$axios.get(url, { headers: { 'Authorization': `${sessionStorage.token}` } });
                if (response.data == null) {
                    return
                }
                this.some_data = response.data;
            } catch (e) {
                this.errorMsg = e.toString();
            }
        },
		async getConversation(convid) {
			localStorage.clear();
			errorMsg = "";
			try {
				const url = `users/${sessionStorage.userId}/conversation/${convid}`;
				let response = await this.$axios.get(url, { headers: { 'Authorization': `${sessionStorage.token}` } });
				if (response.data == null) {
					return;
				}
				localStorage.userId = response.id;
				localStorage.userName= response.name;
				localStorage.userImage = response.image;
				localStorage.users = JSON.stringify(response.participants);
				localStorage.groupFlag = response.group;
				this.$router.push(`/conversation/${response.id}`);
			} catch (e) {
				this.errorMsg = e.toString();
			}
		},
		handleSearchModalToggle() {
			this.searchModalIsVisible = !this.searchModalIsVisible;
		},
		// Funzione per mostare o nascondere il modale per la creazione di un nuovo gruppo
		handleCreateGroupModalToggle() {
			this.createGroupModalIsVisible = !this.createGroupModalIsVisible;
		},
	},
	mounted() {
		if (!sessionStorage.token) {
			this.$router.push("/");
			return;
		}
	},
	components: {ModalUsers, Group}
}
</script>

<template>
	<div>
	  <div
		class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
		<h1 class="h2">Home page</h1>
  
		<!-- Modale utilizzato per la creazione di un nuovo gruppo -->
		<Group :show="createGroupModalIsVisible" @close="handleCreateGroupModalToggle" title="search">
		  <template v-slot:header>
			<h3>Select users</h3>
		  </template>
		</Group>
		<!-- Modale utilzzato per la ricerca degli utenti con cui aprire una nuova conversazione -->
		<ModalUsers :show="searchModalIsVisible" @close="handleSearchModalToggle" title="search">
		  <template v-slot:header>
			<h3>Users</h3>
		  </template>
		</ModalUsers>
  
		<!-- Pulsanti per aggiornare la lista delle conversazioni, creare un nuovo gruppo e cercare nuovi utenti -->
		<div class="btn-toolbar mb-2 mb-md-0">
		  <div class="btn-group me-2">
			<!-- Pulsante per aggiornare la lista delle conversazioni -->
			<button type="button" class="btn btn-sm btn-outline-secondary" @click="getMyConversations">
			  Refresh
			</button>
		  </div>
		  <!-- Pulsante per creare un nuovo gruppo -->
		  <div class="btn-group me-2">
			<button type="button" class="btn btn-sm btn-outline-primary" @click="handleCreateGroupModalToggle">
			  New Group
			</button>
		  </div>
		  <!-- Pulsante per cercare nuovi utenti e aprire un nuova conversazione -->
		  <div class="btn-group me-2">
			<button type="button" class="btn btn-sm btn-outline-primary" @click="handleSearchModalToggle">
			  New Chat
			</button>
		  </div>
		</div>
	  </div>
  
	  <!-- Lista delle conversazioni -->
	  <div v-if="some_data.length != 0">
		<!-- Mostra le conversazioni dell'utente iterando all'interno di some_data dove sono salvate tutte le conversazioni -->
		<div class="conversations" v-for="response in some_data" :key="response.id">
			<!-- Controlla se la conversazione non è con un gruppo -->
			<div v-if="response.group == false">
				<!-- Mostra il nome dell'utente con cui si sta conversando, l'ultimo messaggio e chi lo ha inviato -->
				<!-- Se il messaggio è una foto, mostra "Photo" al posto del testo -->
				<button v-if="response.image== ''" type="button" class="btn btn-sm btn-outline-primary" @click="getConversation(response)">
					{{ response.name }} <br> {{ response.lastmessage.sender.name }}: {{ response.message.text }}
				</button>
					<!-- Altrimenti mostra il contenuto del messaggio -->
				<button type="button" class="btn btn-sm btn-outline-primary" @click="getConversation(response)" v-else>
					{{ response.user.name }} <br> {{ response.lastmessage.sender.name }}: Photo
				</button>
			</div>
		  <hr>
		</div>
	  </div>
  
	  <!-- Se non ci sono conversazioni, mostra un messaggio -->
	  <div v-else>
		<p>Start Chatting!</p>
	  </div>
  
	  <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
	</div>
  </template>
  
  <style></style>

<style>
.profile-picture {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
  object-fit: cover;
}

.username {
  font-size: 14px;
  font-weight: bold;
}
</style>
