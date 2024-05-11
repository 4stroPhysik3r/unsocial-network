<template>
  <main>
    <div class="createPost">
      <form @submit.prevent="createPost">
        <textarea
          class="post-input"
          v-model="content"
          placeholder="What's on your mind?"
          required
        ></textarea>

        <div class="create-form">
          <label for="privacy">Privacy Level: </label>
          <select id="privacy" v-model="privacy" required>
            <option value="public">Public</option>
            <option value="private">Private</option>
            <option value="friends">Friends</option>
          </select>
        </div>

        <div v-if="privacy === 'friends'">
          <label>Select following:</label>
          <div v-for="follow in following" :key="follow.following_id">
            <input
              type="checkbox"
              :value="follow.follower_id"
              :id="'follower_' + follow.follower_id"
              v-model="selectedFollowing"
            />
            {{ follow.full_name }}
            {{ follow.follower_id }}
          </div>
        </div>

        <div class="create-form">
          <label for="image">Image: </label>
          <label for="image" class="custom-file-upload"> Upload Image </label>
          <input id="image" type="file" @change="handleImageChange" />
        </div>

        <button type="submit">Create Post</button>
      </form>
    </div>

    <div class="break"></div>

    <h2>Your Feed:</h2>

    <div v-if="posts && posts.length > 0">
      <div
        v-for="post in posts"
        :key="post.post_id"
        class="post"
        @click="redirectToPost(post.post_id)"
      >
        <div class="post-header">
          <!-- Display the user avatar -->
          <img
            class="user-avatar"
            v-if="userAvatarUrl(post.user_id)"
            :src="userAvatarUrl(post.user_id)"
          />

          <router-link
            class="post-link"
            :to="{ name: 'userid', params: { userID: post.user_id } }"
            @click.stop
          >
            <strong>{{ post.full_name }} </strong>
          </router-link>

          <p class="post-date">{{ formatPostDateTime(post.created_at) }}</p>
          <p class="privacy-level">{{ post.privacy_level }}</p>
        </div>

        <div class="post-content">
          <p>{{ post.content }}</p>
        </div>

        <img v-if="post.post_image" :src="post.post_image" alt="Post Image" />
      </div>
    </div>

    <div v-else>No posts to display.</div>
  </main>
</template>

<script>
import { onMounted, ref, watch } from "vue";
import getFollowingListForPost from "@/composables/getFollowingListForPost";
import getPosts from "@/composables/getPosts";
import getUsers from "@/composables/getUsers";

export default {
  name: "Home",
  computed: {
    // Compute the user avatar URL based on user_id
    userAvatarUrl() {
      return (userId) => {
        const user = this.users.find((user) => user.user_id === userId);
        // If user is found, return the avatar URL, otherwise return null (or an empty string)
        return user
          ? user.avatar
          : "http://localhost:8000/uploads/avatars/default-avatar-profile.jpg";
      };
    },
  },
  methods: {
    async performLogout() {
      const { logout } = useAuth();
      await logout();

      this.$router.push("/login");
    },
    redirectToPost(postID) {
      this.$router.push({ name: "PostID", params: { postID: postID } });
    },
  },
  setup() {
    const { following, loadFollowing } = getFollowingListForPost();
    const { posts, loadPosts } = getPosts();
    const { users, loadUsers } = getUsers();

    // Function to sort posts by date, from newest to oldest
    const sortPostsByDate = () => {
      if (posts.value) {
        posts.value.sort(
          (a, b) => new Date(b.created_at) - new Date(a.created_at)
        );
      }
    };

    onMounted(async () => {
      await loadFollowing();
      await loadPosts();
      await loadUsers();
    });

    watch(async () => {
      sortPostsByDate();
    });

    // TODO: clean up
    // logic for creating the posts
    const content = ref("");
    const privacy = ref("public");
    const image = ref(null);
    const selectedFollowing = ref([]);

    const createPost = async () => {
      const postData = {
        content: content.value,
        privacy_level: privacy.value,
        post_image: image.value,
        viewer_ids: selectedFollowing.value,
      };

      try {
        const response = await fetch("http://localhost:8000/api/create-post/", {
          method: "POST",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(postData),
        });

        if (!response.ok) {
          throw new Error("Failed to create post");
        }

        loadPosts(); // rerender the posts after creating new post

        content.value = "";
        privacy.value = "public";
        image.value = null;
      } catch (error) {
        console.error("Error creating post:", error);
      }
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

    return {
      following,
      selectedFollowing,
      content,
      privacy,
      image,
      createPost,
      posts,
      users,
      handleImageChange,
      formatPostDateTime,
    };
  },
};
</script>

<style >
form {
  display: flex;
  flex-direction: column;
  width: 100%;
  text-align: left;
}

.createPost {
  width: 40vw;
  margin: 0 auto;
  margin-top: 20px;
  padding: 20px 40px;
  border-radius: 1rem;
  background-color: #f0f0f0;
  display: flex;
  flex-direction: column;
  align-items: center;
}

input[type="file"] {
  display: none;
}

.custom-file-upload {
  border: 1px solid #888;
  border-radius: 0.5rem;
  padding: 5px;
  cursor: pointer;
  font-weight: bold;
  font-size: 14px;
}

.custom-file-upload:hover {
  background-color: #ccc;
}

.create-form {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

select {
  width: 105px;
  border: 1px solid #888;
  border-radius: 0.5rem;
  color: #2c3e50;
  font-weight: bold;
  font-size: 14px;
}

label {
  font-weight: bold;
}

.post-input {
  min-height: 55px;
  max-height: 150px;
  border: none;
  border-radius: 0.5rem;
  background-color: white;
  font-size: 14px;
  resize: vertical;
  padding: 15px;
  margin-bottom: 15px;
}

.post-input::placeholder {
  color: #888;
  font-style: italic;
}

#privacy {
  padding: 5px;
}

.break {
  margin: 25px 0;
}

.post-header {
  width: 90%;
  display: flex;
}

.post-link {
  align-self: center;
  text-decoration: none;
  color: inherit;
  margin-right: auto;
  margin-left: 10px;
}

.post-link:hover {
  text-decoration: underline;
  color: #888;
}

.post:hover {
  background-color: #e6e6e6;
  cursor: pointer;
}

.post-date,
.privacy-level {
  font-size: small;
  color: #888;
  margin-right: 10px;
  align-self: flex-end;
}

.post-content {
  background-color: white;
  border-radius: 5px;
  width: 90%;
  text-align: left;
  margin-bottom: 10px;
  padding: 0 20px;
}

.post-content p {
  margin-left: 5px;
}

.post img {
  max-width: 90%;
  height: auto;
  border-radius: 5px;
  margin-bottom: 15px;
}

.user-avatar {
  width: 30px;
  aspect-ratio: 1 / 1;
  clip-path: circle(50%);
  align-self: center;
  margin-bottom: 5px !important;
  margin-left: 5px;
}
</style>