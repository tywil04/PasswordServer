import * as crypto from "/public/js/lib/crypto.js"
import * as utils from "/public/js/lib/utils.js"

async function signup(emailInput, passwordInput) {
  let masterKey = await crypto.generateMasterKey(passwordInput.value, emailInput.value)
  let masterHash = utils.arrayBufferToHex(await crypto.generateMasterHash(passwordInput.value, masterKey))

  let databaseKey = await crypto.generateDatabaseKey()
  let [iv, encryptedDatabaseKey] = await crypto.protectDatabaseKey(masterKey, databaseKey)
  let protectedDatabaseKey = utils.arrayBufferToHex(iv) + ";" + utils.arrayBufferToHex(encryptedDatabaseKey)

  let response = await fetch("/api/v1/auth/signup", {
    method: "POST",
//    credentials: "include",
    body: JSON.stringify({
      Email: emailInput.value,
      MasterHash: masterHash,
      ProtectedDatabaseKey: protectedDatabaseKey,
    })
  })
  let jsonResponse = await response.json()

  let success = jsonResponse.UserId !== undefined

  console.log(success)
  console.log(jsonResponse)
}

window.signup = signup