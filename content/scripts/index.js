function getLoginState() {
  return Math.floor(Math.random() * 1000)
}

function appendLoginState() {
  const loginUrl = document.getElementById('login')
  loginUrl.href = loginUrl.href + '&state=' + getLoginState()
}
