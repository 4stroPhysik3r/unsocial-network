import { ref } from "vue";

const getEvents = () => {
  const events = ref([]);
  const errorEvents = ref(null);

  const loadEvents = async () => {
    try {
      let data = await fetch("http://localhost:8000/api/get-events", {
        credentials: "include",
      });
      if (!data.ok) {
        throw Error("Failed to fetch all events");
      }
      events.value = await data.json();
    } catch (err) {
      errorEvents.value = err.message;
    }
  };

  return { events, errorEvents, loadEvents };
};

export default getEvents;
