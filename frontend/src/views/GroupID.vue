<template>
  <main>
    <div v-if="groupInfo" class="group">
      <div class="group-header">
        <div class="group-title">
          <h3>{{ groupInfo.title }}</h3>
        </div>
        <!-- Display the user avatar -->
        <img
          class="user-avatar"
          v-if="userAvatarUrl(groupInfo.user_id)"
          :src="userAvatarUrl(groupInfo.user_id)"
        />
        <router-link
          class="group-creator"
          :to="{ name: 'userid', params: { userID: groupInfo.user_id } }"
          @click.stop
        >
          <strong>{{ groupInfo.creator_name }}</strong>
        </router-link>
        <p class="group-date">
          {{ formatGroupDateTime(groupInfo.created_at) }}
        </p>
      </div>
      <div class="group-description">
        <p>{{ groupInfo.content }}</p>
      </div>

      <!-- Join Group -->
      <button
        class="group-button"
        v-if="isMember === 'rejected' || isMember === 'not_member'"
        @click="joinGroup"
      >
        Join Group
      </button>
      <!-- Leave Group -->
      <button
        class="group-button"
        v-else-if="isMember === 'accepted' && !isCreator"
        @click="leaveGroup"
      >
        Leave Group
      </button>
      <!-- Invited -->
      <button
        class="group-button"
        v-else-if="isMember === 'invited'"
        disabled
        :class="{ 'button-disabled': isMember === 'invited' }"
      >
        You were invited, check notifications
      </button>
      <!-- Pending request -->
      <button
        v-else-if="isMember === 'request'"
        disabled
        :class="{ 'button-disabled': isMember === 'request' }"
      >
        Your request pending
      </button>
    </div>
    <div v-if="isMember === 'accepted'" class="side-menu">
      <ul>
        <li
          @click="selectMenuItem('posts')"
          :class="{ active: selectedMenuItem === 'posts' }"
        >
          Posts
        </li>
        <li
          @click="selectMenuItem('events')"
          :class="{ active: selectedMenuItem === 'events' }"
        >
          Events
        </li>
        <li
          @click="selectMenuItem('invite')"
          :class="{ active: selectedMenuItem === 'invite' }"
        >
          Invite
        </li>
      </ul>
    </div>
    <div v-if="isMember === 'accepted'" class="content-area">
      <div v-if="selectedMenuItem === 'posts'">
        <h3>Feed</h3>
        <div class="createPost">
          <form @submit.prevent="createNewPost">
            <div>
              <textarea
                v-model="content"
                placeholder="What's on your mind?"
                required
              ></textarea>
            </div>
            <div>
              <label for="image">Image: </label>
              <input type="file" @change="handleImageChange" />
            </div>
            <button type="submit">Create Post</button>
          </form>
        </div>
        <div v-if="groupPosts && groupPosts.length > 0">
          <div
            v-for="post in groupPosts"
            :key="post.id"
            class="post"
            @click="redirectToPost(post.post_id)"
          >
            <div class="post-header">
              <router-link
                :to="{ name: 'userid', params: { userID: post.user_id } }"
                class="post-link"
                @click.stop
              >
                <strong>{{ post.full_name }} </strong>
              </router-link>
              <p class="post-date">
                {{ formatGroupDateTime(post.created_at) }}
              </p>
              <p class="privacy-level">{{ post.privacy_level }}</p>
            </div>
            <div class="post-content">
              <p>{{ post.content }}</p>
            </div>
            <img
              v-if="post.post_image"
              :src="post.post_image"
              alt="Post Image"
            />
          </div>
        </div>
        <div v-else>No posts in this group.</div>
      </div>
      <div v-else-if="selectedMenuItem === 'events'">
        <h3>Events</h3>
        <form class="createEvent" @submit.prevent="handleSubmit">
          <div>
            <input
              v-model="eventTitle"
              placeholder="Event Title"
              class="event-title"
              required
            />
          </div>
          <div class="event-container">
            <textarea
              v-model="eventContent"
              placeholder="What's it about?!"
              required
            ></textarea>
          </div>
          <div>
            <label for="eventDate">Select Time & Date: </label>
            <input
              type="datetime-local"
              id="eventDate"
              v-model="eventDate"
              required
            />
          </div>
          <button type="submit">Create Event</button>
        </form>
        <div v-if="events && events.length > 0">
          <div v-for="event in events" :key="event.event_id" class="event">
            <p>
              <strong>{{ event.title }}</strong>
            </p>
            <p>{{ event.content }}</p>
            <p>Date: {{ formatGroupDateTime(event.date) }}</p>
            <p>Creator: {{ event.creator_name }}</p>
            <p>Your status: {{ event.attendeeStatus }}</p>
            <div class="eventResponse" v-if="event.attendeeStatus === 'nil'">
              <button @click="handleAttendeesStatusChange(event, 'going')">
                Going
              </button>
              <button @click="handleAttendeesStatusChange(event, 'not going')">
                Not Going
              </button>
            </div>
            <div v-else-if="event.attendeeStatus === 'going'">
              <button @click="handleAttendeesStatusChange(event, 'not going')">
                Change to Not Going
              </button>
            </div>
            <div v-else-if="event.attendeeStatus === 'not going'">
              <button @click="handleAttendeesStatusChange(event, 'going')">
                Change to Going
              </button>
            </div>
            <hr />
          </div>
        </div>
        <div v-else>
          <p>No events available.</p>
        </div>
      </div>
      <div v-else-if="selectedMenuItem === 'invite'">
        <h3>Invite new members</h3>
        <form class="newform" @submit.prevent="inviteSelectedMembers">
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
              <div v-for="user in filteredUsers" :key="user.user_id">
                <p @click="addMember(user)">
                  {{ user.firstname }} {{ user.lastname }}
                </p>
              </div>
            </div>
          </div>
          <!-- Member list -->
          <div class="selectedMembers">
            <div v-for="member in selectedMembers" :key="member.user_id">
              {{ member.firstname }} {{ member.lastname }}
            </div>
          </div>
          <button type="submit" @click="clearSearchInput">
            Invite them now
          </button>
        </form>
      </div>
    </div>
    <div v-else>See more infos only for followers</div>
  </main>
