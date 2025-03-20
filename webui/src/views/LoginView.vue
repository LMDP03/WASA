<script>
export default {
    data() {
        return {
            userName: "",
            errorMsg: "",
            usernameValidate: new RegExp('^\\w{0,16}$'),
        }
    },
    methods: {
        async doLogin() {
            try {
                if (this.userName.length < 3 || this.userName.length > 16 || !this.usernameValidate.test(this.userName)) throw "Username must be between 3 and 16 alphanumerical characters";
                let response = await this.$axios.post('/session', {
                    name: this.userName,
                });
                localStorage.userId = response.data.Id;
                localStorage.userName = response.data.Name;
                localStorage.userImage = response.data.Image;
                localStorage.token = response.data.Id;
                this.$router.push('/home');
                this.$emit('login-success');
            } catch (e) {
                this.errorMsg = e.toString();
                document.getElementsByTagName("input")[0].style.outline = "auto";
                document.getElementsByTagName("input")[0].style.outlineColor = "red";
            };
        }
    },
    emits: ['login-success'],
    mounted() {
        if (localStorage.token) {
            this.$router.push('/home');
            return;
        }
        localStorage.clear();
    },
}
</script>

<template>
    <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
    <div class="login-container">
        <form @submit.prevent="doLogin">
            <h1>WasaText</h1>
            <input type="text" v-model="userName" placeholder="Enter your username" />
            <button type="submit">Login</button>
        </form>
    </div>
</template>

<style>
.login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 75vh;
    width: 50vw;
}

.login-container form {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.login-container input{
    margin: 15px;
    padding: 10px;
    border: 1px solid white;
    border-radius: 5px;
}

.login-container button {
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    background-color: blue;
    color: white;
    cursor: pointer;
}

.login-container button:hover {
    background-color: darkblue;
}
</style>