import { ref } from 'vue';

const getGroupInfo = () => {
  const groupInfo = ref(null);
  const errorInfo = ref(null);

  const fetchGroupInfo = async (groupID) => {
    try {
      const response = await fetch(`http://localhost:8000/api/get-group-info/${groupID}`, {
        credentials: 'include',
      });

      if (!response.ok) {
        throw Error('Failed to fetch group information');
      }

      groupInfo.value = await response.json();
    } catch (error) {
      errorInfo.value = error.message;
      console.error("error fetching group info: ", errorInfo.value);
    }
  };


  const groupPosts = ref(null);
  const errorPost = ref(null);

  const fetchGroupPosts = async (groupID) => {
    try {
      let data = await fetch(`http://localhost:8000/api/get-group-posts/${groupID}`, {
        credentials: 'include',
      });

      if (!data.ok) {
        throw Error('Failed to fetch post');
      }

      groupPosts.value = await data.json();
    } catch (err) {
      errorPost.value = err.message;
      console.error("error fetching group posts:", errorPost.value);
    }
  };

  const isMember = ref("");
  const checkMembershipStatus = async (groupID) => {
    try {
      const response = await fetch(`http://localhost:8000/api/check-membership/${groupID}`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (response.ok) {
        const data = await response.json();
        isMember.value = data;
      } else {
        throw new Error('Failed to check membership status');
      }
    } catch (error) {
      console.error('Error checking membership status:', error);
    }
  };
  function updateMembershipStatus(newStatus) {
    isMember.value = newStatus;
  }

  return { groupInfo, errorInfo, fetchGroupInfo, groupPosts, errorPost, fetchGroupPosts, isMember, checkMembershipStatus, updateMembershipStatus };
};

export default getGroupInfo;