</template>

<script>
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import createGroupPost from "@/composables/createGroupPost";
import getGroupInfo from "@/composables/getGroupInfo";
import getUserFromSession from "@/composables/getUserFromSession";
import getUsers from "@/composables/getUsers";
import searchMembers from "@/composables/searchMembers";
import {
  joinGroup,
  leaveGroup,
  inviteMembers,
} from "@/composables/changeGroupMembers";
import {
  createEvent,
  fetchEventData,
  fetchAttendeesStatus,
  updateAttendeesStatus,
} from "@/composables/eventLogic";

export default {
  name: "GroupID",
  computed: {
    userAvatarUrl() {
      return (userId) => {
        const user = this.users.find((user) => user.user_id === userId);
        return user
          ? user.avatar
          : "http://localhost:8000/uploads/avatars/default-avatar-profile.jpg";
      };
    },
    filteredUsers() {
      if (!this.users || !this.groupInfo) return [];

      // Create a set of member IDs for faster lookup
      const memberSet = new Set(this.groupInfo.members);
      console.log("this is group members set", memberSet);

      // Filter out users who are already members of the group
      return this.users.filter((user) => !memberSet.has(user.user_id));
    },
  },
  methods: {
    redirectToPost(postID) {
      this.$router.push({ name: "PostID", params: { postID: postID } });
    },
  },
  setup() {
    const route = useRoute();
    const isCreator = ref(false);
    const selectedMembers = ref([]);

    const {
      groupInfo,
      fetchGroupInfo,
      groupPosts,
      fetchGroupPosts,
      isMember,
      checkMembershipStatus,
      updateMembershipStatus,
    } = getGroupInfo();
    const { content, image, groupPost } = createGroupPost();
    const { user: sessionUserID, fetchUserDataFromSession } =
      getUserFromSession();
    const { users, loadUsers } = getUsers();
    const { events, getEventData } = fetchEventData();
    const { eventTitle, eventContent, eventDate, createNewEvent } =
      createEvent();
    const { searchFunction, clearSearchInput } = searchMembers();

    onMounted(async () => {
      await checkMembershipStatus(route.params.groupID);
      await fetchGroupInfo(route.params.groupID);
      await fetchGroupPosts(route.params.groupID);
      await fetchUserDataFromSession();
      await loadUsers();
      await getEventData(route.params.groupID);

      if (groupInfo.value.user_id === sessionUserID.value.userID) {
        isCreator.value = true;
      }
    });

    const selectedMenuItem = ref("posts"); // Default selected menu item
    const selectMenuItem = (menuItem) => {
      selectedMenuItem.value = menuItem;
    };

    const createNewPost = async () => {
      try {
        await groupPost(route.params.groupID);
      } catch (error) {
        console.error("Error creating post:", error);
      }
      fetchGroupPosts(route.params.groupID);
    };

    const handleSubmit = async () => {
      await createNewEvent(route.params.groupID);
      await getEventData(route.params.groupID);
    };

    const attendeeStatus = ref(null);
    const getAttendeeStatus = async (eventID) => {
      try {
        const status = await fetchAttendeesStatus(eventID);
        attendeeStatus.value = status;
      } catch (error) {
        console.error("Error getting attendee status:", error);
        return null;
      }
    };

    const handleAttendeesStatusChange = async (event, status) => {
      try {
        await updateAttendeesStatus(event.event_id, status);
        // Update attendeeStatus in UI
        event.attendeeStatus = status;
      } catch (error) {
        console.error("Error updating attendees status:", error);
      }
    };

    // Function to format post date and time
    const formatGroupDateTime = (createdAt) => {
      const postDate = new Date(createdAt);
      const time = postDate.toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });
      const day = ("0" + postDate.getDate()).slice(-2);
      const month = ("0" + (postDate.getMonth() + 1)).slice(-2);
      const year = postDate.getFullYear();
      return `${time}, ${day}.${month}.${year}`;
    };

    const handleImageChange = (event) => {
      const file = event.target.files[0];
      if (!file) {
        return;
      }
      const reader = new FileReader();
      reader.onload = (e) => {
        // Set the base64 string to the image ref
        image.value = e.target.result;
      };
      reader.readAsDataURL(file);
    };

    //adds members to selectedMembers list
    const addMember = (user) => {
      if (
        user &&
        !selectedMembers.value.find((member) => member.user_id === user.user_id)
      ) {
        selectedMembers.value.push(user);
      }
    };

    const inviteSelectedMembers = async () => {
      await inviteMembers(route.params.groupID, selectedMembers.value);
      selectedMembers.value = [];
    };

    return {
      groupInfo,
      groupPosts,
      isMember,
      isCreator,
      content,
      formatGroupDateTime,
      selectedMenuItem,
      selectMenuItem,
      createNewPost,
      handleImageChange,
      joinGroup: () => joinGroup(route.params.groupID, updateMembershipStatus),
      leaveGroup: () =>
        leaveGroup(route.params.groupID, updateMembershipStatus),
      eventTitle,
      eventContent,
      eventDate,
      handleSubmit,
      handleAttendeesStatusChange,
      events,
      getAttendeeStatus,
      attendeeStatus,
      sessionUserID,
      users,
      searchFunction,
      clearSearchInput,
      addMember,
      selectedMembers,
      inviteSelectedMembers,
    };
  },
};
</script>

