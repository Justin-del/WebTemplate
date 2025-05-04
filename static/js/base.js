window.is_navbar_expanded = false;
window.maximum_width_for_navbar_to_be_displayed_in_mobile_mode = 600;
window.theme = "light";

function saveThemeIntoLocalStorage(){
    window.localStorage.setItem("theme", window.theme)
}

/**
 * 
 * @param {HTMLAnchorElement} element 
 */
function toggleTheme(element){
    if (window.theme === "light"){
        window.theme = "dark"
    } else {
        window.theme = "light"
    }

    saveThemeIntoLocalStorage()

    if (window.theme == "light"){
        document.body.classList.add("light-theme")
        document.body.classList.remove("dark-theme")

        element.textContent="☀"
    } else {
        document.body.classList.add("dark-theme")
        document.body.classList.remove("light-theme")

        element.textContent="🌙"
    }
}

function expandNavbar(){

    is_navbar_expanded=true;
    const navbarContents = document.querySelector("nav ul");
    const navbarToggler = document.querySelector("nav button");

    navbarContents.classList.remove("hide-contents")

    navbarContents.style.height = navbarContents.scrollHeight+'px';
    navbarContents.inert = false;

    navbarToggler.innerHTML = 'X';
    navbarToggler.setAttribute('aria-expanded','true');

}

function collapseNavbar(){

    is_navbar_expanded=false;
    const navbarContents = document.querySelector("nav ul");
    const navbarToggler = document.querySelector("nav button");

    navbarContents.classList.add("hide-contents")
    
    navbarContents.style.height = '0';
    navbarContents.inert = true;

    navbarToggler.innerHTML = '&#8801'; //unicode icon for hamburger;
    navbarToggler.setAttribute('aria-expanded','false');

}

//Expand navbar if it's collapsed and collapse the navbar if it's expanded.
function toggleNavbar(){
    is_navbar_expanded = !is_navbar_expanded;
    if (is_navbar_expanded){
        expandNavbar();
    } else {
        collapseNavbar();
    }
}

function updateNavbarStyleBasedOnScreenWidth(){
    const navbar = document.querySelector("nav");
    const navbarContents = document.querySelector("nav ul");

    if (window.screen.availWidth > maximum_width_for_navbar_to_be_displayed_in_mobile_mode){
        navbarContents.classList.remove('hide-contents')
        navbar.classList.add('desktop-style');
        navbarContents.inert = false;
    } else {
        navbar.classList.remove('desktop-style');
        if (is_navbar_expanded){
            expandNavbar();
        } else {
            collapseNavbar();
        }
    }
}

window.onresize=updateNavbarStyleBasedOnScreenWidth;

window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', event => {
    window.theme = event.matches ? "dark" : "light";

    saveThemeIntoLocalStorage()

    if (window.theme == "light"){
        document.body.classList.add("light-theme")
        document.body.classList.remove("dark-theme")
    } else {
        document.body.classList.add("dark-theme")
        document.body.classList.remove("light-theme")
    }

});

function setTheme(){
    //Get theme from local storage.
    const savedTheme = window.localStorage.getItem('theme')

    if (savedTheme !== null){
        window.theme=savedTheme
    } else {
        if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
            window.theme="dark"
        } else {
            window.theme = "light"
        }
    }

    const anchorTag = document.querySelector("a.theme")

    if (window.theme === "light"){
        document.body.classList.add("light-theme")
        document.body.classList.remove("dark-theme")
        anchorTag.textContent="☀"
    } else {
        document.body.classList.add("dark-theme")
        document.body.classList.remove("light-theme")
        anchorTag.textContent="🌙"
    }
}

async function registerServiceWorker(){
    if ("serviceWorker" in navigator){
        await navigator.serviceWorker.register("./static/js/ServiceWorker.js", {
            scope:'/'
        });
    }
}

/**
 * 
 * @param {string} message 
 */
function setScreenReaderOnlyMessage(message){
    document.querySelector(".screen-readers-only").textContent = message;
}



