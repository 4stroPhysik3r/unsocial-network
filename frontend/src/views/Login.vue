<template>
  <main>
    <h1>Login</h1>

    <h2>Enter your credentials</h2>
    <div class="post">
      <form @submit.prevent="login">
        <div class="email">
          <label>Email :</label>
        </div>
        <div class="emailInput">
          <input type="email" v-model="email" required />
        </div>
        <div class="pwd">
          <label>Password :</label>
        </div>
        <div class="pwdInput">
          <input type="password" v-model="password" required />
        </div>
        <button class="login-btn" type="submit">Login</button>
      </form>
    </div>
  </main>
</template>

<script>
import { useAuth } from "@/composables/useAuth";

export default {
  name: "Login",
  data() {
    return {
      email: "",
      password: "",
    };
  },
  methods: {
    async login() {
      const userData = {
        email: this.email,
        password: this.password,
      };

      try {
        const response = await fetch("http://localhost:8000/api/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify(userData),
        });

        const data = await response.json();
        if (!response.ok) {
          throw new Error(data.error || "An error occurred");
        }

        const { checkAuthStatus } = useAuth();
        await checkAuthStatus(); // This will update isAuthenticated based on the new session

        this.$router.push("/"); // Redirect to the homepage
      } catch (error) {
        alert(error.message);
      }
    },
  },
};
</script>

<style scoped>
form {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(3, 1fr);
  grid-gap: 10px 40px;
  align-items: center;
  text-align: right;
}

.email {
  grid-area: 1 / 1 / 2 / 2;
}
.emailInput {
  grid-area: 1 / 2 / 2 / 4;
}
.pwd {
  grid-area: 2 / 1 / 3 / 2;
}
.pwdInput {
  grid-area: 2 / 2 / 3 / 4;
}

label {
  font-weight: bold;
  text-align: left;
}

input {
  display: flex;
  padding: 10px;
  width: 90%;
  border-radius: 5px;
  border: 1px solid #ccc;
}

button {
  grid-area: 3 / 2 / 4 / 3;
}

.login-btn {
  padding: 10px 25px;
}

.post {
  max-width: 800px;
  padding: 20px;
}

.post:hover {
  background-color: #f0f0f0;
  cursor: default;
}

.error {
  color: red;
}
</style>
