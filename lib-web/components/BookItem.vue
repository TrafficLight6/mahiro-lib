<template>
    <div class="container">
        <div class="book-item">

            <div class="book-item-img">
                <el-link :href="`/book/`+bookHash">
                    <el-image fit="cover" style="height: 424px;width: 300px;":src="bookCover" />
                </el-link>
            </div>
            <div class="book-item-info">
                <div class="book-item-title">{{ bookTitle }}</div>
            </div>

        </div>
    </div>
</template>

<script setup>
const config = useRuntimeConfig()

const props = defineProps({ bookTitle: { type: String }, bookCover: { type: String }, bookHash: { type: String } })
const bookTitle = props.bookTitle
const bookHerf = props.bookHash
let bookCover = props.bookCover.replaceAll('{proxy}', config.public.apiProxy)
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