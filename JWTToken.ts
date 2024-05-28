import {randomBytes,createHmac} from 'crypto'

function base64URLEncode(string:string){
    let encodedString=btoa(string);
    encodedString=encodedString.replaceAll('+','-').replaceAll("=","").replaceAll("/","_");
    return encodedString;
}

const key=randomBytes(32).toString('base64');

export function createJWTToken(payload:Record<any,any>){
    const header={
        "alg": "HS256",
        "typ": "JWT"
    }
    
    const headerPlusPayLoad=base64URLEncode(JSON.stringify(header))+'.'+base64URLEncode(JSON.stringify(payload));

    const signature=createHmac('sha256',key).update(headerPlusPayLoad).digest('base64url');
    
    return headerPlusPayLoad+'.'+signature;
}

export function isJWTTokenValid(token:string){
    
    const [header,payload,signature]=token.split('.') as [string,string,string];

    //Recreate the signature using header+'.'+payload
    const signature2=createHmac('sha256',key).update(header+'.'+payload).digest('base64url');

    return signature===signature2;
}