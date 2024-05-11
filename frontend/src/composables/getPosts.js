import { ref } from "vue"

const getPosts = () => {
  const posts = ref([])
  const errorPosts = ref(null)

  const loadPosts = async () => {
    try {
      let data = await fetch("http://localhost:8000/api/get-posts-feed", {
        credentials: "include",
      })
      if (!data.ok) {
        throw Error("no data available")
      }
      posts.value = await data.json()
    }
    catch (err) {
      errorPosts.value = err.message
    }
  }

  return { posts, errorPosts, loadPosts }
}

export default getPosts
