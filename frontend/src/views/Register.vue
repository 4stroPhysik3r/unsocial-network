<template>
  <main>
    <h1>Register</h1>

    <h2>Enter your credentials</h2>
    <div class="post">
      <form @submit.prevent="register">
        <div class="input">
          <label for="email">Email*: </label>
          <input type="email" v-model="email" required />
        </div>
        <!-- extra div, so form has gaps -->
        <div class="input"></div>
        <div class="input">
          <label>Password*: </label>
          <input type="password" v-model="password" required />
        </div>
        <div class="input">
          <label>Confirm Password*: </label>
          <input type="password" v-model="confirmPassword" required />
          <p v-if="passwordError" class="error">{{ passwordError }}</p>
        </div>
        <div class="input">
          <label>First Name*: </label>
          <input type="text" v-model="firstName" required />
        </div>
        <div class="input">
          <label>Last Name*: </label>
          <input type="text" v-model="lastName" required />
        </div>
        <div class="input">
          <label>Date of Birth*: </label>
          <input type="date" v-model="dob" required />
        </div>
        <!-- extra div, so form has gaps -->
        <div class="input"></div>
        <div class="input">
          <label>Avatar/Image: </label>
          <input type="file" @change="handleFileInputChange" />
        </div>
        <div class="input">
          <label>Nickname: </label>
          <input type="text" v-model="nickname" />
        </div>
        <div class="input">
          <label>About Me: </label>
          <textarea v-model="aboutMe"></textarea>
        </div>
        <!-- extra div, so form has gaps -->
        <div class="input"></div>
        <div class="input">
          <label>Public Profile: </label>
          <input
            type="checkbox"
            class="profilePublic"
            v-model="profilePublic"
          />
        </div>
        <!-- Display error message if email is already taken -->
        <p v-if="emailError" class="error">{{ emailError }}</p>
        <!-- Display general error message for required fields -->
        <p v-if="generalError" class="error">{{ generalError }}</p>
        <p class="required">*required</p>
        <button type="submit">Register</button>
      </form>
    </div>
  </main>
</template>

<script>
export default {
  name: "Register",
  data() {
    return {
      email: "",
      password: "",
      confirmPassword: "",
      firstName: "",
      lastName: "",
      dob: "",
      avatar: null,
      nickname: "",
      aboutMe: "",
      emailError: "",
      passwordError: "",
      generalError: "",
      profilePublic: false,
    };
  },
  methods: {
    handleFileInputChange(event) {
      const file = event.target.files[0];
      if (!file) {
        return;
      }
      const reader = new FileReader();
      reader.onload = (e) => {
        // Set the base64 string to the avatar
        this.avatar = e.target.result;
      };
      reader.readAsDataURL(file);
    },
    validatePassword() {
      if (this.password !== this.confirmPassword) {
        this.passwordError = "Passwords do not match";
      } else if (this.password.length < 5) {
        this.passwordError = "Password must be at least 5 characters long";
      } else if (this.password === "" || this.confirmPassword === "") {
        this.passwordError = "Password fields cannot be empty";
      } else {
        this.passwordError = "";
      }
    },
    validateEmail() {
      // Check if email field is not empty
      if (!this.email.trim()) {
        this.emailError = "Email field cannot be empty";
        return false; // Exit the validation process
      }

      // Regular expression for email validation
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

      // Check if email matches the format
      if (!emailRegex.test(this.email)) {
        this.emailError = "Invalid email format";
        return false;
      }

      // Clear the error message if email is valid
      this.emailError = "";
      return true;
    },
    register() {
      // Validate email format
      const isEmailValid = this.validateEmail();

      // Validate password match
      this.validatePassword();

      // Check if all required fields are filled
      const areRequiredFieldsFilled =
        this.email &&
        this.password &&
        this.confirmPassword &&
        this.firstName &&
        this.lastName &&
        this.dob;

      // Check if password and confirm password match
      if (this.password !== this.confirmPassword) {
        this.passwordError = "Passwords do not match";
        return;
      }

      // Display error message for required fields if validation fails
      if (!isEmailValid || !areRequiredFieldsFilled) {
        this.generalError = "Required fields cannot be empty";
        return;
      } else {
        // Clear general error message if all required fields are filled
        this.generalError = "";
      }

      // Handle form submission here
      const userData = {
        email: this.email,
        password: this.password,
        confirmPassword: this.confirmPassword,
        firstName: this.firstName,
        lastName: this.lastName,
        dateOfBirth: this.dob,
        nickname: this.nickname,
        aboutMe: this.aboutMe,
        profilePublic: this.profilePublic,
      };

      // Add the avatar to the userData if it exists
      if (this.avatar) {
        userData.avatar = this.avatar;
      }

      // Make an HTTP POST request to your backend API
      fetch("http://localhost:8000/api/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(userData),
      })
        .then(async (response) => {
          if (!response.ok) {
            // Check if the error message indicates that the email is already taken
            const errorMessage = await response.text();
            if (errorMessage.includes("email is already taken")) {
              this.emailError = "Email is already taken";
              throw new Error("Email is already taken");
            } else {
              throw new Error(
                "Failed to register user: " + response.statusText
              );
            }
          }
          // Handle successful registration
          alert("Registration successful");
          this.$router.push("/"); // Redirect to the homepage
        })
        .catch((error) => {
          console.error("Error registering user:", error);
        });
    },
  },
};
</script>

<style scoped>
form {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-gap: 10px 40px;
  padding: 20px;
}

.input {
  display: flex;
  justify-content: space-between;
}

label {
  margin-top: 10px;
  text-align: left;
  font-weight: bold;
}

button {
  grid-column: span 2;
}

input {
  border-radius: 5px;
  border: 1px solid #ccc;
}

input[type="file"] {
  align-self: center;
}

.post {
  width: 80vw;
  max-width: 1200px;
}

.post:hover {
  background-color: #f0f0f0;
  cursor: default;
}

.required {
  color: red;
  font-style: italic;
  font-size: small;
}

.error {
  color: red;
}
</style>
