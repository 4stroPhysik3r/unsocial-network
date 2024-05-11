<template>
  <main>
    <h1>Notifications</h1>
    <div class="notification-container" v-if="notifications.length > 0">
      <div
        class="notification"
        v-for="notification in notifications"
        :key="notification.notification_id"
      >
        <p>{{ notification.message }}</p>
        <div v-if="notification.type === 'new_event'">
          <router-link
            :to="{
              name: 'GroupID',
              params: { groupID: notification.reference_id },
            }"
            @click="respondToNotification(notification.notification_id, true)"
            class="group-link"
          >
            Go to group page
          </router-link>
        </div>

        <div v-else>
          <button
            @click="respondToNotification(notification.notification_id, true)"
          >
            Accept
          </button>
          <button
            @click="respondToNotification(notification.notification_id, false)"
          >
            Reject
          </button>
        </div>
      </div>
    </div>

    <div v-else>
      <p>No new notifications</p>
    </div>
  </main>
</template>

<script>
import { useNotifications } from "@/composables/useNotifications";

export default {
  name: "Notifications",
  setup() {
    const { notifications, respondToNotification } = useNotifications();

    return {
      notifications,
      respondToNotification,
    };
  },
};
</script>

<style scoped>
.notification-container {
  display: grid;
  width: 60vw;
  max-width: 1200px;
  justify-content: center;
  border-radius: 1rem;
  background-color: #f0f0f0;
  margin-left: auto;
  margin-right: auto;
}

.notification {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 20px;
  padding: 0 20px;
}

.notification-dropdown {
  position: absolute;
  right: 0;
  background-color: #f9f9f9;
  min-width: 160px;
  box-shadow: 0px 8px 16px 0px rgba(0, 0, 0, 0.2);
  z-index: 1;
}

.notification-item {
  padding: 12px 16px;
  text-align: left;
  display: block;
  white-space: nowrap;
}

.notification-item:hover {
  background-color: #f1f1f1;
}

button {
  margin: 0;
  min-width: fit-content;
}

button:first-child {
  margin-right: 10px;
}

.group-link {
  text-decoration: none;
  color: #42b983;
}
</style>