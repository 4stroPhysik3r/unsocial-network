import { ref } from 'vue'

const getUserFromID = () => {
    const user = ref(null)
    const error = ref(null)

    const fetchUserDataFromID = async (userID) => {
        try {
            const response = await fetch(`http://localhost:8000/api/userid/${userID}`, {
                method: "GET",
                credentials: "include",
            })
            if (response.ok) {
                const userData = await response.json()
                user.value = userData
            } else {
                throw new Error("Failed to fetch user data")
            }
        } catch (error) {
            error.value = error.message
            console.error("Error fetching user data:", error)
        }
    }

    return { user, error, fetchUserDataFromID }
}

export default getUserFromID