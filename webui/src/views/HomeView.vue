<script>
import Private from '../components/SearchUsers.vue';
import Group from '../components/CreateGroup.vue';
export default {
	data: function() {
		return {
			errorMsg: "",
			showUserSearh: false,
			showGroupCreate: false,
			some_data: [],
			users: [],
		}
	},
	emits: ['successful-login', 'username-changed'],
	methods: {
		async getMyConversations() {
            this.errorMsg = "";
            try {
                const url = `users/${sessionStorage.userId}/conversations`
                let response = await this.$axios.get(url, { headers: { 'Authorization': `${sessionStorage.token}` } });
                this.some_data = response.data;
            } catch (e) {
                this.errorMsg = e.toString();
            }
        },
		async goToConversation(conv) {
			localStorage.clear();
			localStorage.userId = conv.Id;
			localStorage.userName = conv.Name;
			localStorage.userImage = conv.Image;
			localStorage.isGroup = conv.Group;
			this.$router.push(`/conversation/${conv.Id}`);
		},
		handleSearchMod() {
			this.showUserSearh = !this.showUserSearh;
		},
		// Funzione per mostare o nascondere il modale per la creazione di un nuovo gruppo
		handleGroupMod() {
			this.showGroupCreate = !this.showGroupCreate;
		}
	},
	mounted() {
		if (!sessionStorage.token) {
			this.$router.push("/");
			return;
		}
	},
	components: {Private, Group}
}
</script>

<template>
	<div>
	  	<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>
	
			<!-- Modale utilizzato per la creazione di un nuovo gruppo -->
			<Group :show="showGroupCreate" @close="handleGroupMod" title="search">
				<template v-slot:header>
					<h3>Select users</h3>
				</template>
			</Group>
			<!-- Modale utilzzato per la ricerca degli utenti con cui aprire una nuova conversazione -->
			<Private :show="searchModalIsVisible" @close="handleSearchMod" title="search">
				<template v-slot:header>
					<h3>Select User</h3>
				</template>
			</Private>
	
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
					<button type="button" class="btn btn-sm btn-outline-primary" @click="handleGroupMod">
						New Group
					</button>
				</div>
				<!-- Pulsante per cercare nuovi utenti e aprire un nuova conversazione -->
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="handleSearchMod">
						Start Conversation
					</button>
				</div>
			</div>
		</div>
	
		<!-- Lista delle conversazioni -->
		<div v-if="some_data.length !== 0">
			<!-- Mostra le conversazioni dell'utente iterando all'interno di some_data dove sono salvate tutte le conversazioni -->
			<div class="conversations" v-for="response in some_data" :key="response.Id">
				<!-- Mostra il nome dell'utente con cui si sta conversando, l'ultimo messaggio e chi lo ha inviato -->
				<!-- Se il messaggio è un testo, mostra il contenuto -->
				<button v-if="response.LastMessage.Image== ''" type="button" class="btn btn-sm btn-outline-primary" @click="goToConversation(response)">
					{{ response.Image }} {{ response.Name }} <br> {{ response.LastMessage.Sender.Name }}: {{ response.LastMessage.Text }}
				</button>
					<!-- Altrimenti mostra "Photo" -->
				<button v-else type="button" class="btn btn-sm btn-outline-primary" @click="goToConversation(response)">
					{{ response.Image }} {{ response.Name }} <br> {{ response.LastMessage.Sender.Name }}: Photo
				</button>
			</div>
		</div>
	
		<!-- Se non ci sono conversazioni, mostra un messaggio -->
		<div v-else>
			<p>Start Chatting!</p>
		</div>

	  	<ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
	</div>
</template>
  

