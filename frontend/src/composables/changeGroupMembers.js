import { ref } from "vue";

const joinGroup = async (groupID, updateMembershipStatus) => {
  try {
    const response = await fetch(
      `http://localhost:8000/api/join-group/${groupID}`,
      {
        credentials: "include",
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    if (!response.ok) {
      throw new Error("Failed to send join request");
    }
    updateMembershipStatus("request");

  } catch (error) {
    console.error("Error sending join request:", error);
  }
};

const leaveGroup = async (groupID, updateMembershipStatus) => {
  try {
    const response = await fetch(
      `http://localhost:8000/api/leave-group/${groupID}`,
      {
        credentials: "include",
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    if (!response.ok) {
      throw new Error("Failed to send leave request");
    }
    updateMembershipStatus("not_member");

  } catch (error) {
    console.error("Error sending leave request:", error);
  }
};

const inviteMembers = async (groupID, selectedMembers) => {
  try {
    const response = await fetch(`http://localhost:8000/api/invite-to-group/${groupID}`, {
      credentials: "include",
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        selectedMembers: selectedMembers.map(user => user.user_id),
      }),
    });

    if (!response.ok) {
      throw new Error("Failed to invite users to the group");
    }

  } catch (error) {
    console.error("Error inviting users to the group:", error);
  }
};

export { joinGroup, leaveGroup, inviteMembers };
