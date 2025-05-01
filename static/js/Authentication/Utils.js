/**
 * 
 * @param {string|undefined} message 
 */
export function displayUnsuccesfulMessage(message){
    const errorMessageElement = document.querySelector("form .error-message")
    const succesfulMessageElement = document.querySelector("form .succesful-message")

    if (message===undefined){
        message="The operation either timed out or was not allowed."
    } else if (message.includes("The operation is not allowed at this time because the page does not have focus.")){
        message+=" If a new passkey was created, you will need to delete it."
    }

    errorMessageElement.textContent=message;

    if (succesfulMessageElement){
        succesfulMessageElement.textContent  = "";
    }
}

/**
 * 
 * @param {string} message 
 */
export function displaySuccesfulMessage(message){
    const errorMessageElement = document.querySelector("form .error-message")
    const succesfulMessageElement = document.querySelector("form .succesful-message")

    succesfulMessageElement.innerHTML=message;
    if (errorMessageElement){
        errorMessageElement.textContent = "";
    }
}

export function clearMessages(){
    const errorMessageElement = document.querySelector("form .error-message")
    const succesfulMessageElement = document.querySelector("form .succesful-message")
    
    if (errorMessageElement){
        errorMessageElement.textContent = "";
    }

    if (succesfulMessageElement){
        succesfulMessageElement.textContent = "";
    }
}