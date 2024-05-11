<template>
  <main>
    <h2>View Post</h2>

    <div v-if="display">
      <div v-if="post" class="post">
        <div class="post-header">
          <!-- Display the user avatar -->
          <img
            class="user-avatar"
            v-if="userAvatarUrl(post.user_id)"
            :src="userAvatarUrl(post.user_id)"
          />
          <router-link
            :to="
              post.user_id === sessionUserID
                ? { name: 'Profile' }
                : { name: 'userid', params: { userID: post.user_id } }
            "
            class="post-link"
          >
            <strong> {{ post.full_name }} </strong>
          </router-link>
          <p class="post-date">{{ formatDateTime(post.created_at) }}</p>
          <p class="privacy-level">
            <span v-if="post.group_id">
              <router-link
                :to="{ name: 'GroupID', params: { groupID: post.group_id } }"
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

      <!-- Previous comments -->
      <h3>Comments:</h3>
      <div v-if="comments && comments.length > 0">
        <div v-for="comment in comments" :key="comment.id" class="post">
          <div class="post-header">
            <!-- Display the user avatar -->
            <img
              class="user-avatar"
              v-if="userAvatarUrl(comment.user_id)"
              :src="userAvatarUrl(comment.user_id)"
            />
            <router-link
              :to="
                comment.user_id === sessionUserID
                  ? { name: 'Profile' }
                  : { name: 'userid', params: { userID: comment.user_id } }
              "
              class="post-link"
            >
              <strong> {{ comment.full_name }} </strong>
            </router-link>
            <p class="post-date">{{ formatDateTime(comment.created_at) }}</p>
          </div>
          <div class="post-content">
            <p>{{ comment.content }}</p>
          </div>
          <img
            v-if="comment.comment_image"
            :src="comment.comment_image"
            alt="Comment Image"
          />
        </div>
      </div>
      <div v-else>No comments yet.</div>

      <div class="break"></div>

      <!-- New comment form -->
      <div class="comment">
        <h3>Add a New Comment</h3>
        <form @submit.prevent="createComment">
          <textarea
            v-model="newCommentContent"
            placeholder="Write your comment here"
            required
          ></textarea>
          <input type="file" @change="handleCommentImageChange" />
          <button type="submit">Create Comment</button>
        </form>
      </div>
    </div>
    <div v-else>Post not found.</div>
  </main>
</template>
  
<script>
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import getCommentsForPost from "@/composables/getCommentsForPost";
import getPostFromPostID from "@/composables/getPostFromPostID";
import getGroupInfo from "@/composables/getGroupInfo";
import getViewerStatus from "@/composables/getViewerStatus";
import getFollowStatus from "@/composables/getFollowStatus";
import getUsers from "@/composables/getUsers";
import getUserFromSession from "@/composables/getUserFromSession";

export default {
  name: "PostID",
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
  setup() {
    const route = useRoute();
    const { comments, fetchCommentsForPost, addComment } = getCommentsForPost();
    const { post, fetchPostFromID } = getPostFromPostID();
    const { groupInfo, fetchGroupInfo, isMember, checkMembershipStatus } =
      getGroupInfo();
    const { status, loadViewer } = getViewerStatus();
    const { users, loadUsers } = getUsers();
    const { user: sessionUserID, fetchUserDataFromSession } =
      getUserFromSession();
    const { followStatus, fetchFollowStatus } = getFollowStatus();

    const newCommentContent = ref("");
    const commentImage = ref(null);
    const postLoaded = ref(false);
    const display = ref(false);

    onMounted(async () => {
      await loadUsers();
      await fetchCommentsForPost(route.params.postID);
      await fetchPostFromID(route.params.postID);
      await fetchUserDataFromSession();
      if (post.value) {
        await fetchFollowStatus(post.value.user_id);
      }

      postLoaded.value = true;
      display.value = await shouldDisplayPost();
    });

    // Function to format post date and time
    const formatDateTime = (createdAt) => {
      const date = new Date(createdAt);
      const time = date.toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });
      const day = ("0" + date.getDate()).slice(-2);
      const month = ("0" + (date.getMonth() + 1)).slice(-2);
      const year = date.getFullYear();
      return `${time}, ${day}.${month}.${year}`;
    };

    const handleCommentImageChange = (event) => {
      const file = event.target.files[0];
      if (!file) {
        return;
      }
      const reader = new FileReader();
      reader.onload = (e) => {
        // Set the base64 string to the image ref
        commentImage.value = e.target.result;
      };
      reader.readAsDataURL(file);
    };

    const createComment = async () => {
      const commentData = {
        content: newCommentContent.value,
        comment_image: commentImage.value,
      };

      await addComment(route.params.postID, commentData);

      newCommentContent.value = "";
      commentImage.value = null;

      fetchCommentsForPost(route.params.postID);
    };

    const shouldDisplayPost = async () => {
      if (!postLoaded.value || !post.value) {
        return false;
      }

      if (post.value.group_id) {
        await fetchGroupInfo(post.value.group_id);
        await checkMembershipStatus(post.value.group_id);
        return isMember.value;
      } else {
        if (post.value.privacy_level === "public") {
          return true;
        } else if (post.value.privacy_level === "private") {
          if (
            followStatus.value == "accepted" ||
            post.value.user_id === sessionUserID.value.userID
          ) {
            return true;
          } else {
            return false;
          }
        } else if (post.value.privacy_level === "friends") {
          await loadViewer(post.value.post_id);
          return status.value;
        } else {
          return false;
        }
      }
    };

    return {
      post,
      sessionUserID,
      comments,
      groupInfo,
      formatDateTime,
      newCommentContent,
      handleCommentImageChange,
      createComment,
      shouldDisplayPost,
      display,
      users,
    };
  },
};
</script>

<style scoped>
.post {
  width: 60vw;
}

.post img {
  max-width: 90%;
  height: auto;
  border-radius: 5px;
  margin-bottom: 15px;
}

.comment {
  width: 60vw;
  margin: 0 auto;
  margin-bottom: 20px;
  padding: 20px 40px;
  border-radius: 1rem;
  background-color: #f0f0f0;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.comment textarea {
  min-height: 75px;
  max-height: 300px;
  border: none;
  border-radius: 0.5rem;
  background-color: white;
  font-size: 14px;
  resize: vertical;
  padding: 15px;
  margin-bottom: 15px;
}

.comment input[type="file"] {
  text-align: center;
  margin-bottom: 15px;
}

.comment textarea::placeholder {
  color: #888;
  font-style: italic;
}

.createComment-btn {
  display: flex;
  margin: auto;
  margin-bottom: 10px;
}

.break {
  margin: 30px;
}
</style>
  