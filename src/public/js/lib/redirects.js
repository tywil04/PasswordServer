export function refreshPage() {
  window.location.reload()
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