import { ref } from "vue";

export function getFollower() {
  const follower = ref([]);
  const errorFollower = ref(null);

  const loadFollower = async (userID) => {
    errorFollower.value = null;

    try {
      const response = await fetch(
        `http://localhost:8000/api/follower/${userID}`,
        {
          method: "GET",
          credentials: "include",
        }
      );
      if (!response.ok) {
        throw new Error("Failed to fetch follower");
      }
      follower.value = await response.json();
    } catch (err) {
      errorFollower.value = err.message;
    }
  };

  return { follower, errorFollower, loadFollower };
}

export default getFollower;
