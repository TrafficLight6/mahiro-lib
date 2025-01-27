<template>
  <div>
    <NuxtPage />
  </div>
</template>

<script setup lang="ts">
const config = useRuntimeConfig()
const guest_allow: string = await $fetch(config.api_proxy + '/config/get?key=guest-allow')
const guest_allow_json = JSON.parse(guest_allow)
const router = useRouter()
const cookie = useCookie('token')



if (guest_allow_json.resultMessage!='true'){
  router.push('/login')
}
</script>