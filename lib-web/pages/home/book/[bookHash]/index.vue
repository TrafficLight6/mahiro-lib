<template>
<!-- 



Image
                        List...

Name






-->
</template>


<script setup>
const router = useRoute()
const bookHash = router.params.bookHash
if (!bookHash){
    useRouter().push('/')
}
</script>