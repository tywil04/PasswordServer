import * as crypto from "/public/js/lib/crypto.js"
import * as utils from "/public/js/lib/utils.js"

window.auth = {}

window.auth.signup = async () => {
  const emailInput = document.querySelector("#email")
  const passwordInput = document.querySelector("#password")

  let masterKey = await crypto.generateMasterKey(passwordInput.value, emailInput.value) // Derive a key via pbkdf2 from the users password and email using
  let masterHash = utils.arrayBufferToHex(await crypto.generateMasterHash(passwordInput.value, masterKey)) // Derive bits via pbkdf2 from the masterkey and the users password (this is used for server-side auth)

  let databaseKey = await crypto.generateDatabaseKey() // generate random AES-256-CBC key
  let [iv, encryptedDatabaseKey] = await crypto.protectDatabaseKey(masterKey, databaseKey) // encrypt the key with masterkey
  let protectedDatabaseKey = utils.arrayBufferToHex(iv) + ";" + utils.arrayBufferToHex(encryptedDatabaseKey)

  let response = await fetch(window.backendRoutes.signup, {
    method: "POST",
    body: JSON.stringify({
      Email: emailInput.value,
      MasterHash: masterHash,
      ProtectedDatabaseKey: protectedDatabaseKey,
    })
  })
  let jsonResponse = await response.json()

  let success = jsonResponse.UserId !== undefined // an error response would not contain "UserId" instead it would contain "Error"

  if (success) {
    window.location = window.frontendRoutes.signin
    console.log(jsonResponse.UserId)
  } else {
    window.location.reload()
  }
}

window.auth.signin = async () => {
  const emailInput = document.querySelector("#email")
  const passwordInput = document.querySelector("#password")

  let masterKey = await crypto.generateMasterKey(passwordInput.value, emailInput.value) // Derive a key via pbkdf2 from the users password and email using
  let masterHash = utils.arrayBufferToHex(await crypto.generateMasterHash(passwordInput.value, masterKey)) // Derive bits via pbkdf2 from the masterkey and the users password (this is used for server-side auth)

  let response = await fetch(window.backendRoutes.signin, {
    method: "POST",
    body: JSON.stringify({
      Email: emailInput.value,
      MasterHash: masterHash,
    })
  })
  let jsonResponse = await response.json()

  // jsonResponse.Authenticated is only used as a quick way to see if a user is authenticated, authentication is used server-side, this value means nothing
  if (jsonResponse.Authenticated) {
    window.location = window.frontendRoutes.home
  } else {
    window.location.reload()
  }
}

window.auth.signout = async () => {
  let response = await fetch(window.backendRoutes.signout, {
    method: "DELETE",
  })
  let jsonResponse = await response.json()

  if (jsonResponse.SignedOut) {
    window.location = window.frontendRoutes.signin
  }
}