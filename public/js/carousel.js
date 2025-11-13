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
    if (event.target.classList.contains('border-yellow-500')) {
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

function previewNewCardImage(event) {
    var output = document.querySelector('.card-image');
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

function updateInstructorPreview() {
    // Update name
    const nameInput = document.getElementById('name-input');
    const previewName = document.getElementById('preview-name');
    if (nameInput && previewName) {
        previewName.textContent = nameInput.value || 'Instructor Name';
    }

    // Update bio
    const bioInput = document.getElementById('bio-input');
    const previewBio = document.getElementById('preview-bio');
    if (bioInput && previewBio) {
        previewBio.textContent = bioInput.value || 'Instructor bio will appear here...';
    }
}

function previewInstructorImage(input) {
    const preview = document.getElementById('preview-picture');
    const currentPicture = document.getElementById('current-picture');

    if (input.files && input.files[0]) {
        const reader = new FileReader();

        reader.onload = function (e) {
            preview.src = e.target.result;
            preview.classList.remove('hidden');
            if (currentPicture) {
                currentPicture.src = e.target.result;
            }
        }

        reader.readAsDataURL(input.files[0]);
    } else if (currentPicture && currentPicture.src) {
        // If no new file is selected but there's a current picture, use that
        preview.src = currentPicture.src;
        preview.classList.remove('hidden');
    } else {
        // If no picture is available, hide the preview
        preview.classList.add('hidden');
    }
}

function resetImageTransform() {
    // Reset form inputs
    const zoomInput = document.querySelector('input[name="ZoomLevel"]');
    const offsetXInput = document.querySelector('input[name="OffsetX"]');
    const offsetYInput = document.querySelector('input[name="OffsetY"]');
    
    // Set default values
    zoomInput.value = 100;
    offsetXInput.value = 0;
    offsetYInput.value = 0;
    
    // Update the displayed values
    document.getElementById('zoomValue').textContent = '100%';
    document.getElementById('offsetXValue').textContent = '0px';
    document.getElementById('offsetYValue').textContent = '0px';
    
    // Update the preview
    const preview = document.getElementById('preview-picture');
    updateTransform(preview, {
        scale: 1.0,
        translateX: '0px',
        translateY: '0px'
    });
}

function updateTransform(preview, updates) {
    // Get current transform string and parse it into an object
    const transformStr = preview.style.transform || '';
    const transformMap = {};
    
    // Parse existing transforms into an object
    transformStr.split(' ').forEach(transform => {
        if (!transform) return;
        const match = transform.match(/(\w+)\(([^)]+)\)/);
        if (match) {
            transformMap[match[1]] = match[2];
        }
    });
    
    // Update with new values
    if (updates.scale !== undefined) transformMap.scale = updates.scale;
    if (updates.translateX !== undefined) transformMap.translateX = updates.translateX;
    if (updates.translateY !== undefined) transformMap.translateY = updates.translateY;
    
    // Build new transform string
    const newTransform = [
        transformMap.translateX !== undefined ? `translateX(${transformMap.translateX})` : '',
        transformMap.translateY !== undefined ? `translateY(${transformMap.translateY})` : '',
        transformMap.scale !== undefined ? `scale(${transformMap.scale})` : ''
    ].filter(Boolean).join(' ');
    
    preview.style.transform = newTransform;
}

// Update zoom value display and transform
function updateZoomValue(input) {
    document.getElementById('zoomValue').textContent = `${input.value}%`;
    const preview = document.getElementById('preview-picture');
    updateTransform(preview, {
        scale: input.value / 100.0
    });
}

// Update X offset value display and transform
function updateOffsetXValue(input) {
    document.getElementById('offsetXValue').textContent = `${input.value}px`;
    const preview = document.getElementById('preview-picture');
    updateTransform(preview, {
        translateX: `${input.value}px`
    });
}

// Update Y offset value display and transform
function updateOffsetYValue(input) {
    document.getElementById('offsetYValue').textContent = `${input.value}px`;
    const preview = document.getElementById('preview-picture');
    updateTransform(preview, {
        translateY: `${input.value}px`
    });
}


// Add event listeners
window.addEventListener('scroll', handleScroll);
// Handle window resize to enable/disable effect
window.addEventListener('resize', function () {
    // Only update height if we're crossing the md breakpoint
    const isNowDesktop = window.innerWidth >= 1024;
    const wasDesktop = image.offsetHeight !== 0 && image.offsetHeight !== height;

    if (!isNowDesktop) {
        image.style.height = '';
    }


    // if (isNowDesktop !== wasDesktop) {
    //     height = image.offsetHeight; // Update the base height
    //     handleScroll();
    // }
});