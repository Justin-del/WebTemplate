import {displaySuccesfulMessage, displayUnsuccesfulMessage, clearMessages} from './Utils.js'

/**
 * @param {string} username
 */
export async function signUp(username){
    clearMessages()
    let response = await fetch("/SignUp/RegistrationData");
    
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
            authenticatorAttachment:"cross-platform"
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
            publicKey:publicKeyCredentialCreationOptions
        })
    }catch(error){
        if (error.message.includes("The operation either timed out or was not allowed.")){
            return;
        }
        displayUnsuccesfulMessage(error.message)
        return;
    }

    response=await fetch(`/SignUp/${registrationData.Challenge.Id}/${user_id}`,{
        method:'POST',
        body:JSON.stringify(credential),
        headers:{
            'Content-Type':'application/json'
        }
    })

    if (response.status!==200){
        displayUnsuccesfulMessage()
    } else {
        displaySuccesfulMessage('Succesfully signed up! You can now proceed to <a hx-get="/login" hx-replace-url="true" hx-target="body" href="javascript:;">login</a>')
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
            alert("Failed to connect to the server. Please ensure that you are connected to the Internet.")
        }
    }
};