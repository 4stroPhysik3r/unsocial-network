<template>
  <main>
    <h1>User Profile</h1>

    <div v-if="user">
      <div class="user-info">
        <div class="avatar-stats-container">
          <div class="avatar">
            <img :src="user.avatar" v-if="user.avatar" />
            <img
              v-else
              src="http://localhost:8000/uploads/avatars/default-avatar-profile.jpg"
              alt="Avatar"
            />
          </div>

          <div class="stats">
            <div class="follower" v-if="follower">
              <h2>Followers: {{ follower.length }}</h2>
              <ul>
                <li
                  v-for="followedUser in follower"
                  :key="followedUser.user_id"
                >
                  {{ followedUser.full_name }}
                </li>
              </ul>
            </div>
            <div v-else class="follower"><h2>Followers: 0</h2></div>

            <div class="following" v-if="following">
              <h2>Following {{ following.length }}</h2>
              <ul>
                <li
                  v-for="followingUser in following"
                  :key="followingUser.user_id"
                >
                  {{ followingUser.full_name }}
                </li>
              </ul>
            </div>
            <div v-else class="following"><h2>Following: 0</h2></div>
          </div>
        </div>

        <div class="details">
          <p><strong>Email:</strong> {{ user.email }}</p>
          <p><strong>First Name:</strong> {{ user.firstName }}</p>
          <p><strong>Last Name:</strong> {{ user.lastName }}</p>

          <div v-if="user.profilePublic">
            <p><strong>Date of Birth:</strong> {{ formattedDateOfBirth() }}</p>
            <p v-if="user.nickname">
              <strong>Nickname:</strong> {{ user.nickname }}
            </p>
            <p v-if="user.aboutMe">
              <strong>About me:</strong> {{ user.aboutMe }}
            </p>
          </div>

          <p>
            <strong>Profile:</strong>
            {{ user.profilePublic ? "Public" : "Private" }}
          </p>

          <!-- Show "Follow" button for users with 'rejected' or 'not_follower' status -->
          <div v-if="!isCreator" class="follow-button">
            <button
              v-if="followStatus === 'rejected'"
              @click="follow"
              class="follow-button"
            >
              Follow
            </button>

            <!-- Show "Unfollow" button for users with 'accepted' status -->
            <button
              v-else-if="followStatus === 'accepted'"
              @click="unfollow"
              class="unfollow-button"
            >
              Unfollow
            </button>

            <!-- Show disabled "Pending" button for users with 'pending' status -->

            <button
              v-else-if="followStatus === 'pending'"
              disabled
              class="button-disabled"
            >
              Pending
            </button>
          </div>
        </div>
      </div>

      <div class="break"></div>

      <!-- Display posts -->
      <div v-if="posts">
        <h2>User activity</h2>
        <div
          v-for="post in filteredPosts"
          :key="post.id"
          class="post"
          @click="redirectToPost(post.post_id)"
        >
          <div class="post-header">
            <img v-if="user.avatar" :src="user.avatar" class="user-avatar" />
            <img
              v-else
              src="http://localhost:8000/uploads/avatars/default-avatar-profile.jpg"
              class="user-avatar"
              alt="Avatar"
            />
            <strong class="post-link">
              {{ user.firstName + " " + user.lastName }}
            </strong>
            <p class="post-date">{{ formatPostDateTime(post.created_at) }}</p>
            <p class="privacy-level">
              <span v-if="post.group_id">
                <router-link
                  :to="{
                    name: 'GroupID',
                    params: { groupID: post.group_id },
                  }"
                >
                  {{ groupInfo.title }}
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

      <div
        v-else-if="
          user.profilePublic !== 'private' && followStatus !== 'accepted'
        "
      >
        <p>Profile is private.</p>
      </div>

      <div v-else>
        <p>No posts yet from this user</p>
      </div>
    </div>

    <div v-else>
      <p>Loading user data...</p>
    </div>
  </main>
</template>

<script>
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useAuth } from "@/composables/useAuth";
import useFollow from "@/composables/useFollow";
import getPostsForProfile from "@/composables/getPostsForProfile";
import getUserFromID from "@/composables/getUserFromID";
import getUserFromSession from "@/composables/getUserFromSession";
import getViewerStatus from "@/composables/getViewerStatus";
import getFollowing from "@/composables/getFollowing";
import getFollower from "@/composables/getFollower";

export default {
  name: "User",
  methods: {
    redirectToPost(postID) {
      this.$router.push({ name: "PostID", params: { postID: postID } });
    },
  },
  setup() {
    const { isAuthenticated } = useAuth();

    const router = useRouter();
    const { userID } = router.currentRoute.value.params;
    const { user, fetchUserDataFromID } = getUserFromID(userID);

    const { user: sessionUserID, fetchUserDataFromSession } =
      getUserFromSession();
    const { posts, loadPosts } = getPostsForProfile(userID);
    const { status, loadViewer } = getViewerStatus();
    const { followStatus, fetchFollowStatus, follow, unfollow } =
      useFollow(userID);
    const { following, loadFollowing } = getFollowing();
    const { follower, loadFollower } = getFollower();

    const isCreator = ref(false);

    onMounted(async () => {
      await fetchUserDataFromID(userID);
      await fetchUserDataFromSession();
      await loadPosts();
      await fetchFollowStatus();
      await loadFollowing(user.value.userID);
      await loadFollower(user.value.userID);

      if (user.value.userID === sessionUserID.value.userID) {
        isCreator.value = true;
        router.push({ name: "Profile" });
      }
    });

    const filteredPosts = computed(() => {
      if (!posts.value) {
        return [];
      }

      return posts.value.filter((post) => {
        if (!post.group_id) {
          if (post.privacy_level === "public") {
            return true;
          } else if (
            post.privacy_level === "private" &&
            followStatus.value === "accepted"
          ) {
            return true;
          } else if (post.privacy_level === "friends") {
            loadViewer(post.post_id);
            if (status.value === "accepted") {
              return true;
            }
          }
        }
      });
    });

    const formattedDateOfBirth = () => {
      if (user.value.dateOfBirth) {
        const dateOfBirth = new Date(user.value.dateOfBirth);
        const day = ("0" + dateOfBirth.getDate()).slice(-2);
        const month = ("0" + (dateOfBirth.getMonth() + 1)).slice(-2);
        const year = dateOfBirth.getFullYear();
        return `${day}.${month}.${year}`;
      }
      return "";
    };

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

    return {
      isAuthenticated,
      user,
      posts,
      filteredPosts,
      formattedDateOfBirth,
      loadPosts,
      formatPostDateTime,
      followStatus,
      follow,
      unfollow,
      following,
      follower,
      isCreator,
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
  width: 120px;
  aspect-ratio: 50%;
  border-radius: 50%;
}

.details {
  justify-content: center;
  text-align: left;
  padding: 20px;
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

.profile-status {
  margin-top: 20px;
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

.post-link {
  margin-right: auto;
  align-self: center;
  padding-left: 3%;
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
  margin-bottom: 10px;
}

.user-avatar {
  width: 30px;
  aspect-ratio: 1 / 1;
  clip-path: circle(50%);
  align-self: center;
  margin-bottom: 5px !important;
  margin-left: 3%;
}

.button-disabled {
  pointer-events: none;
  background-color: grey;
}
</style>