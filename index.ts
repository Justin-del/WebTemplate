import { Database } from "bun:sqlite";
import nunjucks from "nunjucks";
import { handlesPostRequestForTheLoginRoute, handlesPostRequestForTheSignUpRoute, isUserAuthorized, createALogoutResponse } from "./Auth";
import { createAHTTPMethodNotAllowedResponse, createAHTTPNotFoundResponse } from "./HTTPResponseHelpers";
import { nameOfDatabaseFile } from "./Settings";

nunjucks.configure("views",{autoescape:true})
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


function createUsersTableIfITDoesNOTExists(){
  {
    using db=new Database(nameOfDatabaseFile,{create:true});
    //email and hashedPassword  are must have fields in order for the authentication and authorization
    // functions to work.
    db.run("create table if not exists users(email TEXT PRIMARY KEY, hashedPassword TEXT)"); 
  }
}


createUsersTableIfITDoesNOTExists();


Bun.serve({
    async fetch(request){
      const url = new URL(request.url);
      let isAuthorized=isUserAuthorized(request);

      if (url.pathname=="/login"){
        if (request.method=="GET") return renderTemplate("Login",{isAuthorized});
        if (request.method=="POST") return await handlesPostRequestForTheLoginRoute(request);
        return createAHTTPMethodNotAllowedResponse();
      } else if (url.pathname=="/signUp"){
        if (request.method=="GET")  return renderTemplate("SignUp",{isAuthorized});
        if (request.method=="POST") return await handlesPostRequestForTheSignUpRoute(request);
        return createAHTTPMethodNotAllowedResponse();
      } else if (url.pathname=="/authorized"){
        if (!isAuthorized)
          return renderTemplate("NotAuthorized",{isAuthorized});
        if (request.method=="GET")  return renderTemplate("Authorized",{isAuthorized})
        return createAHTTPMethodNotAllowedResponse();
      } else if (url.pathname=="/logout"){
        return createALogoutResponse();
      }

      if (request.method=="GET"){
        return renderTemplate("PageNotFound");
      }
      return createAHTTPNotFoundResponse();
  }

});
