import { ref } from "vue";

const getGroups = () => {
  const myGroups = ref([]);
  const allGroups = ref([]);
  const errorGroups = ref(null);

  const loadGroups = async () => {
    try {
      // Fetch my groups
      let myData = await fetch("http://localhost:8000/api/get-my-groups", {
        credentials: "include",
      });
      if (!myData.ok) {
        throw Error("Failed to fetch my groups");
      }
      myGroups.value = await myData.json();

      // Fetch all groups
      let allData = await fetch("http://localhost:8000/api/get-all-groups", {
        credentials: "include",
      });
      if (!allData.ok) {
        throw Error("Failed to fetch all groups");
      }
      allGroups.value = await allData.json();
    } catch (err) {
      errorGroups.value = err.message;
    }
  };

  return { myGroups, allGroups, errorGroups, loadGroups };
};

export default getGroups;
