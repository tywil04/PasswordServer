import * as crypto from "/public/js/lib/crypto.js"
import * as utils from "/public/js/lib/utils.js"

async function signup() {
  const emailInput = document.querySelector("#email")
  const passwordInput = document.querySelector("#password")

  let masterKey = await crypto.generateMasterKey(passwordInput.value, emailInput.value) // Derive a key via pbkdf2 from the users password and email using
  let masterHash = utils.arrayBufferToHex(await crypto.generateMasterHash(passwordInput.value, masterKey)) // Derive bits via pbkdf2 from the masterkey and the users password (this is used for server-side auth)

  let databaseKey = await crypto.generateDatabaseKey() // generate random AES-256-CBC key
  let [iv, encryptedDatabaseKey] = await crypto.protectDatabaseKey(masterKey, databaseKey) // encrypt the key with masterkey
  let protectedDatabaseKey = utils.arrayBufferToHex(iv) + ";" + utils.arrayBufferToHex(encryptedDatabaseKey)

  let response = await fetch("/api/v1/auth/signup", {
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
    utils.redirectSignin()
    console.log(jsonResponse.UserId)
  } else {
    utils.refresh()
  }
}

async function signin() {
  const emailInput = document.querySelector("#email")
  const passwordInput = document.querySelector("#password")

  let masterKey = await crypto.generateMasterKey(passwordInput.value, emailInput.value) // Derive a key via pbkdf2 from the users password and email using
  let masterHash = utils.arrayBufferToHex(await crypto.generateMasterHash(passwordInput.value, masterKey)) // Derive bits via pbkdf2 from the masterkey and the users password (this is used for server-side auth)

  let response = await fetch("/api/v1/auth/signin", {
    method: "POST",
    body: JSON.stringify({
      Email: emailInput.value,
      MasterHash: masterHash,
    })
  })
  let jsonResponse = await response.json()

  // jsonResponse.Authenticated is only used as a quick way to see if a user is authenticated, authentication is used server-side, this value means nothing
  if (jsonResponse.Authenticated) {
    utils.redirectHome()
  } else {
    utils.refresh()
  }
}

// Export functions so they can be used
window.signup = signup
window.signin = signin