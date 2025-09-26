// src/api.js
import axios from 'axios'

// Vite의 base 설정값을 그대로 사용
const index = axios.create({
  baseURL: import.meta.env.BASE_URL, // vite.config.js > defineConfig.base
  withCredentials: true // 필요 시 세션 쿠키 등 포함
})

export default index
