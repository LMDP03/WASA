<script>
export default {
    props: {
        show: Boolean,
        conversations: Array,
        title: String,
    },
    data() {
        return {
            errorMsg: "",
            searchText: "",
            usernameValidate: new RegExp('^\\w{0,16}$'),
            filteredConvs: [],
        };
    },
    methods: {
        closeModal() {
            this.searchText = "";
            this.$emit('close');
        },
        async filterConvs() {
            this.errorMsg = "";
            this.filteredConvs = this.conversations;
            if (this.searchText.length > 0) {
                if (this.searchText.length > 16 || !this.usernameValidate.test(this.searchText)) {
                    this.errorMsg = "Invalid username, it can contain at most 16 alphanumerical characters.";
                    this.filteredConvs = [];
                    return;
                }

                if (this.title === "search") {
                    try {
                        const url = `/users/${localStorage.userId}/conversations?srcName=${this.searchText}`;
                        let response = await this.$axios.get(url, {headers: {'Authorization': localStorage.token }});
                        if (response.data == null) {
                            this.filteredConvs = [];
                            return;
                        }
                        this.filteredConvs = response.data;
                    } catch (e) {
                        this.errorMsg = e.toString();
                        this.filteredConvs = [];
                    }
                } else {
                    this.filteredConvs = this.conversations.filter(conv => conv.Name.toLowerCase().includes(this.searchText.toLowerCase()));
                }
            }
        },
        async selectConv(conv) {
            sessionStorage.clear();
            sessionStorage.userId = 0;
            sessionStorage.convId = conv.Id;
            sessionStorage.convImg = conv.Image;
            sessionStorage.convName = conv.Name;
            sessionStorage.isGroup = conv.Group;
            this.$router.push('/conversation');
            this.closeModal();
        },
    },
    watch: {
        searchText() {
            this.filterConvs();
        },
        show() {
            this.filteredConvs = this.conversations;
        }
    },
};
</script>

<template>
    <Transition name="modal">
        <div v-if="show" class="modal-mask">
            <div class="modal-wrapper">
                <div class="modal-container">
                    <div class="modal-header">
                        <slot name="header">deafult header</slot>
                        <button class="like-btn" @click="closeModal">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#x" />
                            </svg>
                        </button>
                    </div>

                    <div class="modal-body">
                        <slot name="body">
                            <div class="search-input">
                                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                                <input type="text" v-model="searchText" placeholder="Search" />
                            </div>
                            <div class="search-results">
                                <div v-for="conv in filteredConvs" :key="conv.Id" @click="selectConv(conv)">
                                    <div class="user">
                                        <p>{{ conv.Name }}</p>
                                    </div>
                                </div>
                            </div>
                        </slot>
                    </div>
                </div>
            </div>
        </div>
    </Transition>
</template>

<style>
.custom-link{
    color: inherit;
    text-decoration: none;
}
.modal-mask {
    position: fixed;
    z-index: 9998;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: table;
    transition: opacity 0.3 ease;
}
.modal-wrapper {
    display: table-cell;
    vertical-align: middle;
}

.modal-container {
    width: 350px;
    margin: 0px auto;
    background-color: white;
    border-radius: 2px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
    transition: all 0.3s ease;
}

.modal-header {
    height: 70px;
    padding: 20px 15px 10px 15px;
}

.modal-header h3 {
    margin-top: 0;
    font-size: 25px;
    color: black;
}

.modal-header button {
    color: gray;
    background: none;
    border: none;
    padding: 5px;
    line-height: 12px;
    font-size: 15px;
}

.modal-header button svg {
    width: 20px;
    height: 20px;
}

.search-input {
    padding: 0 15px;
}

.search-input input {
    height: 30px;
    width: 100%;
    outline: none;
    border-radius: 3px;
    border: 1px solid gray;
}

.search-results {
    font-size: 15px;
    padding: 10px 15px;
    border-bottom: 1px solid white;
    cursor: pointer;
    max-height: 200px;
    overflow-y: scroll;
}

.modal-deafult-button {
    float: right;
}

.username-form {
    display: flex;
    flex-direction: column;
    padding: 0 15px;
}

.username-form input {
    margin-bottom: 10px;
    margin-top: 5px;
    outline: none;
    border-radius: 3px;
    border: 1px solid gray;
}

.username-form button {
    margin-bottom: 15px;
}
</style>