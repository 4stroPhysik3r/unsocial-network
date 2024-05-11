import { ref } from 'vue';

const getPostFromPostID = () => {
  const post = ref(null);
  const errorPost = ref(null);

  const fetchPostFromID = async (postID) => {
    try {
      let data = await fetch(`http://localhost:8000/api/get-post/${postID}`, {
        credentials: 'include',
      });

      if (!data.ok) {
        throw Error('Failed to fetch post');
      }

      post.value = await data.json();
    } catch (err) {
      errorPost.value = err.message;
      console.error(errorPost.value);
    }
  };

  return { post, errorPost, fetchPostFromID };
};

export default getPostFromPostID;
