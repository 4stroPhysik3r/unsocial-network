import { ref, reactive, onUnmounted } from "vue";

const openChat = () => {
  const messages = ref([]);
  const chatID = ref(null);
  let socket = null;
  const state = reactive({
    isLoading: false,
    error: null,
  });

  async function fetchMessages(chatId) {
    state.isLoading = true;
    chatID.value = chatId;
    state.error = null;

    try {
      const response = await fetch("http://localhost:8000/api/get-messages", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ chat_id: chatId }),
      });

      if (!response.ok) {
        throw new Error("Failed to fetch messages");
      }

      const data = await response.json();
      if (data !== null) {
        messages.value = data;
      } else {
        messages.value = [];
      }

      // After fetching, establish a WebSocket connection for real-time updates
      connectWebSocket(chatId);
    } catch (error) {
      state.error = error.message;
      console.error("Fetch messages error:", error);
    } finally {
      state.isLoading = false;
    }
  }

  function connectWebSocket(chatId) {
    if (socket) {
      socket.close();
    }

    socket = new WebSocket("ws://localhost:8000/api/chat/ws");

    socket.onopen = () => {
      socket.send(JSON.stringify({ chat_id: chatId }));
    };

    socket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      if (message.chat_id === chatId) {
        messages.value.push(message);
      }
    };

    socket.onerror = (error) => {
      console.error("WebSocket Error:", error);
    };

    onUnmounted(() => {
      if (socket) {
        socket.close();
      }
    });
  }

  function sendMessage(userId, messageText) {
    if (socket) {
      const created_at = new Date().toISOString();
      socket.send(
        JSON.stringify({
          chat_id: chatID.value,
          sender_id: userId,
          content: messageText,
          created_at: created_at,
        })
      );
    }
  }

  function closeChat() {
    if (socket) {
      socket.close();
    }
  }

  return {
    messages,
    chatID,
    state,
    fetchMessages,
    sendMessage,
    closeChat,
  };
};

export default openChat;
