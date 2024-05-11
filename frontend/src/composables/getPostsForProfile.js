import { ref } from 'vue';

const getPostsForProfile = (id = null) => {
   const posts = ref(null);
   const postsError = ref(null);

   const loadPosts = async () => {
      try {
         let url = 'http://localhost:8000/api/posts';
         if (id) {
            url = `http://localhost:8000/api/postsFromID/${id}`;
         }
         const response = await fetch(url, {
            credentials: "include",
         });
         if (!response.ok) {
            throw Error(id ? 'Post does not exist' : 'Failed to fetch posts');
         }
         posts.value = await response.json();
      }
      catch (err) {
         postsError.value = err.message;
         console.error('Error fetching post:', postsError.value);
      }
   };

   return { posts, postsError, loadPosts };
};

export default getPostsForProfile;
