import {displaySuccesfulMessage, displayUnsuccesfulMessage, clearMessages} from './Utils.js'


/**
 * @param {string} username
 */
export async function signUp(username){
    clearMessages()

    //Check if username exists in the database.
    let response = await fetch("/signUp/isUsernameTaken",{
        body:JSON.stringify({username}),
        method:'POST'
    })

    /**
     * @type {{isUsernameTaken: boolean}}
     */
    let {isUsernameTaken} = await response.json()

    if (isUsernameTaken){
        displayUnsuccesfulMessage("This username is already taken. Please choose another one.")
        return;
    }

    response = await fetch("/signUp/RegistrationData");
    
    /**@type {{Challenge:{
      Id:number,
     Challenge:string
     },
     RP:{
     Id:string,
     Name:string
     },
     SupportedCoseAlgorithms: number[],
     TimeoutInMinutes: number }}  */
    let registrationData = await response.json()
    
    let pubKeyCredParams=registrationData.SupportedCoseAlgorithms.map(algorithm=>{
        return {
            alg:algorithm,
            type:"public-key"
        }
    })

    const user_id = self.crypto.randomUUID()

    const publicKeyCredentialCreationOptions = {
        challenge:Uint8Array.from(atob(registrationData.Challenge.Challenge), c=>c.charCodeAt(0)),
        rp:{
            name:registrationData.RP.Name,
            id:registrationData.RP.Id
        },
        user:{
            id:new TextEncoder().encode(user_id),
            name:username,
            displayName:username
        },
        authenticatorSelection:{
            userVerification:"required",
            residentKey:"required"
        },
        pubKeyCredParams:pubKeyCredParams,
        timeout:registrationData.TimeoutInMinutes*60*1000, //timeout will be in milliseconds. That's why timeoutInMinutes is multiplied by 60 to get the number of seconds and further multiplied by 1000 to get the number of milliseconds.
        attestation:"none",
    }

    /**
     * @type PublicKeyCredential
     */
    let credential;
    try{
        credential = await navigator.credentials.create({
            publicKey:publicKeyCredentialCreationOptions,
        })
    }catch(error){
        if (error.message.includes("The operation either timed out or was not allowed.")){
            return;
        }

        if (error.message.includes("denied permission")){
            return;
        }
        displayUnsuccesfulMessage(error.message)
        return;
    }

    response=await fetch(`/signUp/${registrationData.Challenge.Id}/${user_id}`,{
        method:'POST',
        body:JSON.stringify({credential,username}),
        headers:{
            'Content-Type':'application/json'
        }
    })

    if (response.status!==200){
        displayUnsuccesfulMessage('Failed to sign up due to an unknown error. If a passkey was created, you will need to delete it.')
    } else {
        displaySuccesfulMessage('Succesfully signed up! You can now proceed to <a hx-on::send-error="Cannot connect to the server.  Please check your internet connection. If your connection is stable, there might be a temporary issue with the website. Please try again in a few minutes." hx-get="/login" hx-replace-url="true" hx-target="body" href="javascript:;">login</a>')
    }
}

/**
 * 
 * @param {string} username 
 */
window.signUp=async(username)=>{
    try{
        await signUp(username)
    } catch (error){
        if (error.message === "Failed to fetch"){
            alert("Cannot connect to the server.  Please check your internet connection. If your connection is stable, there might be a temporary issue with the website. Please try again in a few minutes.")
        }
    }
};