import { ref } from 'vue';

const getFollowStatus = () => {
  const followStatus = ref('');
  const errorStatus = ref(null);

  const fetchFollowStatus = async (userID) => {
    try {
      const response = await fetch(`http://localhost:8000/api/get-follow-status/${userID}`, {
        method: "GET",
        credentials: "include",
      });

      if (!response.ok) {
        throw Error('Failed to fetch follow status');
      }
      const data = await response.json();
      followStatus.value = data;
    } catch (err) {
      errorStatus.value = err.message;
    }
  };

  return { followStatus, errorStatus, fetchFollowStatus };
};

export default getFollowStatus;