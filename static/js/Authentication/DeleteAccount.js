import {displayUnsuccesfulMessage, clearMessages, displaySuccesfulMessage} from './Utils.js'

export async function deleteAccount(){
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
    
    response = await fetch(`/deleteAccount/${authenticationData.Challenge.Id}`, {
        method:"POST",
        body:JSON.stringify(credential),
        headers:{
            "Content-Type":"application/json"
        }
    })

    if (response.status === 200){
        displaySuccesfulMessage("Account deletion was succesful. Don't forget to delete the passkeys associated with this account as well!")
    } else {
        displayUnsuccesfulMessage("Account deletion failed. Please try again.")
    }

}

window.deleteAccount = deleteAccount;