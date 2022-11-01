// Get the modal for image 1 
var modal = document.getElementById("myModal");

// Get the image and insert it inside the modal - use its "alt" text as a caption
var img = document.getElementById("myImg");
var modalImg = document.getElementById("img01");
var captionText = document.getElementById("caption");
img.onclick = function(){
  modal.style.display = "block";
  modalImg.src = this.src;
  captionText.innerHTML = this.alt;
}

// Get the element that closes the modal
var span = document.getElementsByClassName("close")[0];

// When the user clicks on (x), close the modal
span.onclick = function(){
  modal.style.display = "none";
}

// Get the modal for image 2 
var modal = document.getElementById("myModal2");

// Get the image and insert it inside the modal - use its "alt" text as a caption
var img = document.getElementById("myImg2");
var modalImg = document.getElementById("img02");
var captionText = document.getElementById("caption2");
img.onclick = function(){
  modal.style.display = "block";
  modalImg.src = this.src;
  captionText.innerHTML = this.alt;
}

// Get the element that closes the modal
var span = document.getElementsByClassName("close")[0];

// When the user clicks on (x), close the modal
span.onclick = function(){
  modal.style.display = "none";
}

// Get the modal for image 3 
var modal = document.getElementById("myModal3");

// Get the image and insert it inside the modal - use its "alt" text as a caption
var img = document.getElementById("myImg3");
var modalImg = document.getElementById("img03");
var captionText = document.getElementById("caption3");
img.onclick = function(){
  modal.style.display = "block";
  modalImg.src = this.src;
  captionText.innerHTML = this.alt;
}

// Get the element that closes the modal
var span = document.getElementsByClassName("close")[0];

// When the user clicks on (x), close the modal
span.onclick = function(){
  modal.style.display = "none";
}

// Get the modal for image 4 
var modal = document.getElementById("myModal4");

// Get the image and insert it inside the modal - use its "alt" text as a caption
var img = document.getElementById("myImg4");
var modalImg = document.getElementById("img04");
var captionText = document.getElementById("caption4");
img.onclick = function(){
  modal.style.display = "block";
  modalImg.src = this.src;
  captionText.innerHTML = this.alt;
}


// Get the element that closes the modal
var span = document.getElementsByClassName("close")[0];

// When the user clicks on (x), close the modal
span.onclick = function(){
  modal.style.display = "none";
}

