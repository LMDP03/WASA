<script>
export default {
    data() {
        return {
            username: "",
            errorMsg: "",
        }
    },
    emits: ['login-success'],
    methods: {
        async doLogin() {
        try {
            if (this.username.length < 3 || this.username.length > 16) throw "Invalid username, it must contains min 3 characters and max 16 characters"

            let response = await this.$axios.post('/session', {name: this.username});

            sessionStorage.userId = response.data.Id;
            sessionStorage.userName = response.data.Name;
            sessionStorage.token = response.data.Id;
            sessionStorage.userImage = response.data.Image;

            this.$router.push("/home");
            this.$emit('login-success');
        } catch (e) {
            this.errorMsg = e.toString();
            document.getElementsByTagName("input")[0].style.outline = "auto";
            document.getElementsByTagName("input")[0].style.outlineColor = "red";
        };
        }
    },
    mounted() {
        if (sessionStorage.token) {
            this.$router.push("/home");
            return;
        }
        sessionStorage.clear();
    },
}

</script>

<template>
    <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
    <div class="login-container">
        <form @submit.prevent="doLogin">
            <h1>WasaText</h1>
            <input type="text" v-model="username" placeholder="Enter your username" />
            <button type="submit">Login</button>
        </form>
    </div>
</template>


<style>
.login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 50vh;
}

.login-container form {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.login-container input {
    margin: 15px;
    padding: 10px;
    border: 1px solid white;
    border-radius: 5px;
}

.login-container button {
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    background-color: black;
    color: white;
    cursor: pointer;
}

.login-container button:hover {
    background-color: white;
    color: black;
}
</style>
