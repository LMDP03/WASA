<script>
export default {
	data: function() {
		return {
			errorMsg: "",
			showSearh: false,
			showGroup: false,
			conversations: [],
		}
	},
	emits: ['successful-login'],
	methods: {
		async getMyConversations() {
            this.errorMsg = "";
            try {
                const url = `users/${sessionStorage.userId}/conversations`
                let response = await this.$axios.get(url, { headers: { 'Authorization': `${sessionStorage.token}` } });
                this.conversations = response.data;
				sessionStorage.conversations = response.data;
            } catch (e) {
                this.errorMsg = e.toString();
            }
        },
		async goToConversation(conv) {
			localStorage.clear();
			localStorage.convId = conv.Id;
			localStorage.convName = conv.Name;
			localStorage.convImage = conv.Image;
			localStorage.isGroup = conv.Group;
			localStorage.receiverId = 0;
			this.$router.push(`/conversation`);
		}
	},
	mounted() {
		if (!sessionStorage.token) {
			this.$router.push("/");
			return;
		}
		this.getMyConversations();
	}
}
</script>

<template>
	<div>
	  	<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">
				Welcome!
			</h1>

			<!-- Pulsanti per aggiornare la lista delle conversazioni, creare un nuovo gruppo e cercare nuovi utenti -->
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<!-- Pulsante per aggiornare la lista delle conversazioni -->
					<button type="button" class="btn btn-sm btn-outline-primary" @click="getMyConversations">
						Sync
					</button>
				</div>
			</div>

		</div>
	
		<!-- Lista delle conversazioni -->
		<div v-if="conversations.length !== 0">

			<div class="conversations" v-for="preview in conversations" :key="preview.Id">
				<!-- Mostra il nome dell'utente con cui si sta conversando, l'ultimo messaggio e chi lo ha inviato -->
				<!-- Se il messaggio è un testo, mostra il contenuto -->
				<button v-if="preview.LastMessage.Image == ''" type="button" class="btn btn-sm btn-outline-primary" @click="goToConversation(preview)">
					{{ preview.Image }} {{ preview.Name }} <br> {{ preview.LastMessage.Sender.Name }}: {{ preview.LastMessage.Text }}
				</button>
					<!-- Altrimenti mostra "Photo" -->
				<button v-else type="button" class="btn btn-sm btn-outline-primary" @click="goToConversation(preview)">
					{{ preview.Image }} {{ preview.Name }} <br> {{ preview.LastMessage.Sender.Name }}:
					<svg class="feather">
						<use href="/feather-sprite-v4.29.0.svg#image" />
					</svg>
					{{ preview.LastMessage.Text }}
				</button>
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
  

