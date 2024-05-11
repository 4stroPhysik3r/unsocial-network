import { ref } from 'vue'

const getUsers = () => {
   const users = ref([])
   const error = ref(null)

   const loadUsers = async () => {
      try {
         let data = await fetch('http://localhost:8000/api/get-users', {
            credentials: "include",
         });
         if (!data.ok) {
            throw Error('No users found')
         }
         users.value = await data.json()
      }
      catch (err) {
         error.value = err.message
      }
   }

   return { users, error, loadUsers }
}

export default getUsers