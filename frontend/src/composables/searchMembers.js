const searchMembers = () => {
   const searchFunction = () => {
      var filter = document.getElementById("myInput2").value.toUpperCase();
      var options = document.getElementById("options2");
      var links = document.getElementById("memberSearch").getElementsByTagName("p");

      // Show search container if there's input or if the input field is focused, otherwise hide it
      options.style.display =
         filter || document.activeElement === document.getElementById("myInput2")
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
      let input = document.getElementById("myInput2");
      let options = document.getElementById("options2");

      if (input === null) input = document.createElement("myInput2");
      if (options === null) options = document.createElement("options2");

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
      document.getElementById("myInput2").value = ""; // Clear the input field

      // Hide search container
      var options = document.getElementById("options2");
      options.style.display = "none";

   };

   return { searchFunction, clearSearchInput }
}

export default searchMembers;
