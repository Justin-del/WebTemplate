
const cacheFirstWithCacheRefresh  = {
    "urls":["https://unpkg.com/htmx.org@2.0.4"],
    "paths":["/","/SignUp","/login","/deleteAccount"],
    "directories":["/static"]
}

/**
 * This origin should match the origin of your application.
 */
const origin = "http://localhost:8080"

/**
 * 
 * @param {string} url_string
 */
function isCacheFirstWithCacheRefresh(url_string){
    if (cacheFirstWithCacheRefresh["urls"].includes(url_string)){
        return true;
    }
    const url = new URL(url_string)

    if (url.origin == origin){
        if (cacheFirstWithCacheRefresh["paths"].includes(url.pathname)){
            return true;
        }

        for (const directory of cacheFirstWithCacheRefresh["directories"]) {
            if (url.pathname.startsWith(directory)) {
                return true;
            }
        }    
    }

    return false;
}

self.addEventListener("fetch", (event) => {
    event.respondWith(
        (async()=>{
            if (event.request.method !== "GET"){
                return await fetch(event.request)
            }

            const responseFromNetwork = (async()=>{
                try{
                    const response = await fetch(event.request,{mode:'no-cors'});
                    const cache = await caches.open("cache");
                    await cache.put(event.request, response.clone());
                    return response;
                } catch(error){
                    return undefined;
                }
            })();

            if (isCacheFirstWithCacheRefresh(event.request.url)){
                return await caches.match(event.request) || await responseFromNetwork || new Response("You are offline.");
            }

            const response = await fetch(event.request)
            return response
            
        
        })(),
    );
});