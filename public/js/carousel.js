let slideIndex = 0;
let timeoutId = 0;
let showingSlides = false;

var startX;
var startY;

if (location.pathname == "/") {
    slideIndex = 0
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
        slideIndex = 0
        showSlides()
    } else {
        showingSlides = false
    }

  })

function handleTouchStart(event) {
    startX = event.changedTouches[0].screenX;
    startY = event.changedTouches[0].screenY;
}

function handleTouchEnd(event) {
    endX = event.changedTouches[0].screenX;
    endY = event.changedTouches[0].screenY;
    if(Math.abs(startX - endX) > 30 && Math.abs(startY - endY) < 100) {
        if(endX < startX) {
            clickShowSlides(true)
        } else {
            clickShowSlides(false)
        }
    }
}

function handleCalendarTouchEnd(event) {
    endX = event.changedTouches[0].screenX;
    endY = event.changedTouches[0].screenY;
    if(Math.abs(startX - endX) > 50 && Math.abs(startY - endY) < 100) {
        if(endX < startX) {
            document.querySelector("#NextMonth").click();
        } else {
            document.querySelector("#PrevMonth").click();
        }
    }
}

function addNewAnnotationForm() {
    newAnnotationRow = `<tr>
                <td>
                    <input type="text" placeholder="Class Annotation" class="input input-bordered w-full"
                        name="Annotations" />
                </td>
                <td>
                    <div class="btn btn-warning btn-sm md:btn-md" onclick="this.parentElement.parentElement.remove()">
                        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" width="24" height="24">
                                <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
                                <g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
                                <g id="SVGRepo_iconCarrier">
                                    <path d="M6 12L18 12" stroke="#000000" stroke-width="2" stroke-linecap="round"
                                        stroke-linejoin="round"></path>
                                </g>
                            </svg>
                    </div>
                </td>
            </tr>`
    document.querySelector('tbody').insertAdjacentHTML("beforeend", newAnnotationRow)
}