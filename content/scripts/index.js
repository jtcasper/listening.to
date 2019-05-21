function getLoginState() {
  return Math.floor(Math.random() * 1000)
}

function appendLoginState() {
  const loginUrl = document.getElementById('login')
  loginUrl.href = loginUrl.href + '&state=' + getLoginState()
}

function getAccountCookie() {
  const listeningCookieName = "account_info"
  const cookies = document.cookie
  return cookies.split(';')
    .filter((item) => item.includes(`${listeningCookieName}=`))
    .map((cookie) => cookie.split('='))
    .reduce((prev, [key, value]) => (prev[key] = value,prev), {})
}

function checkCookieExists() {
  const cookie = getAccountCookie()
  return !(Object.keys(cookie).length === 0 && obj.constructor === Object)
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
