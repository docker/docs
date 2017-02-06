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
      arrowOffset: "center",
      width: "570px",
      onShow: function () {
        // Show the arrow again
        $('.hopscotch-bubble-arrow-container').css('visibility', 'visible');
      }
    },
    {
      title: "Guides",
      content: "Use the <b>Guides</b> tab to learn how to install, configure, and manage Docker as a whole, or to view the docs archives for previous Docker versions.",
      target: "top-nav",
      placement: "bottom",
      width: "570px"
    },
    {
      title: "Product Manuals",
      content: "Use the <b>Product Manuals</b> tab to learn detailed information about a specific Docker product, such as Docker Cloud or UCP.",
      target: "top-nav",
      placement: "bottom",
      width: "570px",
      arrowOffset: "140px"
    },
    {
      title: "Glossary",
      content: "Use the <b>Glossary</b> tab to quickly define and learn about terminology specific to Docker.",
      target: "top-nav",
      placement: "bottom",
      width: "570px",
      arrowOffset: "280px"
    },
    {
      title: "Reference",
      content: "Use the <b>Reference</b> tab to go straight to reference information about Docker, such as API and CLI reference topics.",
      target: "top-nav",
      placement: "bottom",
      width: "570px",
      arrowOffset: "390px"
    },
    {
      title: "Samples",
      content: "Use the <b>Samples</b> tab to learn about Docker using self-paced tutorials, labs, and sample Docker applications.",
      target: "top-nav",
      placement: "bottom",
      width: "570px",
      arrowOffset: "490px"
    },
    {
      title: "Left Navigation",
      content: "Use the left navigation for a structured view of content within a top-level category.",
      target: "left-nav",
      placement: "right",
      yOffset: "100px",
      arrowOffset: "center",
    },
    {
      title: "Feedback Links",
      content: "Use the feedback links to edit the page, provide feedback, or find out how to get support.",
      target: "feedback-links",
      placement: "left",
      arrowOffset: "center",
      multipage: "true",
      onNext: function() {
        window.location = "/learn/index.html";
      }
    },
    {
      title: "In-page navigation",
      content: "Use the in-page navigation links to jump to specific areas within the page you are viewing.",
      target: "side-toc",
      placement: "left",
      arrowOffset: "center"
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
    if (hopscotch.getState() === "hello-hopscotch:9") {
      window.location = document.referrer;
    }
  }
};

// Start tour if button is pressed
$("#start-tour").click(function(){
  hopscotch.startTour(tour);
});

// Resume tour if already in progress
if (hopscotch.getState() === "hello-hopscotch:9") {
  hopscotch.startTour(tour);
}
