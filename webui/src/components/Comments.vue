<script>
export default {
    props: {
        show: Boolean,
        msg: Object,
        comments: Object,
    },
    data() {
        return {
            errorMsg: "",
            userId: sessionStorage.userId,
            convId: parseInt(this.$route.params.convId),
            emojis: ["😀", "😂", "😍", "😎", "😭", "😡", "🎉", "❤️", "👍", "🔥"],
        };
    },
    methods: {
        closeMod() {
            window.location.reload();
            this.$emit('close');
        },
        async commentMessage(emoji) {
            this.errorMsg = "";
            const url = `/users/${sessionStorage.userId}/conversation/${this.convId}/messages/${this.msg.Id}/reactions`;
            this.$axios.post(url, emoji, {headers: { 'Authorization': `${sessionStorage.token}`}}).then(() => {
                this.closeMod();
            }).catch(e => {this.errorMsg = e.toString()});
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
                        <button class="like-btn" @click="closeMod">
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
                    </div>
                    <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
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
    background-color: rgba(0, 0, 0, 0.5);
    display: table;
    transition: opacity 0.3s ease;
}

.modal-wrapper {
    display: table-cell;
    vertical-align: middle;
}

.modal-container {
    width: 350px;
    margin: 0px auto;
    background-color: #fff;
    border-radius: 2px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.33);
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
    color: #42b983;
}

.modal-header button {
    color: rgb(86, 86, 86);
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
    grid-template-columns: repeat(5, 1fr);
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

  