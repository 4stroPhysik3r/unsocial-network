import { ref } from 'vue';

const getCommentsForPost = () => {
  const comments = ref([]);
  const error = ref(null);

  const fetchCommentsForPost = async (postID) => {
    try {
      const response = await fetch(`http://localhost:8000/api/get-comments-for-post/${postID}`, {
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error('Failed to fetch comments');
      }

      comments.value = await response.json();
    } catch (err) {
      error.value = err.message;
    }
  };

  const addComment = async (postID, commentData) => {
    try {
      const response = await fetch(`http://localhost:8000/api/add-comment/${postID}`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(commentData),
      });

      if (!response.ok) {
        throw new Error('Failed to add comment');
      }

    } catch (err) {
      error.value = err.message;
    }
  };

  return { comments, error, fetchCommentsForPost, addComment };
};

export default getCommentsForPost;
