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
const user_check: string = await $fetch(config.api_proxy + '/user/check?token=' + cookie.value)
const user_check_json = JSON.parse(user_check)

if (guest_allow_json.resultMessage != 'true' && user_check_json.success) {
  router.push('/login')
} else {
  router.push('/home')
}
</script>