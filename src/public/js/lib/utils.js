export function stringToArrayBuffer(string) {
  return new TextEncoder().encode(string)
}

export function arrayBufferToHex(byteArray) {
  return [...new Uint8Array(byteArray)].map(x => x.toString(16).padStart(2, '0')).join('');
}

export function hexToArrayBuffer(hex) {
  var typedArray = new Uint8Array(hex.match(/[\da-f]{2}/gi).map((h) => parseInt(h, 16)))
  return typedArray.buffer
}

export function refreshPage() {
  window.location = window.location // hacky way to refresh page
}

export function redirectSignin() {
  window.location = "/auth/signin"
}

export function redirectSignup() {
  window.location = "/auth/signup"
}

export function redirectHome() {
  window.location = "/home"
}