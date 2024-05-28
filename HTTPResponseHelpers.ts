export function createAHTTPUnprocessableContentResponse(reason:string){
    return new Response(JSON.stringify({reason}),{
      headers:{
        "Content-Type":'application/json',
      },
      status:422,
      statusText:"Unprocessable content"
    })
  }
  
export function createAHTTPUnauthorizedResponse(reason:string){
    return new Response(JSON.stringify({reason}),{
      headers:{
        "Content-Type":'application/json',
      },
      status:401,
      statusText:"Unauthorized"
    })
}

export function createAHTTPConflictResponse(reason:string){
    return new Response(JSON.stringify({reason}),{
        headers:{
          "Content-Type":'application/json',
        },
        status:409,
        statusText:"Conflict"
      })
}

export function createAHTTPBadRequestResponse(reason:string){
    return new Response(JSON.stringify({reason}),{
        headers:{
          "Content-Type":'application/json',
        },
        status:400,
        statusText:"Bad request."
      })
}

export function createAHTTPForbiddenResponse(reason:string){
    return new Response(JSON.stringify({reason}),{
        headers:{
          "Content-Type":'application/json',
        },
        status:403,
        statusText:"Forbidden"
      })
}

export function createAHTTPMethodNotAllowedResponse(){
    return new Response("Method not allowed.",{
        headers:{
          "Content-Type":'application/json',
        },
        status:405,
        statusText:"Method not allowed."
      })
}

export function createAHTTPNotFoundResponse(){
    return new Response("Not found",{
        headers:{
          "Content-Type":'application/json',
        },
        status:404,
        statusText:"Not found"
      })
}

export function createAHTTPUnsupportedMediaTypeResponse(){
    return new Response("Unsupported media type.",{
        headers:{
          "Content-Type":'application/json',
        },
        status:415,
        statusText:"Unsupported media type."
      })
}


