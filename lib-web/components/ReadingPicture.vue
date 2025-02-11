<template>
    <div class="reading-container">
        <div class="book-picture-list" v-for="(bookPicture, index) in props.bookPictureList">
            <el-image :src="bookPicture" v-show="index == pageNumber"></el-image>
        </div>
        <div class="button">
            <el-button @click="prevPage">prev page</el-button>
            <el-button @click="nextPage">next page</el-button>
        </div>
    </div>
</template>

<script setup>
let pageNumber = ref(0)
let props = defineProps({
    bookPictureList: { type: Array, default: () => [] }
})

// const nextPageDiraction = ref("right")

const prevPage = () => {
    if (pageNumber.value == 0) {
        return
    }
    pageNumber.value--
}

const nextPage = () => {
    if (pageNumber.value == props.bookPictureList.length - 1) {
        return
    }
    pageNumber.value++
}
</script>


<style>
.button {
    display: flex;
    justify-content: space-between;
}

.reading-container {
    width: 70%;
}
</style>