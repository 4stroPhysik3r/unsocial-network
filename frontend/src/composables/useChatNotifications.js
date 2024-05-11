import { ref, onMounted, onUnmounted } from "vue";

export function useChatNotifications(webSocketUrl) {
  const isConnected = ref(false);
  const unreadMessages = ref([]);
  let webSocket = null;

  const handleIncomingMessage = async (unreadMessageIds) => {
    unreadMessages.value = unreadMessageIds;
  };

  // Setup WebSocket connection
  onMounted(() => {
    webSocket = new WebSocket(webSocketUrl);

    webSocket.onopen = () => {
      isConnected.value = true;
    };

    webSocket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      handleIncomingMessage(message);
    };

    webSocket.onerror = function (event) {
      console.error("WebSocket error observed:", event);
    };

    webSocket.onclose = () => {
      isConnected.value = false;
    };
  });

  onUnmounted(() => {
    if (webSocket) {
      webSocket.close();
    }
  });

  return {
    isConnected,
    unreadMessages,
    handleIncomingMessage
  };
}
