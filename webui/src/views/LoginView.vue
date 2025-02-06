<script>
export default {
    data() {
        return {
            username: "",
            errorMsg: "",
        }
    },
    emits: ['successful-login'],
    methods: {
        async doLogin() {
            try {
                if (this.username.length < 3 || this.username.length > 16) {
                    throw "Invalid username: must be between 3 and 16 characters."
                }
                let response = await this.$axios.post('/session', {
                    name: this.username
                }, {headers: {'Content-Type': 'application/json'}});

                sessionStorage.userId = response.data.id;
                sessionStorage.userName = response.data.name;
                sessionStorage.token = response.data.id;
                sessionStorage.userImage = response.data.image;

                this.$router.push("/home");
                this.$emit('successful-login');
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
    <ErrorMsg v-if="errorMsg" :msg="errorMsg" ></ErrorMsg>
    <div class="container">
        <form v-on:submit="doLogin">
            <h1>Welcome to WASA-Text!</h1>
            <input type="text" v-model="username" placeholder="Enter your username">
            <button type="submit">Login</button>
        </form>
    </div>
</template>

<style>
.container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 50vh;
}

.container form {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.container input {
    margin: 15px;
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 5px;
}

.container button {
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    background-color: blue;
    color: white;
    cursor: pointer;
}

.container button:hover {
    background-color: rgb(0, 102, 255);
}
</style>