<style scoped>
.group {
  width: 60vw;
  margin: 0 auto;
  margin-top: 20px;
  margin-bottom: 20px;
  padding-bottom: 15px;
  display: flex;
  flex-direction: column;
  align-items: left;
}

.group-header {
  width: 100%;
  display: flex;
  justify-content: flex-end;
  border-bottom: 2px solid #bbb;
}

.group-creator {
  margin-right: 10px;
  align-self: center;
  text-decoration: none;
  color: inherit;
}

.group-creator:hover {
  text-decoration: underline;
  color: #888;
  cursor: pointer;
}

.group-date {
  font-size: 14px;
  color: #888;
  margin-right: 10px;
  align-self: center;
}

.group-title {
  text-align: left;
  display: flex;
  margin-right: auto;
}

.group-description {
  background-color: white;
  text-align: left;
  margin-bottom: 15px;
}

.group-button {
  border: none;
  width: 15%;
  border-radius: 5px;
  cursor: pointer;
  font-size: 16px;
  padding: 7px 10px;
  margin-bottom: 0;
  color: #fff;
  margin-left: auto;
  background-color: #42b983;
}

.side-menu {
  position: fixed;
  top: 165px;
  left: 4.2%;
  overflow: auto;
  width: 14.5%;
}

.side-menu li:hover {
  background-color: #e6e6e6;
  mix-blend-mode: multiply;
  border-radius: 1rem;
}

