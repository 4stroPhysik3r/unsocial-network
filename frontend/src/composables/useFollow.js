import { ref } from 'vue';

const useFollow = (userID) => {
  const followStatus = ref(""); // Possible values: 'rejected', 'accepted', 'pending', 'not_following'
  const errorStatus = ref(null);

  const updateFollowStatus = (status) => {
    followStatus.value = status;
  };

  const fetchFollowStatus = async () => {
    try {
      const response = await fetch(`http://localhost:8000/api/get-follow-status/${userID}`, {
        method: "GET",
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error('Failed to fetch follow status');
      }
      const data = await response.json();


      updateFollowStatus(data);

    } catch (err) {
      errorStatus.value = err.message;
    }
  };

  const follow = async () => {
    try {
      const response = await fetch(`http://localhost:8000/api/follow-user/${userID}`, {
        credentials: "include",
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error("Failed to send follow request");
      }

      await fetchFollowStatus()

    } catch (error) {
      console.error("Error sending follow request:", error);
    }
  };

  const unfollow = async () => {
    try {
      const response = await fetch(`http://localhost:8000/api/unfollow-user/${userID}`, {
        credentials: "include",
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error("Failed to send unfollow request");
      }

      await fetchFollowStatus()

    } catch (error) {
      console.error("Error sending unfollow request:", error);
    }
  };

  return { followStatus, errorStatus, fetchFollowStatus, follow, unfollow, updateFollowStatus };
};

export default useFollow;
