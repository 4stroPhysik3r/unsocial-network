import { ref } from "vue";

const getChatNamesIDs = () => {
  const chatNamesIDs = ref([]);
  const errorChatNamesIDs = ref(null);

  const loadChatNamesIDs = async () => {
    try {
      let response = await fetch("http://localhost:8000/api/get-chats", {
        credentials: "include",
      });
      if (!response.ok) {
        throw Error("no data available");
      }
      chatNamesIDs.value = await response.json();
    } catch (err) {
      errorChatNamesIDs.value = err.message;
    }
  };

  return { chatNamesIDs, errorChatNamesIDs, loadChatNamesIDs };
};

export default getChatNamesIDs;
