function getLoginState() {
  return Math.floor(Math.random() * 1000)
}

function appendLoginState() {
  const loginUrl = document.getElementById('login')
  loginUrl.href = loginUrl.href + '&state=' + getLoginState()
}

function checkCookieExists() {
  const listeningCookieName = "account_info"
  const cookies = document.cookie
  return cookies.split(';').filter((item) => item.includes('account_info=')).length >= 1
}

function hideLoginUrl() {
  const loginUrl = document.getElementById('login')
  loginUrl.style.display = 'none'
}

function callAnalyze() {
  const analyzeUrl = 'analyze'
  return fetch(analyzeUrl)
  .then(data => data.json())
  .then(res => res)
}
