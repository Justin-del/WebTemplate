window.is_navbar_expanded = false;
window.maximum_width_for_navbar_to_be_displayed_in_mobile_mode = 400;

function expandNavbar(){
    is_navbar_expanded=true;
    const navbarContents = document.querySelector("nav ul");
    const navbarToggler = document.querySelector("nav button");

    navbarContents.style.height = navbarContents.scrollHeight+'px';
    navbarContents.inert = false;

    navbarToggler.innerHTML = 'X';
    navbarToggler.setAttribute('aria-expanded','true');
}

function collapseNavbar(){
    is_navbar_expanded=false;
    const navbarContents = document.querySelector("nav ul");
    const navbarToggler = document.querySelector("nav button");
    
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