let slideIndex = 0;
let timeoutId = 0;
let showingSlides = false;
if (location.pathname == "/") {
    showSlides();
}

function clickShowSlides(plus) {
    clearTimeout(timeoutId)
    if(plus) {
        showSlides();
    } else {
        showSlides(false);
    }
}

function showSlides(plus = true) {
    showingSlides = true
    let i;
    let slides = document.getElementsByClassName("mySlides");
    let dots = document.getElementsByClassName("dot");
    for (i = 0; i < slides.length; i++) {
        slides[i].style.display = "none";
    }
    if (plus) {
        slideIndex++;
    } else {
        slideIndex--;
    }
    
    if (slideIndex > slides.length) { slideIndex = 1 }
    if (slideIndex <= 0) { slideIndex = slides.length }
    for (i = 0; i < dots.length; i++) {
        dots[i].className = dots[i].className.replace(" dot-active", "");
    }
    if (slides.length != 0) {
        slides[slideIndex - 1].style.display = "block";
        dots[slideIndex - 1].className += " dot-active";
        timeoutId = setTimeout(showSlides, 5000); // Change image every 2 seconds
    } else {
        clearTimeout(timeoutId)
    }
}

htmx.on('htmx:afterRequest', (evt) => {
    // check which element triggered the htmx request. If it's the one you want call the function you need
  //you have to add htmx: before the event ex: 'htmx:afterRequest'
  console.log(evt)
  
    if(evt.detail.pathInfo.requestPath == "/") {
        console.log(evt.detail.pathInfo.requestPath)
        clearTimeout(timeoutId)
        showSlides()
    } else {
        showingSlides = false
    }

  })