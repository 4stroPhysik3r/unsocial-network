import { ref } from "vue";

export function getViewerStatus() {
  const status = ref([]);
  const errorViewer = ref(null);

  const loadViewer = async (postID) => {
    errorViewer.value = null;

    try {
      const response = await fetch(
        `http://localhost:8000/api/viewer-status/${postID}`,
        {
          method: "GET",
          credentials: "include",
        }
      );
      if (!response.ok) {
        throw new Error("Failed to fetch viewer status");
      }
      status.value = await response.json();
    } catch (err) {
      errorViewer.value = err.message;
    }
  };

  return { status, errorViewer, loadViewer };
}

export default getViewerStatus;
