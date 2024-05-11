<template>
  <main>
    <div id="userSearch" v-show="isAuthenticated">
      <input
        id="myInput"
        type="text"
        placeholder="Search users.."
        @keyup="searchFunction"
        @focus="searchFunction"
      />
      <div id="options">
        <div v-for="user in sortedUsers" :key="user.user_id">
          <a :href="'/userid/' + user.user_id" @click="clearSearchInput">
            {{ user.firstname }} {{ user.lastname }}</a
          >
        </div>
      </div>
    </div>
    <header class="nav">
      <p>unsocial worknet</p>
      <!-- Conditional Rendering Based on isAuthenticated -->
      <div>
        <router-link class="nav-link" to="/" v-show="isAuthenticated"
          >Home</router-link
        >
        <router-link
          class="nav-link"
          :to="{ name: 'Groups' }"
          v-show="isAuthenticated"
          >Groups</router-link
        >
        <router-link
          class="nav-link"
          :to="{ name: 'Chats' }"
          v-show="isAuthenticated"
          >Chats
          <span
            v-if="unreadMessages && unreadMessages.length > 0"
            class="red-dot-header"
          >
          </span
        ></router-link>
        <router-link
          class="nav-link"
          :to="{ name: 'Profile' }"
          v-show="isAuthenticated"
          >Profile</router-link
        >
        <router-link
          class="nav-link"
          :to="{ name: 'Login' }"
          v-show="!isAuthenticated"
          >Login</router-link
        >
        <router-link
          class="nav-link"
          :to="{ name: 'Register' }"
          v-show="!isAuthenticated"
          >Register</router-link
        >
      </div>
    </header>
    <div id="notifications" v-show="isAuthenticated">
      <router-link class="nav-link" :to="{ name: 'Notifications' }"
        >Notifications
        <span v-if="notifications.length" class="red-dot-header"></span>
      </router-link>
    </div>
    <router-view />
    <footer>
      <p>
        ©
        <a href="https://01.kood.tech/git/4stroPhysik3r" target="_blank"
          >4stroPhysik3r</a
        >, <a href="https://01.kood.tech/git/Freyby" target="_blank">Freyby</a>,
        <a href="https://01.kood.tech/git/KristjanM/" target="_blank"
          >KristjanM</a
        >
      </p>
    </footer>
  </main>
</template>

<script>
import { computed, onMounted, watch, watchEffect } from "vue";
import { useAuth } from "@/composables/useAuth";
import { useRoute } from "vue-router";
import getUsers from "@/composables/getUsers";
import searchUsers from "@/composables/searchUsers";
import { useChatNotifications } from "@/composables/useChatNotifications";
import { useNotifications } from "@/composables/useNotifications";

export default {
  setup() {
    const { isAuthenticated, checkAuthStatus } = useAuth();
    const route = useRoute();
    const { users, loadUsers } = getUsers();
    const { searchFunction, clearSearchInput } = searchUsers();
    const { unreadMessages } = useChatNotifications(
      "ws://localhost:8000/api/chat-notifications/ws"
    );
    const { connect, notifications } = useNotifications(false);

    onMounted(async () => {
      await checkAuthStatus();
      await loadUsers();
    });

    watchEffect(() => {
      if (isAuthenticated.value) {
        connect();
      }
    });

    watch(
      () => route.path,
      async () => {
        await loadUsers();
      }
    );

    const sortedUsers = computed(() => {
      return [...users.value].sort((a, b) => {
        const nameA = `${a.firstname} ${a.lastname}`.toUpperCase();
        const nameB = `${b.firstname} ${b.lastname}`.toUpperCase();
        if (nameA < nameB) return -1;
        if (nameA > nameB) return 1;
        return 0;
      });
    });

    return {
      isAuthenticated,
      searchFunction,
      clearSearchInput,
      sortedUsers,
      unreadMessages,
      notifications,
    };
  },
};
</script>

<style>
body {
  margin: 0;
}
main {
  max-width: 1400px;
  margin: 0 auto;
}
.nav {
  padding: 30px;
  border-bottom: 3px solid #f0f0f0;
}

.nav .nav-link {
  font-weight: bold;
  color: #2c3e50;
  text-decoration: none;
  margin: 10px;
}

.nav .nav-link:hover:not(.router-link-exact-active) {
  color: #aeaeae;
  text-decoration: underline 2px;
}

.router-link-exact-active {
  color: #42b983 !important;
}

#userSearch {
  position: absolute;
  top: 65px;
  left: 4%;
  padding-bottom: 10px;
  overflow: auto;
  z-index: 1;
}

#notifications {
  position: absolute;
  top: 80px;
  right: 4%;
  font-weight: bold;
  text-decoration: none;
  color: #2c3e50;
  padding-bottom: 10px;
  overflow: auto;
  z-index: 1;
}

#notifications a {
  text-decoration: none;
  color: #2c3e50;
}

#notifications a:hover {
  color: #aeaeae;
  text-decoration: underline 2px;
}

#userSearch a {
  display: block;
  color: inherit;
  text-decoration: none;
}

#myInput {
  padding: 5px 0;
  margin-top: 10px;
  border: 1px solid #aeaeae;
  border-radius: 0.5rem;
  background-color: #f6f6f6;
}

#myInput::placeholder {
  color: #888;
  font-style: italic;
  text-align: center;
}

#options {
  display: none;
  padding: 5px 15px;
  margin: auto;
  border-radius: 0 0 1rem 1rem;
  box-shadow: #aeaeae 0px 0px 10px;
  background-color: #f6f6f6;
  text-align: left;
}

#options a {
  position: relative;
  padding: 5px 0;
  z-index: 2;
}

#options a:hover {
  text-decoration: underline;
  color: #309668;
}

h1 {
  margin-top: 2rem;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin: 4rem;
  margin-top: 0;
}

.post {
  margin: 0 auto;
  margin-bottom: 20px;
  padding: 10px 20px 5px 20px;
  border-radius: 1rem;
  background-color: #f0f0f0;
  display: flex;
  flex-direction: column;
  align-items: center;
}

button {
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 16px;
  padding: 7px 15px;
  margin: auto;
  color: #fff;
  background-color: #42b983;
}

button:hover {
  background-color: #309668;
}

.break {
  border-bottom: 2px solid #f0f0f0;
  margin: 0 2rem;
}

footer {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 40px;
  text-align: center;
  background-color: #aeaeae;
  color: #000000;
  font-size: small;
}

footer a {
  text-decoration: none;
  color: #000000;
}

footer a:hover {
  text-decoration: underline;
}

.red-dot-header {
  width: 10px;
  height: 10px;
  background-color: #c92a2a;
  border-radius: 50%;
  display: inline-block;
  margin-left: auto;
  margin-right: 5px;
  margin-bottom: 5px;
}
</style>
