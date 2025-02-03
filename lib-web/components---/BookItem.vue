<template>
    <div class="container">
        <div class="book-item">
            <div class="book-item-img">
                <el-image fit="cover" style="height: 424px;width: 300px;" :src="bookItem.book_cover" />
            </div>
            <div class="book-item-info">
                <div class="book-item-title">{{ bookItem.book_name }}</div>
            </div>
        </div>
    </div>
</template>

<script setup>
const props = defineProps({
    book: Object
})
const bookItem = props.book
const config = useRuntimeConfig()

let cover = bookItem.book_cover.prototype.replaceAll("{proxy}", config.public.apiProxy)
bookItem.book_cover = cover
</script>