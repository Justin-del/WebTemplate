//This file contains functions related to authentication and authorization.

import { createAHTTPUnsupportedMediaTypeResponse, createAHTTPUnprocessableContentResponse, createAHTTPBadRequestResponse, createAHTTPForbiddenResponse } from "./HTTPResponseHelpers";
import { createJWTToken, isJWTTokenValid } from "./JWTToken";
import { nameOfDatabaseFile } from "./Settings";
import { isEmailValid, isPasswordOk } from "./ValidationFunctionsForEmailAndPassword";
import { Database } from "bun:sqlite";
import nunjucks from "nunjucks";
/**
 * 
 * @param fileName Name of the template file without any extensions. Template files are defined in the views folder.
 */
function renderTemplate(fileName:string,context:Record<any,any>={}){
    const renderedString = nunjucks.render(fileName+'.njk', context);
    return new Response(renderedString, {headers:{
      "Content-Type":'text/html'
    }} )
}

export async function handlesPostRequestForTheLoginRoute(request:Request){
    if (!(request.headers.get('content-type')?.includes("multipart/form-data") ||request.headers.get('content-type')?.includes("application/x-www-form-urlencoded"))){
      return createAHTTPUnsupportedMediaTypeResponse();
    }
    const data = await request.formData();
  
    let email=data.get('email');
    let password = data.get('password');
  
    if (email==null){
      return createAHTTPUnprocessableContentResponse("Missing email field in body.");
    }
  
    if (password==null){
      return createAHTTPUnprocessableContentResponse("Missing password field in body.");
    }
  
    //Do database related stuff.
    {
      using db = new Database(nameOfDatabaseFile,{create:true});
  
      //Get the user from database.
      const query=db.query("select email,hashedPassword from users where email=$email");
      const user=query.get({$email:email.toString()}) as {email:string,hashedPassword:string};
      query.finalize();
  
      if (!user){
        return renderTemplate("Login",{isAuthorized:isUserAuthorized(request),email:email.toString(),password:password.toString(), emailExistsInDatabase:false});
      }
  
      if (!(await Bun.password.verify(password.toString(),user.hashedPassword))){
        return renderTemplate("Login",{isAuthorized:isUserAuthorized(request),email:email.toString(), 
          password:password.toString(), isPasswordIncorrect:true});
      }
    }
  
    return new Response("<body onload=window.location.href='/authorized'></body>", {
      headers: {
        'Set-Cookie': `token=${createJWTToken({email})}; HttpOnly;Secure`,
        'Content-type':'text/html'
      }
    });
  
  }

  export async function handlesPostRequestForTheSignUpRoute(request:Request){
    if (!(request.headers.get('content-type')?.includes("multipart/form-data") ||request.headers.get('content-type')?.includes("application/x-www-form-urlencoded"))){
      return createAHTTPUnsupportedMediaTypeResponse();
    }
  
    const data = await request.formData();
  
    let email=data.get('email');
    let password = data.get('password');
    let confirmPassword=data.get('confirmPassword');
  
    if (email==null)
      return createAHTTPUnprocessableContentResponse("Missing email field in body.");
  
    if (password==null)
      return createAHTTPUnprocessableContentResponse("Missing password field in body.");
  
    if (confirmPassword==null)
      return createAHTTPUnprocessableContentResponse("Missing confirmPassword field in body.");
  
    if (!isEmailValid(email.toString())){
      return createAHTTPBadRequestResponse("Invalid email.")
    }
  
    if (!isPasswordOk(password.toString())){
      return createAHTTPForbiddenResponse("Password does not meet requirements. The requirements for a password are at least 11 characters, at least one uppercase letter, at least one lowercase letter, at least one digit and at least one symbol.")
    }
  
    if (confirmPassword!==password){
      return createAHTTPForbiddenResponse("Confirm password does not match password.");
    }
  
  
    //Do database related stuff.
    {
      using db = new Database(nameOfDatabaseFile,{create:true});
  
      //See if there's already a user with the given email in the database.
      let query = db.query("select email,hashedPassword from users where email=$email");
      const user=query.get({$email:email.toString()}) as {email:string,hashedPassword:string};
      query.finalize();
      if (user){
        return renderTemplate("SignUp",{isAuthorized:isUserAuthorized(request),email:email.toString(),password:password.toString(), confirmPassword:confirmPassword.toString(), emailAlreadyExistsInDatabase:true});
      }
  
      //Add user into database.
      query=db.query(`INSERT INTO users (email,hashedPassword) values (?,?)`);
      query.run(email.toString(), await Bun.password.hash(password.toString()));
    }
  
    return new Response("<body onload=window.location.href='/authorized'></body>", {
      headers: {
        'Set-Cookie': `token=${createJWTToken({email})}; HttpOnly;Secure`,
        'Content-type':'text/html'
      }
    });
  }



export function isUserAuthorized(request:Request){
    const cookie=request.headers.get('cookie')?.split(';')?.reduce((res:Record<any,any>,item)=>{
      const [name, value] = item.trim().split('=');
      res[name] = value;
      return res;
    },{})
  
    if (!cookie){
      return false;
    }

    if (cookie['token']===undefined) return false;
  
    return isJWTTokenValid(cookie['token']);
}

export function createALogoutResponse(){
        //logs the user out by just setting their token to an empty string.
        return new Response("<body onload=window.location.href='/login'></body>", {
        headers: {
            'Set-Cookie': `token=''; HttpOnly;Secure`,
            'Content-type':'text/html'
        }
        });
}