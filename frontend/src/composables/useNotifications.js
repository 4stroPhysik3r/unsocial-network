import { ref } from "vue";
import { useWebSocket } from "./useWebSocket.js";

const notifications = ref([]);
const wsUrl = "ws://localhost:8000/api/notifications/ws";

const handleIncomingMessage = (messages) => {
  notifications.value = [...notifications.value, ...messages];
};

const { isConnected, sendMessage, connect } = useWebSocket(
  wsUrl,
  handleIncomingMessage,
  true
);

const respondToNotification = (notificationId, accepted) => {
  const response = {
    action: "notification_response",
    notification_id: notificationId,
    accepted,
  };
  sendMessage(response);

  // Remove the notification from the list
  notifications.value = notifications.value.filter(
    (notification) => notification.notification_id !== notificationId
  );
};

export function useNotifications() {

  return {
    isConnected,
    notifications,
    respondToNotification,
    handleIncomingMessage,
    connect,
  };
}
