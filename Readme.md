This is a template for creating websites in go. It includes authentication and session management.
It uses sqlite database as its default database.
It also uses HTMX for a smoother page experience.
<br/><br/>
You may need to change the RP variable in [WebAuthn.go](/Utils/WebAuthn/WebAuthn.go).
<br><br/>
Also, see [Authorized.go](/RoutesHandler/Authorized/Authorized.go) for an example on how to handle access to authorized pages.

Also, you might want to change /authorized to the route that you would like the user to access after logging in at [Authentication.js](/static/js/Authentication.js)

Also, often times, you might want the navigation bar to display different links conditionally. You can have a look at [base.html](/templates/base.html) for how this is done.

Also, you might want to change the session timeout at [Sessions.go](/Database/Sessions/Sessions.go).