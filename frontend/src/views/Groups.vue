<template>
  <main>
    <h1>Groups</h1>

    <router-link :to="{ name: 'CreateGroup' }">
      <button class="create-group-button" @click="createNewGroup">
        Create a New Group
      </button></router-link
    >

    <div class="break"></div>

    <h2>My Groups:</h2>

    <div v-if="myGroups">
      <div
        v-for="group in myGroups"
        :key="group.id"
        class="group"
        @click="redirectToGroup(group.group_id)"
      >
        <div class="group-header">
          <div class="group-title">
            <h3>{{ group.title }}</h3>
          </div>
          <div class="group-creator">
            <img
              class="user-avatar"
              v-if="userAvatarUrl(group.user_id)"
              :src="userAvatarUrl(group.user_id)"
            />
            <router-link
              class="creator-link"
              :to="{ name: 'userid', params: { userID: group.user_id } }"
              @click.stop
            >
              {{ group.creator_name }}
            </router-link>
          </div>
          <p class="group-date">{{ formatGroupDateTime(group.created_at) }}</p>
        </div>

        <div class="group-description">
          <p>{{ group.content }}</p>
        </div>
      </div>
    </div>

    <div class="break"></div>

    <!-- Display all groups -->
    <h2>All Groups:</h2>
    <div v-if="allGroups">
      <div
        v-for="group in allGroups"
        :key="group.id"
        class="group"
        @click="redirectToGroup(group.group_id)"
      >
        <div class="group-header">
          <div class="group-title">
            <h3>{{ group.title }}</h3>
          </div>
          <div class="group-creator">
            <img
              class="user-avatar"
              v-if="userAvatarUrl(group.user_id)"
              :src="userAvatarUrl(group.user_id)"
            />
            <router-link
              class="creator-link"
              :to="{ name: 'userid', params: { userID: group.user_id } }"
              @click.stop
            >
              {{ group.creator_name }}
            </router-link>
          </div>
          <p class="group-date">{{ formatGroupDateTime(group.created_at) }}</p>
        </div>
        <div class="group-description">
          <p>{{ group.content }}</p>
        </div>
      </div>
    </div>
  </main>
</template>

<script>
import { onMounted } from "vue";
import getGroups from "@/composables/getGroups";
import getUsers from "@/composables/getUsers";

export default {
  name: "Groups",
  computed: {
    userAvatarUrl() {
      return (userId) => {
        const user = this.users.find((user) => user.user_id === userId);
        return user
          ? user.avatar
          : "http://localhost:8000/uploads/avatars/default-avatar-profile.jpg";
      };
    },
  },
  methods: {
    redirectToGroup(groupID) {
      this.$router.push({ name: "GroupID", params: { groupID: groupID } });
    },
  },
  setup() {
    const { myGroups, allGroups, loadGroups } = getGroups();
    const { users, loadUsers } = getUsers();

    onMounted(async () => {
      await loadGroups();
      await loadUsers();
    });

    const formatGroupDateTime = (createdAt) => {
      const groupDate = new Date(createdAt);
      const time = groupDate.toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });
      const day = ("0" + groupDate.getDate()).slice(-2);
      const month = ("0" + (groupDate.getMonth() + 1)).slice(-2);
      const year = groupDate.getFullYear();
      return `${time}, ${day}.${month}.${year}`;
    };

    return { myGroups, allGroups, formatGroupDateTime, users };
  },
};
</script>

<style scoped>
header {
  display: flex;
}

.create-group-button {
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 5px;
  padding: 7px 15px;
  cursor: pointer;
}

.create-group-button:hover {
  background-color: #309668;
}

.group {
  margin: 0 auto;
  margin-bottom: 20px;
  padding: 5px 20px;
  border-radius: 1rem;
  background-color: #f0f0f0;
  display: flex;
  flex-direction: column;
}

.group:hover {
  background-color: #e6e6e6;
  cursor: pointer;
}

.group-header {
  display: flex;
  border-bottom: 1px solid #bbb;
  justify-content: space-between;
  align-items: center;
}

.group-title {
  text-align: left;
  display: flex;
  margin-left: 20px;
}

.group-creator {
  display: flex;
}

.group-date {
  font-size: small;
  color: #888;
  margin-right: 10px;
  align-self: center;
}

.group-description {
  text-align: left;
  margin-left: 20px;
}

.creator-link {
  font-weight: bold;
  text-decoration: none;
  color: inherit;
  margin-right: 10px;
  align-self: center;
}

.creator-link:hover {
  text-decoration: underline;
  color: #888;
}

.break {
  margin: 40px 0;
}

.break2 {
  width: 90%;
  border-bottom: 1px solid #bbb;
}

.user-avatar {
  width: 30px;
  aspect-ratio: 1 / 1;
  clip-path: circle(50%);
  margin-right: 10px;
  margin-bottom: 5px;
}
</style>
