// Get the modal for image 1 
var modal = document.getElementById("img-modal1");

// If modal exists, get the image and insert it inside the modal - use its "alt" text as a caption
if (modal) {
  const images = document.querySelectorAll("main img");
  var modalImg = document.getElementById("img-modal-img1");
  var captionText = document.getElementById("img-modal-caption1");
  function handleImageClick(event) {
    modal.style.display = "block";
    modalImg.src = event.target.src;
    modalImg.alt = event.target.alt;
    captionText.innerHTML = event.target.alt;
    window.addEventListener("keydown", handleModalClose)
  }
  images.forEach(image => image.addEventListener("click", handleImageClick))


  // Get the element that closes the modal
  var span = document.getElementById("img-modal-close1");
  function handleModalClose(event) {
    if (event.type==="click"||event.key==="Escape"){
      modal.style.display = "none"
      window.removeEventListener("keydown", handleModalClose)
    }
  }
  modal.addEventListener("click", handleModalClose)

  // When the user clicks on (x), close the modal
  span.onclick = function(){
    modal.style.display = "none";
  } 
}
