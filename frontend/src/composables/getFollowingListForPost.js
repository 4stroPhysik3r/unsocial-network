import { ref } from "vue";

export function getFollowingListForPost() {
  const following = ref([]);
  const errorFollowingList = ref(null);

  const loadFollowing = async () => {
    errorFollowingList.value = null;

    try {
      const response = await fetch(
        "http://localhost:8000/api/following-list/",
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
      errorFollowingList.value = err.message;
    }
  };

  return { following, errorFollowingList, loadFollowing };
}

export default getFollowingListForPost;
