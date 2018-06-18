// /v1/cache with basic auth
    // then pass token
    $(document).ready(function () {
        generateCards();
      });
  
      $("#search-box").change(function () {
        alert("Handler for .change on searchbox called.");
        searchFilter();
      });
  
      function generateCards() {
        
      }
  
  
      function displayLogin() {
        alert("login")
      }
  
      function unixTimeToJSDate(unixTime) {
        // Create a new JavaScript Date object based on the timestamp
        return new Date(unixTime * 1000);
      }
  