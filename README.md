# WebTemplate

To install dependencies:

```bash
bun install
```

To run:

```bash
bun run index.ts
```

This project was created using `bun init` in bun v1.1.9. [Bun](https://bun.sh) is a fast all-in-one JavaScript runtime.

## List of Directories And/Or Files of The Framework and Their Purposes
|Name of directory or file|Purpose                                                                        |
|-------------------------|-------------------------------------------------------------------------------|
|Views                    | Used for storing njk (template) files.                                        |
|-------------------------|-------------------------------------------------------------------------------|
|.gitignore               | Used for storing files that should be ignored by Github.                      |
|auth.ts                  | Contains a bunch of useful authentication and authorization related functions.|
|                         | Here are a list of the function signatures that the file exposes.             |
|                         | export async function handlesPostRequestForTheLoginRoute(request:Request)     |
|                         | export async function handlesPostRequestForTheSignUpRoute(request:Request)    |
|                         | export function isUserAuthorized(request:Request)                             |
|                         | export function createALogoutResponse()                                       |

