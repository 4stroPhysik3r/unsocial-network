import { ref } from 'vue'

const getUserFromSession = () => {
   const user = ref(null)
   const errorSession = ref(null)

   const fetchUserDataFromSession = async () => {
      try {
         const response = await fetch("http://localhost:8000/api/my-profile", {
            method: "GET",
            credentials: "include",
         });
         if (response.ok) {
            const userData = await response.json();
            user.value = userData;
         } else {
            throw new Error("Failed to fetch user data");
         }
      } catch (errorSession) {
         errorSession.value = errorSession.message;
         console.errorSession("Error fetching user data:", errorSession);
      }
   };

   return { user, errorSession, fetchUserDataFromSession }
}

export default getUserFromSession
