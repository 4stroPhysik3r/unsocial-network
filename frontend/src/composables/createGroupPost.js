import { ref } from 'vue';

const createGroupPost = () => {
   const content = ref("");
   const image = ref(null);
   const error = ref(null);

   const groupPost = async (groupID) => {
      try {
         const groupPostData = {
            content: content.value,
            post_image: image.value,
         };

         let data = await fetch(`http://localhost:8000/api/create-group-post/${groupID}`, {
            method: "POST",
            credentials: "include",
            headers: {
               "Content-Type": "application/json",
            },
            body: JSON.stringify(groupPostData),
         });

         if (!data.ok) {
            throw Error("Failed to create post");
         }

         content.value = "";
         image.value = null;

      } catch (error) {
         console.error("Error creating post:", error);
         error.value = error.message;
      }
   };

   return { groupPost, image, error, content };
};

export default createGroupPost;
