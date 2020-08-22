import axios from 'axios'
//export const restPath = 'http://localhost:3230/api/'
//export const restPath = 'https://localhost:8000/api/'
export const restPath = '/api/'

export function authHeader() {
  var tokens = JSON.parse(localStorage.getItem("tokens"))
  if (tokens && tokens.token) {
    return {
      'Authorization': `Bearer ${tokens.token}`
    }
  }
  return {}
}

export const request = ( method, path, data = {} ) => {
  
  const options = {
    method,
    data,
    headers: authHeader(),
    url: restPath + path,
    timeout: 50000,
  }
  return axios(options)
}
