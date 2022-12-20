import * as utils from "{{ .FrontendRoutesJS.utils }}"

export async function importRawKey(key, extractable=false, algorithm="PBKDF2", usages=["deriveKey", "deriveBits"]) {
  return await crypto.subtle.importKey(
    "raw",
    key,
    algorithm,
    extractable,
    usages,
  )
}

export async function exportRawKey(key) {
  return await crypto.subtle.exportKey("raw", key)
}

export async function generateMasterKey(masterPassword, email) {
  return await crypto.subtle.deriveKey(
    {
      name: "PBKDF2",
      hash: "SHA-512",
      salt: utils.stringToArrayBuffer(email),
      iterations: 100000
    },
    await importRawKey(utils.stringToArrayBuffer(masterPassword)),
    {name: "AES-CBC", length: 256},
    true,
    ["wrapKey", "unwrapKey"]
  );
}

export async function generateMasterHash(masterPassword, masterKey) {
  const masterHash = await crypto.subtle.deriveBits(
    {
      name: "PBKDF2",
      hash: "SHA-512",
      salt: utils.stringToArrayBuffer(masterPassword),
      iterations: 100000
    },
    await importRawKey(await exportRawKey(masterKey)),
    512
  )
  return masterHash
}

export async function generateDatabaseKey() {
  return await crypto.subtle.generateKey(
    {name: "AES-CBC", length: 256},
    true,
    ["encrypt", "decrypt"]
  )
}

export async function protectDatabaseKey(masterKey, databaseKey) {
  const iv = crypto.getRandomValues(new Uint8Array(16))
  const protectedDatabaseKey = await crypto.subtle.wrapKey(
    "raw", 
    databaseKey, 
    masterKey, 
    {name: "AES-CBC", iv}
  )
  return [iv, protectedDatabaseKey]
}

export async function unprotectDatabaseKey(masterKey, protectedDatabaseKey) {
  const [iv, wrappedDatabaseKey] = encryptedDatabaseKey
  return await crypto.subtle.unwrapKey(
    "raw", 
    wrappedDatabaseKey, 
    masterKey, 
    {name: "AES-CBC", iv},
    {name: "AES-CBC", length: 256},
    false,
    ["encrypt", "decrypt"]
  )
}