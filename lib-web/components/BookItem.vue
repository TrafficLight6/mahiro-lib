<template>
    <div class="container">
        <el-link :href="`/home/book/` + bookHash">
            <div class="book-item">

                <div class="book-item-img">
                    <el-image fit="cover" style="height: 424px;width: 300px;" :src="bookCover" />
                </div>
                <div class="book-item-info">
                    <div class="book-item-title">
                        <!-- <p>{{ cutString(bookTitle,20) }}</p> -->
                        <p v-html="cutString(bookTitle, 20)"></p>
                    </div>
                </div>

            </div>
        </el-link>

    </div>
</template>

<script setup>
const config = useRuntimeConfig()

const props = defineProps({ bookTitle: { type: String }, bookCover: { type: String }, bookHash: { type: String } })
const bookTitle = props.bookTitle
let bookCover = props.bookCover.replaceAll('{proxy}', config.public.apiProxy)

const cutString = (rawString, length) => {
    // return rawString.length > length ? rawString.substring(0, length) + '...' : rawString;
    if (rawString.length > length) {
        if (rawString.length > length * 2) {
            return rawString.substring(0, length) + '<br>' + rawString.substring(length, length * 2) + '...'
        } else {
            return rawString.substring(0, length) + '<br>' + rawString.substring(length)

        }
    } else {
        return rawString;
    }
}


</script>

<style scoped>
.book-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    transition: box-shadow 0.3s ease;
}

.book-item:hover {
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}

.book-item-img {
    width: 100%;
}

.book-item-info {
    width: 100%;
    text-align: center;
}
</style>