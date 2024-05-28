function hasAtLeastOneCapitalLetter(string:string){
    return /[A-Z]/.test(string);
}

function hasAtLeastOneLowecaseLetter(string:string){
    return /[a-z]/.test(string);
}

function hasAtLeastOneDigit(string:string){
    return /[0-9]/.test(string);
}

function hasAtLeastOneSymbol(string:string){
    return /[ `!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]/.test(string);
}

function hasAtLeast11Characters(string:string){
    return string.length>11;
}

export function isEmailValid(email:string){
    return email.match(/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/);
}

export function isPasswordOk(password:string){
    return hasAtLeastOneCapitalLetter(password) && hasAtLeastOneLowecaseLetter(password) && hasAtLeastOneDigit(password) &&
    hasAtLeastOneSymbol(password) && hasAtLeast11Characters(password);
}