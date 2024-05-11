<template>
  <main>
    <h1>Create New Group</h1>
    <div class="createPost">
      <form @submit.prevent="createGroup">
        <div>
          <label for="groupTitle">Group Name:</label>
          <input type="text" id="groupTitle" v-model="groupTitle" required />
        </div>
        <div>
          <label for="groupDescription">Group Description:</label>
          <textarea
            id="groupDescription"
            v-model="groupContent"
            required
          ></textarea>
        </div>
        <!-- Search field -->
        <div id="memberSearch">
          <input
            id="myInput2"
            type="text"
            placeholder="Search to select Users.."
            @keyup="searchFunction"
            @focus="searchFunction"
          />
          <!-- Options -->
          <div id="options2">
            <div v-for="userSelection in users" :key="userSelection.user_id">
              <div v-if="userSelection.user_id !== user.userID">
                <p @click="addMember(userSelection)">
                  {{ userSelection.firstname }} {{ userSelection.lastname }}
                </p>
              </div>
            </div>
          </div>
        </div>
        <h3>Selected Members:</h3>
        <!-- Member list -->
        <div class="selectedMembers">
          <div v-for="member in selectedMembers" :key="member.user_id">
            {{ member.firstname }} {{ member.lastname }}
          </div>
        </div>
        <button type="submit" @click="clearSearchInput">Create Group</button>
      </form>
    </div>
  </main>
</template>

<script>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import getFollowingList from "@/composables/getFollowing";
import searchMembers from "@/composables/searchMembers";
import getUsers from "@/composables/getUsers";
import getUserFromSession from "@/composables/getUserFromSession";

export default {
  name: "CreateGroup",
  setup() {
    const router = useRouter();
    const { following, loadFollowing } = getFollowingList();
    const { users, loadUsers } = getUsers();
    const { searchFunction, clearSearchInput } = searchMembers();
    const { user, fetchUserDataFromSession } = getUserFromSession();

    onMounted(async () => {
      await loadFollowing();
    });

    onMounted(() => {
      loadUsers();
      fetchUserDataFromSession();
    });

    const groupTitle = ref("");
    const groupContent = ref("");
    const selectedMembers = ref([]);

    const createGroup = async () => {
      try {
        const memberIDs = selectedMembers.value.map((member) => member.user_id);

        const response = await fetch(
          "http://localhost:8000/api/create-group/",
          {
            credentials: "include",
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              Title: groupTitle.value,
              Content: groupContent.value,
              Members: memberIDs,
            }),
          }
        );

        console.log(memberIDs);
        if (!response.ok) {
          throw new Error("Failed to create group");
        }
        // Redirect to groups.vue page after successful creation
        router.push("/groups");
      } catch (error) {
        console.error("Error creating group:", error);
      }
    };

    const addMember = (user) => {
      if (
        user &&
        !selectedMembers.value.find((member) => member.user_id === user.user_id)
      ) {
        selectedMembers.value.push(user);
      }
    };

    return {
      users,
      groupTitle,
      groupContent,
      following,
      selectedMembers,
      createGroup,
      addMember,
      searchFunction,
      clearSearchInput,
      user,
    };
  },
};
</script>
  
<style scoped>
form {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 20px;
}

.createPost {
  width: 60vw;
  margin: 0 auto;
  margin-bottom: 20px;
  padding: 5px 10px;
  border-radius: 1rem;
  background-color: #f0f0f0;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.createPost label {
  font-weight: bold;
}

.createPost textarea {
  width: 98%;
  min-height: 100px;
  max-height: 300px;
  border: none;
  border-radius: 0.5rem;
  background-color: #ffffff;
  font-size: 14px;
  resize: vertical;
  padding: 15px;
  padding-right: 0;
  margin-bottom: 15px;
}

.createPost textarea::placeholder {
  color: #888;
  font-style: italic;
}

#groupTitle {
  width: 98%;
  border: none;
  border-radius: 0.5rem;
  background-color: #ffffff;
  font-size: 14px;
  padding: 15px;
  padding-right: 0;
  margin-bottom: 15px;
}

.createPost-btn {
  display: flex;
  margin: auto;
  margin-bottom: 10px;
}

#memberSearch {
  position: relative;
  padding-bottom: 10px;
  overflow: auto;
  z-index: 1;
}

#memberSearch a {
  display: block;
  color: black;
  text-decoration: none;
}

#myInput2 {
  width: 80%;
  padding: 5px 10px;
  margin-top: 10px;
  border: 2px solid #aeaeae;
  border-radius: 0.5rem;
  background-color: #f6f6f6;
}

#myInput2::placeholder {
  color: #888;
  font-style: italic;
}

#options2 {
  display: none;
  width: 80%;
  padding: 5px 0;
  margin: auto;
  border-radius: 0 0 1rem 1rem;
  box-shadow: #aeaeae 0px 0px 10px;
  background-color: #f6f6f6;
}

#options2 p {
  position: relative;
  z-index: 2;
}

#options2 p:hover {
  cursor: pointer;
  text-decoration: underline;
  color: #309668;
}

.selectedMembers {
  margin: 10px;
}
</style>@/composables/getFollowing