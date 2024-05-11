<template>
  <main>
    <h1>Chats</h1>

    <div class="chat-container">
      <div v-if="chatNamesIDs && chatNamesIDs.length" class="chat-nav">
        <div
          v-for="chat in chatNamesIDs"
          :key="chat.chat_id"
          @click="selectChat(chat.chat_id)"
          :class="{
            'chat-item': true,
            'selected-chat-item': selectedChatId === chat.chat_id,
          }"
        >
          <p>{{ chat.full_name }}</p>
          <span v-if="hasUnreadMessages(chat.chat_id)" class="red-dot"></span>
        </div>
      </div>

      <div v-else class="chat-nav chat-nav-error">
        <p>
          You have no chats to see. Start following someone or join a group.
        </p>
      </div>

      <div v-if="selectedChatId" class="chat">
        <div
          v-if="allMessages && allMessages.length"
          class="chat-content"
          :style="{ maxHeight: '480px' }"
        >
          <div
            v-for="message in allMessages"
            :key="message.message_id"
            :class="{
              'send-messages': message.sender_id === sessionUserID.userID,
              'receive-messages': message.sender_id !== sessionUserID.userID,
            }"
          >
            <p class="message-content">
              <span class="chat-header">
                {{ getSenderFullName(message.sender_id) }},
                {{ formatDateTime(message.created_at) }}
              </span>
              <br />
              {{ message.content }}
            </p>
          </div>
        </div>

        <div v-else>
          <p class="no-messages">
            No messages in this chat. Start the conversation!
          </p>
        </div>

        <div class="chat-input-container">
          <textarea
            v-model="chatInput"
            class="chat-input"
            placeholder="Type message..."
            @keyup.enter.prevent="send"
          ></textarea>

          <button class="message-send-btn" @click="send">Send</button>
        </div>
      </div>
    </div>
  </main>
</template>

<script>
import { ref, computed, onMounted } from "vue";
import getChatNamesIDs from "@/composables/getChatNamesIDs";
import openChat from "@/composables/openChat";
import getUserFromSession from "@/composables/getUserFromSession";
import getUsers from "@/composables/getUsers";
import { useChatNotifications } from "@/composables/useChatNotifications";

export default {
  setup() {
    const { users, loadUsers } = getUsers();
    const { chatNamesIDs, loadChatNamesIDs } = getChatNamesIDs();
    const { user: sessionUserID, fetchUserDataFromSession } =
      getUserFromSession();
    const { messages, fetchMessages, sendMessage } = openChat();
    const selectedChatId = ref(null);
    const chatInput = ref("");
    const { isConnected, unreadMessages } = useChatNotifications(
      "ws://localhost:8000/api/chat-notifications/ws"
    );

    onMounted(async () => {
      await loadUsers();
      await fetchUserDataFromSession();
      await loadChatNamesIDs();
    });

    const allMessages = computed(() => {
      if (!messages.value) {
        return [];
      }
      return messages.value
        .slice()
        .sort((a, b) => new Date(a.created_at) - new Date(b.created_at));
    });

    const send = () => {
      if (chatInput.value.trim() !== "") {
        sendMessage(sessionUserID.value.userID, chatInput.value);
        chatInput.value = "";
      }
    };

    const hasUnreadMessages = (chatId) => {
      if (unreadMessages.value) {
        return unreadMessages.value.includes(chatId);
      }
      return false;
    };

    const formatDateTime = (createdAt) => {
      const chatDate = new Date(createdAt);
      const time = chatDate.toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });
      const day = ("0" + chatDate.getDate()).slice(-2);
      const month = ("0" + (chatDate.getMonth() + 1)).slice(-2);
      const year = ("0" + chatDate.getFullYear()).slice(-2);
      return `${time}, ${day}.${month}.${year}`;
    };

    const getSenderFullName = (senderId) => {
      const user = users.value.find((user) => user.user_id === senderId);
      return user ? user.firstname + " " + user.lastname : "Unknown";
    };

    const selectChat = (chatId) => {
      selectedChatId.value = chatId;
      fetchMessages(chatId);
    };

    return {
      chatNamesIDs,
      selectedChatId,
      selectChat,
      fetchMessages: async (chatId) => {
        await fetchMessages(chatId);
        selectedChatId.value = chatId;
      },
      allMessages,
      sessionUserID,
      sendMessage,
      chatInput,
      send,
      isConnected,
      hasUnreadMessages,
      getSenderFullName,
      formatDateTime,
    };
  },
};
</script>

<style>
.chat-container {
  display: flex;
}

.chat-nav {
  display: block;
  min-width: 210px;
  background-color: #f0f0f0;
  border-radius: 1rem;
  margin-right: 20px;
  text-align: left;
  padding: 20px;
  font-size: 18px;
}

.chat {
  display: block;
  width: 90%;
  padding: 20px;
  border-radius: 1rem;
  background-color: #f0f0f0;
}

.chat-content {
  display: flex;
  flex-direction: column;
  border-radius: 1rem;
  background-color: #fff;
  padding: 20px;
}

.message-send-btn {
  padding: 10px;
  margin: 10px;
}

.chat-input-container {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 20px;
}

.chat-input-container > * {
  flex: 1;
  margin-right: 0;
}

.chat-input {
  min-width: 80%;
  min-height: 15px;
  border-radius: 0.5rem;
  border: none;
  padding: 10px;
  resize: vertical;
}

.chat-input::placeholder {
  font-style: italic;
  font-size: 14px;
  margin-top: 100px;
}

.chat-header {
  font-style: italic;
  font-size: small;
}

.message-content {
  font-size: 18px;
  margin: 0;
  margin-bottom: 10px;
  padding: 5px 20px;
}

.messages {
  width: 20em;
  border-radius: 1rem;
  background-color: #f0f0f0;
  margin-bottom: 10px;
}

.receive-messages {
  text-align: left;
  align-self: flex-start;
}

.no-messages {
  text-align: left;
  align-self: flex-start;
  margin-left: 75px;
}

.send-messages {
  text-align: right;
  align-self: flex-end;
}

.chat-item {
  cursor: pointer;
  padding: 10px;
  margin: 5px 0;
  border: 1px solid #ddd;
  border-radius: 0.5rem;
  transition: background-color 0.3s, box-shadow 0.3s;
  display: flex;
}

.chat-item:hover {
  background-color: #ddd;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.chat-item p {
  margin: 0;
  font-size: 1rem;
}

.selected-chat-item {
  background-color: #ccc;
}

.chat-nav-error {
  text-transform: lowercase;
  font-size: medium;
  font-weight: 500;
}

.red-dot {
  width: 10px;
  height: 10px;
  background-color: #c92a2a;
  border-radius: 50%;
  display: inline-block;
  margin-left: auto;
  margin-right: 5px;
  align-self: center;
}

.chat-content {
  overflow-y: auto;
}
</style>