const searchUsers = () => {
   const searchFunction = () => {
      var filter = document.getElementById("myInput").value.toUpperCase();
      var options = document.getElementById("options");
      var links = document
         .getElementById("userSearch")
         .getElementsByTagName("a");

      // Show search container if there's input or if the input field is focused, otherwise hide it
      options.style.display =
         filter || document.activeElement === document.getElementById("myInput")
            ? "block"
            : "none";

      for (var i = 0; i < links.length; i++) {
         links[i].style.display = links[i].textContent
            .toUpperCase()
            .includes(filter)
            ? ""
            : "none";
      }
   };

   document.body.addEventListener("click", (event) => {
      let input = document.getElementById("myInput");
      let options = document.getElementById("options");

      if (input === null) input = document.createElement("myInput");
      if (options === null) options = document.createElement("options");

      // Check if input field is not active
      if (
         event.target !== input &&
         !input.contains(event.target) &&
         document.activeElement !== input
      ) {
         options.style.display = "none";
      }
   });

   const clearSearchInput = () => {
      document.getElementById("myInput").value = ""; // Clear the input field

      // Hide search container
      var options = document.getElementById("options");
      options.style.display = "none";

   };

   return { searchFunction, clearSearchInput }
}

export default searchUsers;
