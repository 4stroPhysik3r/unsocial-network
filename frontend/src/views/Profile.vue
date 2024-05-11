<template>
  <main>
    <h1>Profile Page</h1>

    <!-- Check if user data is available before rendering -->
    <div v-if="user">
      <div class="user-info">
        <div class="avatar-stats-container">
          <div class="avatar">
            <img :src="user.avatar" v-if="user.avatar" alt="Avatar" />
          </div>

          <div class="stats">
            <div class="follower" v-if="follower">
              <h2>Followers: {{ follower.length }}</h2>
              <!-- <ul>
                <li
                  v-for="followedUser in sortedFollowers"
                  :key="followedUser.user_id"
                >
                  {{ followedUser.full_name }}
                </li>
              </ul> -->
            </div>
            <div v-else class="follower"><h2>Followers: 0</h2></div>

            <div class="following" v-if="following">
              <h2>Following: {{ following.length }}</h2>
              <!-- <ul>
                <li
                  v-for="followingUser in sortedFollowing"
                  :key="followingUser.user_id"
                >
                  {{ followingUser.full_name }}
                </li>
              </ul> -->
            </div>
            <div v-else class="following"><h2>Following: 0</h2></div>
          </div>
        </div>

        <div class="details">
          <p><strong>Email:</strong> {{ user.email }}</p>
          <p><strong>First Name:</strong> {{ user.firstName }}</p>
          <p><strong>Last Name:</strong> {{ user.lastName }}</p>
          <p><strong>Date of Birth:</strong> {{ formattedDateOfBirth() }}</p>
          <p v-if="user.nickname">
            <strong>Nickname:</strong> {{ user.nickname }}
          </p>
          <p v-if="user.aboutMe">
            <strong>About me:</strong> {{ user.aboutMe }}
          </p>
          <p>
            <strong>Profile:</strong>
            {{ user.profilePublic ? "Public" : "Private" }}
          </p>
          <!-- Hide toggleProfileStatus button if viewing another user's profile -->
          <button v-if="isAuthenticated" @click="toggleProfileStatus">
            {{
              user.profilePublic
                ? "Make Profile Private"
                : "Make Profile Public"
            }}
          </button>
          <button v-if="isAuthenticated" @click="performLogout">Logout</button>
        </div>
      </div>

      <div class="break"></div>

      <!-- Display posts -->
      <div v-if="posts">
        <h2>Your activity</h2>

        <div
          v-for="post in posts"
          :key="post.id"
          class="post"
          @click="redirectToPost(post.post_id)"
        >
          <div class="post-header">
            <p class="post-date">{{ formatPostDateTime(post.created_at) }}</p>
            <p class="privacy-level">
              <span v-if="post.group_id">
                <router-link
                  :to="{ name: 'GroupID', params: { groupID: post.group_id } }"
                >
                  GroupID: {{ post.group_id }}<br />
                </router-link>
              </span>
              <span v-else>{{ post.privacy_level }}</span>
            </p>
          </div>

          <div class="post-content">
            <p>{{ post.content }}</p>
          </div>

          <img v-if="post.post_image" :src="post.post_image" alt="Post Image" />
        </div>
      </div>

      <div v-else-if="postsError">
        <p>{{ postsError }}</p>
      </div>

      <div v-else>
        <p>No posts yet from this users</p>
      </div>
    </div>

    <div v-else>
      <p>Loading user data...</p>
    </div>
  </main>
</template>
  
<script>
import { computed, onMounted, watch } from "vue";
import { useRouter } from "vue-router";
import { useAuth } from "@/composables/useAuth";
import getUserFromSession from "@/composables/getUserFromSession";
import getPostsForProfile from "@/composables/getPostsForProfile";
import getFollowing from "@/composables/getFollowing";
import getFollower from "@/composables/getFollower";

