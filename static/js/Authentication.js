/**
 * 
 * @param {string|undefined} message 
 */
function displayUnsuccesfulMessage(message){
    if (message===undefined){
        message="The operation either timed out or was not allowed."
    } else if (message.includes("The operation is not allowed at this time because the page does not have focus.")){
        message+=" If a new passkey was created, you will need to delete it."
    }
    document.querySelector("form .error-message").textContent=message;
    document.querySelector("form .succesful-message").textContent="";
}

function displaySuccesfulMessage(){
    document.querySelector("form .succesful-message").innerHTML="Succesfully signed up! You can now proceed to <a href='/login'>login</a>.";
    document.querySelector("form .error-message").textContent="";
}

/**
 * @param {string} username
 */
async function signUp(username){

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
        displaySuccesfulMessage()
    }
}

async function login(){
    let response = await fetch("/login/AuthenticationData")

    /**@type{{
      Challenge:{
      Id:number,
      Challenge:string
      },
      RelyingPartyId:string,
      TimeoutInMinutes:number
    }} */
    let authenticationData = await response.json()

    const credential = await navigator.credentials.get({
        publicKey:{
            challenge:Uint8Array.from(atob(authenticationData.Challenge.Challenge), c=>c.charCodeAt(0)),
            rpId: authenticationData.RelyingPartyId,
            timeout:authenticationData.TimeoutInMinutes*1000*60, //timeout need to be in milliseconds
            userVerification:"required"
    }
    })
    
    response = await fetch(`/login/${authenticationData.Challenge.Id}`, {
        method:"POST",
        body:JSON.stringify(credential),
        headers:{
            "Content-Type":"application/json"
        }
    })

    if (response.status === 200){
        //Change this to where you would like the user to be redirected after the user logs in.
        window.location.href="/authorized"
    } else {
        displayUnsuccesfulMessage("Login failed. Please try again.")
    }

}