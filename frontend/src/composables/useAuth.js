import { ref } from "vue";

const isAuthenticated = ref(false);

async function checkAuthStatus() {
  try {
    const response = await fetch("http://localhost:8000/api/auth/status", {
      credentials: "include",
    });

    if (response.ok) {
      const data = await response.json();
      isAuthenticated.value = data.isAuthenticated;
    } else {
      isAuthenticated.value = false;
    }
  } catch (error) {
    console.error("Error checking authentication status", error);
    isAuthenticated.value = false;
  }
}

async function logout() {
  await fetch("http://localhost:8000/api/logout", {
    method: "POST",
    credentials: "include",
  });

  isAuthenticated.value = false;
}

export function useAuth() {
  return { isAuthenticated, checkAuthStatus, logout };
}