export default {
  name: "Profile",
  methods: {
    redirectToPost(postID) {
      this.$router.push({ name: "PostID", params: { postID: postID } });
    },
  },
  setup() {
    const router = useRouter();
    const { isAuthenticated, logout } = useAuth();
    const { posts, loadPosts } = getPostsForProfile();
    const { user, fetchUserDataFromSession } = getUserFromSession();
    const { following, loadFollowing } = getFollowing();
    const { follower, loadFollower } = getFollower();

    // Fetch user data when the component is mounted
    onMounted(async () => {
      await loadPosts();
      await fetchUserDataFromSession();
      await loadFollowing(user.value.userID);
      await loadFollower(user.value.userID);
    });

    // Function to sort posts by date, from newest to oldest
    const sortPostsByDate = () => {
      if (posts.value) {
        posts.value.sort(
          (a, b) => new Date(b.created_at) - new Date(a.created_at)
        );
      }
    };

    watch(async () => {
      sortPostsByDate();
    });

    const toggleProfileStatus = async () => {
      if (
        window.confirm("Are you sure you want to change your profile status?")
      ) {
        user.value.profilePublic = !user.value.profilePublic;

        try {
          const response = await fetch("http://localhost:8000/api/my-profile", {
            method: "PUT",
            credentials: "include",
            body: JSON.stringify(user.value.profilePublic),
          });
          if (!response.ok) {
            throw new Error("Failed to update profile status");
          }
          console.log("Profile status updated successfully");
        } catch (error) {
          console.error("Error updating profile status:", error);
          // Revert the profile status change if there's an error
          user.value.profilePublic = !user.value.profilePublic;
        }
      }
    };

    // Computed property to format date of birth
    const formattedDateOfBirth = () => {
      if (user.value.dateOfBirth) {
        const dateOfBirth = new Date(user.value.dateOfBirth);
        const year = dateOfBirth.getFullYear();
        const month = ("0" + (dateOfBirth.getMonth() + 1)).slice(-2);
        const day = ("0" + dateOfBirth.getDate()).slice(-2);
        return `${day}.${month}.${year}`;
      }
      return "";
    };

    // Function to format post date and time
    const formatPostDateTime = (createdAt) => {
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

    // Function to log out the user
    const performLogout = async () => {
      await logout();
      router.push("/login"); // Redirect to login page after logout
    };

    const sortUsers = (users) => {
      return users.slice().sort((a, b) => {
        return a.full_name.localeCompare(b.full_name);
      });
    };

    const sortedFollowers = computed(() => {
      return sortUsers(follower.value);
    });

    const sortedFollowing = computed(() => {
      return sortUsers(following.value);
    });

    return {
      isAuthenticated,
      user,
      posts,
      performLogout,
      toggleProfileStatus,
      formattedDateOfBirth,
      formatPostDateTime,
      following,
      follower,
      sortedFollowers,
      sortedFollowing,
    };
  },
};
</script>

<style scoped>
.avatar-stats-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.user-info {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  align-items: center;
}

.avatar img {
  width: 100px;
  height: 100px;
  border-radius: 50%;
}

.details {
  text-align: left;
}

.user-activity {
  margin-top: 4rem;
}

.stats {
  display: flex;
}

.follower {
  margin-right: 2rem;
}

.following ul,
.follower ul {
  list-style: none;
  padding: 0;
}

.following li,
.follower li {
  margin-bottom: 5px;
}

button {
  margin-right: 12px;
}

.post-header {
  width: 90%;
  display: flex;
  justify-content: flex-end;
  align-items: flex-end;
}

.post-date {
  margin-right: 10px;
  font-size: 14px;
  color: #888;
  text-align: left;
  align-self: flex-end;
}

.privacy-level {
  font-size: 14px;
  color: #888;
  align-self: flex-end;
  padding-right: 3%;
}

.post-content {
  background-color: white;
  border-radius: 5px;
  width: 90%;
  text-align: left;
  margin-bottom: 10px;
}

.post-content p {
  margin-left: 3%;
  text-align: left;
}

.post:hover {
  background-color: #ccc;
  cursor: pointer;
}

.post img {
  max-width: 90%;
  height: auto;
  border-radius: 5px;
  margin-bottom: 15px;
}
</style>