// Define the tour!
var tour = {
  id: "hello-hopscotch",
  steps: [
    {
      title: "Navigation improvements!",
      content: "We've improved the navigation for the Docker documentation!<br /><br />This short tour takes less than a minute. It highlights the changes, then returns you to the current page. You'll learn about:<br /><ul><li>New top navigation</li><li>Left-hand navigation</li><li>Feedback links</li><li>In-page navigation</li></ul>",
      target: "main-content",
      placement: "top",
      xOffset: "center",
      yOffset: "300px",
      onShow: function () {
        // Hide the arrow on the first navigation bubble
        $('.hopscotch-bubble-arrow-container').css('visibility', 'hidden');
      }
    },
    {
      title: "Top Navigation",
      content: "Use the top navigation to discover the types of content available.",
      target: "top-nav",
      placement: "bottom",
      onShow: function () {
        // Show the arrow again
        $('.hopscotch-bubble-arrow-container').css('visibility', 'visible');
      }
    },
    {
      title: "Left Navigation",
      content: "Use the left navigation for a structured view of content within a top-level category.",
      target: "left-nav",
      placement: "right"
    },
    {
      title: "Feedback Links",
      content: "Use the feedback links to edit the page, provide feedback, or find out how to get support.",
      target: "feedback-links",
      placement: "left",
      multipage: "true",
      onNext: function() {
        window.location = "/learn/index.html";
      }
    },
    {
      title: "In-page navigation",
      content: "Use the in-page navigation links to jump to specific areas within the page you are viewing.",
      target: "side-toc",
      placement: "left"
    }
  ],
  showPrevButton: true,
  scrollTopMargin: 200,
  skipIfNoElement: false,
  onEnd: function() {
    // Return them back where they came from when the tour ends
    if (hopscotch.getState() === null) {
      window.location = document.referrer;
    }
  },
  onClose: function() {
    // Return them back where they came from if they end the tour early
    if (hopscotch.getState() === "hello-hopscotch:4") {
      window.location = document.referrer;
    }
  }
};

// Start tour if button is pressed
$("#start-tour").click(function(){
  hopscotch.startTour(tour);
});

// Resume tour if already in progress
if (hopscotch.getState() === "hello-hopscotch:4") {
  hopscotch.startTour(tour);
}
