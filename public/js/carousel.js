let slideIndex = 0;
let timeoutId = 0;
let showingSlides = false;

var touchStartX;
var touchStartY;


slideIndex = 0
showSlides();


function clickShowSlides(plus) {
    clearTimeout(timeoutId)
    if (plus) {
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
        if (dots.length != 0) {
            dots[slideIndex - 1].className += " dot-active";
        }
        timeoutId = setTimeout(showSlides, 5000); // Change image every 2 seconds
    } else {
        clearTimeout(timeoutId)
    }
}

htmx.on('htmx:afterRequest', (evt) => {
    // check which element triggered the htmx request. If it's the one you want call the function you need
    //you have to add htmx: before the event ex: 'htmx:afterRequest'
    console.log(evt)

    if (evt.detail.pathInfo.requestPath == "/") {
        console.log(evt.detail.pathInfo.requestPath)
        clearTimeout(timeoutId)
        slideIndex = 0
        showSlides()
    } else {
        showingSlides = false
    }

})

function handleHoverStart(event) {
    if (!event.currentTarget.getAttribute('open')) {
        event.currentTarget.setAttribute('open', true)
    }
}

function handleMenuHoverOut(event) {
    hoverEndX = event.clientX
    hoverEndY = event.clientY
    var elementRect = event.currentTarget.getBoundingClientRect()

    if (hoverEndX >= elementRect.right - 2 || hoverEndX <= elementRect.left || hoverEndY <= elementRect.top || hoverEndY >= elementRect.bottom + 10) {
        event.currentTarget.removeAttribute('open')
    }
}

function handleTouchStart(event) {
    touchStartX = event.changedTouches[0].screenX;
    touchStartY = event.changedTouches[0].screenY;
}

function handleTouchEnd(event) {
    endX = event.changedTouches[0].screenX;
    endY = event.changedTouches[0].screenY;
    if (Math.abs(touchStartX - endX) > 30 && Math.abs(touchStartY - endY) < 100) {
        if (endX < touchStartX) {
            clickShowSlides(true)
        } else {
            clickShowSlides(false)
        }
    }
}

function handleCalendarTouchEnd(event) {
    endX = event.changedTouches[0].screenX;
    endY = event.changedTouches[0].screenY;
    if (Math.abs(touchStartX - endX) > 50 && Math.abs(touchStartY - endY) < 100) {
        if (endX < touchStartX) {
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
                    <div class="btn btn-primary btn-sm md:btn-md" onclick="this.parentElement.parentElement.remove()">
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

function selectExistingBannerImage(event) {
    if (event.target.classList.contains('border-yellow-500')){
        event.target.classList.remove('border-yellow-500');
        event.target.classList.add('border-sky-500');
        document.querySelector('input[name=\'ExistingCoverPhoto\']').value = '';
    } else {
        if (document.querySelector('.cover-photo.border-yellow-500') != null) {
            document.querySelector('.cover-photo.border-yellow-500').classList.add('border-sky-500');
            document.querySelector('.cover-photo.border-yellow-500').classList.remove('border-yellow-500');
        }
        event.target.classList.add('border-yellow-500');
        event.target.classList.remove('border-sky-500');
        document.querySelector('input[name=\'ExistingCoverPhoto\']').value = event.target.getAttribute('src');
    }
    previewExistingBannerImage(event);
}

function previewNewBannerImage(event) {
    var output = document.querySelector('.banner-image');
    output.src = URL.createObjectURL(event.target.files[0]);
    output.onload = function () {
        URL.revokeObjectURL(output.src) // free memory
    }
};

function previewExistingBannerImage(event) {
    var output = document.querySelector('.banner-image');
    console.log(event)
    output.src = event.target.src;
};

const image = document.getElementById('scaling-image');
const height = image.offsetHeight;

function handleScroll() {
    if (window.innerWidth >= 1024) { // Tailwind's md breakpoint is 768px
        const scrollPosition = window.scrollY;
        let newScale = 100 - (scrollPosition * 0.1);
        if (newScale < 50) newScale = 50;
        if (newScale > 100) newScale = 100;
        
        const newHeight = (newScale / 100) * height;
        image.style.height = newHeight + 'px';
    } else {
        // Reset to original height on smaller screens
        image.style.height = '';
    }
}

// Add event listeners
window.addEventListener('scroll', handleScroll);
// Handle window resize to enable/disable effect
window.addEventListener('resize', function() {
    // Only update height if we're crossing the md breakpoint
    const isNowDesktop = window.innerWidth >= 1024;
    const wasDesktop = image.offsetHeight !== 0 && image.offsetHeight !== height;
    
    if(!isNowDesktop){
        image.style.height = '';
    }


    // if (isNowDesktop !== wasDesktop) {
    //     height = image.offsetHeight; // Update the base height
    //     handleScroll();
    // }
});