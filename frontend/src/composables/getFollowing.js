import { ref } from "vue";

export function getFollowing() {
  const following = ref([]);
  const errorFollowing = ref(null);

  const loadFollowing = async (userID) => {
    errorFollowing.value = null;

    try {
      const response = await fetch(
        `http://localhost:8000/api/following/${userID}`,
        {
          method: "GET",
          credentials: "include",
        }
      );
      if (!response.ok) {
        throw new Error("Failed to fetch following");
      }
      following.value = await response.json();
    } catch (err) {
      errorFollowing.value = err.message;
    }
  };

  return { following, errorFollowing, loadFollowing };
}

export default getFollowing;
