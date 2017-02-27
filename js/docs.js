/*
 *
 * swapStyleSheet*********************************************************************
 *
 */

// Cookie functions
function createCookie(name,value,days) {
    var expires = "";
    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + (days*24*60*60*1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + value + expires + "; path=/";
}

function readCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
    for(var i=0;i < ca.length;i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1,c.length);
        if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
    }
    return null;
}

function eraseCookie(name) {
    createCookie(name,"",-1);
}
if (readCookie("night") == "true") {
  document.getElementById('pagestyle').setAttribute('href', '/css/style-alt.css');
  $('#switch-style').prop('checked', true);
} else {
  document.getElementById('pagestyle').setAttribute('href', '/css/style.css');
  $('#switch-style').prop('checked', false);
}
/*
 *
 * toggle menu *********************************************************************
 *
 */

$("#menu-toggle").click(function(e) {
        e.preventDefault();
        $(".wrapper").toggleClass("right-open");
    });


var navHeight = $('.navbar').outerHeight(true) + 80;

$(document.body).scrollspy({
	target: '#leftCol',
	offset: navHeight
});


$(document).ready(function(){
  // Add smooth scrolling to all links
  // $( ".toc-nav a" ).addClass( "active" );
  $(".toc-nav a").on('click', function(event) {
    // $(this).addClass('active');
    // Make sure this.hash has a value before overriding default behavior
    if (this.hash !== "") {
      // Prevent default anchor click behavior
      event.preventDefault();

      // Store hash
      var hash = this.hash;

      // Using jQuery's animate() method to add smooth page scroll
      // The optional number (800) specifies the number of milliseconds it takes to scroll to the specified area
      $('html, body').animate({
        scrollTop: $(hash).offset().top-80
      }, 800, function(){

        // Add hash (#) to URL when done scrolling (default click behavior)
        window.location.hash = hash;
      });
    } // End if
  });
});


$(document).ready(function(){
  // Add smooth scrolling to all links
  $(".nav-sidebar ul li a").on('click', function(event) {

    // Make sure this.hash has a value before overriding default behavior
    if (this.hash !== "") {
      // Prevent default anchor click behavior
      event.preventDefault();

      // Store hash
      var hash = this.hash;

      // Using jQuery's animate() method to add smooth page scroll
      // The optional number (800) specifies the number of milliseconds it takes to scroll to the specified area
      $('html, body').animate({
        scrollTop: $(hash).offset().top-80
      }, 800, function(){

        // Add hash (#) to URL when done scrolling (default click behavior)
        window.location.hash = hash;
      });
    } // End if
  });
});



/*
 *
 * make dropdown show on hover *********************************************************************
 *
 */

$('ul.nav li.dropdown').hover(function() {
  $(this).find('.dropdown-menu').stop(true, true).delay(200).fadeIn(500);
}, function() {
  $(this).find('.dropdown-menu').stop(true, true).delay(200).fadeOut(500);
});

/*
 *
 * swapStyleSheet*********************************************************************
 *
 */

function swapStyleSheet(sheet) {
    document.getElementById('pagestyle').setAttribute('href', sheet);
}


$('#switch-style').change(function() {

    if ($(this).is(':checked')) {
        swapStyleSheet('/css/style-alt.css');
        createCookie("night",true,999)
    } else {
        swapStyleSheet('/css/style.css');
        eraseCookie("night")
    }
});


// $(function (){
//     $(".nav-secondary-tabs .search-form input[type=search]").focus(function() {
//         $(".tabs").fadeOut();
//     }).blur(function() {
//         $(".tabs").fadeIn();
//     });
// })
//
// $(function (){
//     $(".nav-secondary input[type=search]").focus(function() {
//         $(".tabs").fadeOut();
//     }).blur(function() {
//         $(".tabs").fadeIn();
//     });
// })

/*
 *
 * TEMP HACK For side menu*********************************************************************
 *
 */

$('.nav-sidebar ul li a').click(function() {
    $(this).addClass('collapse').siblings().toggleClass('in');
});

if($('.nav-sidebar ul a.active').length != 0)
{
  $('.nav-sidebar ul').click(function() {
      $(this).addClass('collapse in').siblings;
  });
}

// $( ".nav-secondary" ).fadeIn( 500 );
$( ".page-content-wrapper" ).fadeIn( 500 );
