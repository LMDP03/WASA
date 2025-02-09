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

                sessionStorage.userId = response.data.Id;
                sessionStorage.userName = response.data.Name;
                sessionStorage.token =  response.data.Id;
                sessionStorage.userImage = response.data.Id;

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
    }
}
</script>

<template>
    <ErrorMsg v-if="errorMsg != ''" :msg="errorMsg" ></ErrorMsg>
    <div class="login-container">
        <form @submit.prevent="doLogin">
            <h1>Welcome to WASA-Text!</h1>
            <input type="text" v-model="username" placeholder="Enter your username">
            <button type="submit">Login</button>
        </form>
    </div>
</template>

<style>
.login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 400%;
    width: 80%
}
.login-container form {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.login-container input {
  margin: 10px;
  padding: 7px;
  border: 1px solid black;
  border-radius: 20px;
  justify-items: center;
}
.login-container button {
  padding: 10px 20px;
  border-radius: 5px;
  background-color: black;
  color: white;
  cursor: pointer;
  border: 1px solid black
}
.login-container button:hover {
  background-color: white;
  color: black;
}
</style>