.side-menu ul {
  list-style-type: none;
  padding: 0;
  font-weight: bold;
}

.side-menu ul li {
  padding: 10px;
  cursor: pointer;
}

.side-menu ul li.active {
  background-color: #ddd;
  color: #42b983;
  text-decoration: underline 2px;
  border-radius: 1rem;
}

form,
.event-container {
  margin-top: 15px;
  width: 90%;
}

.newform {
  margin-left: 80px;
}

.createPost,
.createEvent {
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

.createEvent textarea,
.createPost textarea {
  width: 98%;
  min-height: 100px;
  max-height: 300px;
  border: none;
  border-radius: 0.5rem;
  background-color: #ffffff;
  font-size: 14px;
  resize: vertical;
  padding-left: 15px;
  padding-top: 15px;
  margin-bottom: 15px;
}

.event-title {
  border: none;
  border-radius: 0.5rem;
  padding: 5px;
  margin-top: 10px;
  font-size: 14px;
  font-style: italic;
}

.createEvent input[type="datetime-local"] {
  border: none;
  border-radius: 1rem;
  padding: 5px;
  margin-bottom: 10px;
}

.createPost input[type="file"] {
  width: 25vw;
  padding: 15px;
  margin-left: 40px;
}

.createEvent textarea::placeholder,
.createPost textarea::placeholder {
  color: #888;
  font-style: italic;
}

.post-header {
  width: 90%;
  display: flex;
  justify-content: flex-end;
  align-items: flex-end;
}

.post:hover {
  background-color: #e6e6e6;
  cursor: pointer;
}

.post-link {
  margin-right: auto;
  align-self: center;
  text-decoration: none;
  color: inherit;
  padding-left: 3%;
}

.post-link:hover {
  text-decoration: underline;
  color: #888;
}

.post-date {
  font-size: 14px;
  color: #888;
  margin-right: 10px;
  text-align: left;
  align-self: flex-end;
}

.post-content {
  background-color: rgb(255, 255, 255);
  border-radius: 5px;
  width: 90%;
  text-align: left;
  margin-bottom: 10px;
}

.post-content p {
  margin-left: 3%;
  text-align: left;
}

.post img {
  max-width: 90%;
  height: auto;
  border-radius: 5px;
  margin-bottom: 15px;
}

.button-disabled {
  background-color: grey;
  color: white;
}

.user-avatar {
  width: 30px;
  aspect-ratio: 1 / 1;
  clip-path: circle(50%);
  align-self: center;
  margin-bottom: 0px !important;
  margin-right: 0.7%;
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
  width: 25%;
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
  width: 25%;
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
  margin: 15px;
}

.eventResponse button:first-child {
  margin-right: 10px;
}
</style>
