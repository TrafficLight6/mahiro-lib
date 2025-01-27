// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  modules: [
    '@element-plus/nuxt',
    // '@element-plus/icons-vue'
  ],
  runtimeConfig :{
    api_proxy : 'http://127.0.0.1:7622/'
  }
})
