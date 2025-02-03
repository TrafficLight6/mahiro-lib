<template>
    <div class="container">
        <div class="book-item">
            <div class="book-item-img">
                <el-image fit="cover" style="height: 424px;width: 300px;" :src="bookCover" />
            </div>
            <div class="book-item-info">
                <div class="book-item-title">{{bookTitle}}</div>
            </div>
        </div>
    </div>
</template>

<script setup>
    const config= useRuntimeConfig()

    const props = defineProps({bookTitle: {type: String},bookCover: {type: String}})
    const bookTitle = props.bookTitle
    let bookCover = props.bookCover.replaceAll('{proxy}',config.public.apiProxy)
</script>

<style scoped>
.book-item {
    display: flex;
    flex-direction: column;
    align-items: center;
}
.book-item-img {
    width: 100%;
}
.book-item-info {
    width: 100%;
    text-align: center;
}
</style>