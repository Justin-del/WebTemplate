import {displayUnsuccesfulMessage, clearMessages} from './Utils.js'

export async function login(){
    clearMessages();
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

window.login = login;
