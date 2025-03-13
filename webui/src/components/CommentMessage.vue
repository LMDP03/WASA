<script>
    export default {
        props: {
            show: Boolean,
            msg: Object,
        },
        data() {
            return {
                userId: localStorage.userId,
                convId: sessionStorage.convId,
                emojis: ["😀", "😂", "😍", "😎", "😭", "😡", "🎉", "❤️", "👍", "🔥"],
                errorMsg: "",
            };
        },
        methods: {
            closeModal() {
                window.location.reload();
                this.$emit('close');
            },
            async commentMessage(emoji) {
                this.errorMsg = "";
                const url = `users/${this.userId}/conversation/${this.convId}/messages/${this.msg.MsgId}/reactions`;
                try {
                    let response = await this.$axios.post(url, emoji, { headers: { 'Authorization': localStorage.token } });
                    this.msg.Reactions = response.data;
                } catch (e) {
                    this.errorMsg = e.toString();
                }
            },
        },
    };
</script>

<template>
    <Transition name="modal">
        <div v-if="show" class="modal-mask">
            <div class="modal-wrapper">
                <div class="modal-container">
                    <div class="modal-header">
                        <h3>Choose a Reaction</h3>
                        <button class="like-btn" @click="closeModal">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#x" />
                            </svg>
                        </button>
                    </div>
                    <div class="modal-body">
                        <div class="emoji-grid">
                            <div v-for="emoji in emojis" :key="emoji" class="emoji" @click="commentMessage(emoji)">
                                {{ emoji }}
                            </div>
                        </div>
                        <ErrorMsg v-if="errorMsg" :msg="errormsg"></ErrorMsg>
                    </div>
                </div>
            </div>
        </div>
    </Transition>
</template>

<style>
    .modal-mask {
        position: fixed;
        z-index: 9998;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: black;
        display: table;
        transition: opeacity 0.3s ease;
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
        box-shadow: 0 2px 8px black;
        transition: all 0.3s ease;
    }

    .modal-header {
        height: 70px;
        padding: 20px 15px 10px 15px;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .modal-header h3 {
        margin: 0;
        font-size: 20px;
        color: lightgreen;
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

    .modal-body {
        padding: 15px;
        text-align: center;
    }

    .emoji-grid {
        display: grid;
        grid-template-columns: repeat(5, lfr);
        gap: 10px;
        justify-items: center;
        align-items: center;
    }

    .emoji {
        font-size: 24px;
        cursor: pointer;
        transition: transform 0.2s;
    }

    .emoji:hover {
        transform: scale(1.2);
    }
</style>