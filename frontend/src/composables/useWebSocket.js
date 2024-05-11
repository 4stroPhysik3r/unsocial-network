import { ref, onMounted, onUnmounted } from "vue";

export function useWebSocket(url, onMessageReceived, autoConnect) {
  const isConnected = ref(false);
  let ws = null;

  // Setup WebSocket connection
  onMounted(() => {
    if (autoConnect && !ws) {
      ws = new WebSocket(url);

      ws.onopen = () => {
        isConnected.value = true;
      };

      ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        if (typeof onMessageReceived === "function") {
          onMessageReceived(message); // Invoke the callback with the received message
        }
      };

      ws.onerror = function (event) {
        console.error("WebSocket error observed:", event);
      };

      ws.onclose = () => {
        isConnected.value = false;
      };
    }
  });

  const connect = () => {
    if (!ws) {
      ws = new WebSocket(url);
      ws.onopen = () => {
        isConnected.value = true;
      };

      ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        if (typeof onMessageReceived === "function") {
          onMessageReceived(message); // Invoke the callback with the received message
        }
      };

      ws.onerror = function (event) {
        console.error("WebSocket error observed:", event);
      };

      ws.onclose = () => {
        isConnected.value = false;
      };
    }
  };

  // Function to send a message through the WebSocket connection
  const sendMessage = (message) => {
    if (ws && isConnected.value) {
      ws.send(JSON.stringify(message));
    } else {
      console.error("WebSocket not connected. Message not sent:", message);
    }
  };

  // Cleanup on component unmount
  onUnmounted(() => {
    if (ws) {
      ws.close();
    }
  });

  return { isConnected, sendMessage, connect };
}
