
import { ref } from 'vue';

const createEvent = () => {
  const eventTitle = ref('');
  const eventContent = ref('');
  const eventDate = ref('');

  const createNewEvent = async (groupID) => {

    const parsedGroupID = parseInt(groupID, 10);
    try {
      const response = await fetch(`http://localhost:8000/api/create-event/${groupID}`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: eventTitle.value,
          content: eventContent.value,
          date: eventDate.value,
          group_id: parsedGroupID
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to create event');
      }

      eventTitle.value = '';
      eventContent.value = '';
      eventDate.value = '';

    } catch (error) {
      console.error('Error creating event:', error);
    }
  };

  return {
    eventTitle,
    eventContent,
    eventDate,
    createNewEvent,
  };
};

const fetchEventData = () => {
  const events = ref([]);

  const getEventData = async (groupID) => {
    try {
      const response = await fetch(`http://localhost:8000/api/get-events/${groupID}`, {
        credentials: 'include',
      });
      if (!response.ok) {
        throw new Error('Failed to fetch events');
      }
      const eventData = await response.json();

      // Fetch attendee status for each event
      const eventsWithStatus = await Promise.all(eventData.map(async (event) => {
        try {
          const attendeeStatusResponse = await fetchAttendeesStatus(event.event_id);
          return { ...event, attendeeStatus: attendeeStatusResponse };
        } catch (error) {
          console.error('Error fetching attendee status:', error);
          return { ...event, attendeeStatus: null };
        }
      }));

      events.value = eventsWithStatus;
    } catch (error) {
      console.error('Error fetching events:', error);
    }
  };

  return {
    events,
    getEventData,
  };
};

const fetchAttendeesStatus = async (eventId) => {
  try {
    const response = await fetch(`http://localhost:8000/api/get-attendees-status/${eventId}`, {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch attendees status');
    }

    const data = await response.json();
    return data;

  } catch (error) {
    console.error('Error fetching attendees status:', error);
    return null;
  }
};

const updateAttendeesStatus = async (eventId, attendeesStatus) => {
  const statusString = String(attendeesStatus)
  try {
    const response = await fetch(`http://localhost:8000/api/update-attendees-status/${eventId}`, {
      method: 'PUT',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ attendees_status: statusString }),
    });

    if (!response.ok) {
      throw new Error('Failed to update attendees status');
    }

  } catch (error) {
    console.error('Error updating attendees status:', error);
  }
};

export { createEvent, fetchEventData, fetchAttendeesStatus, updateAttendeesStatus };
