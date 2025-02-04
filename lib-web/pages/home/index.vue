<template>
                <div v-if="books.length > 0" class="book-item">
                    <BookItem v-for="book in books" :bookTitle="book.book_name" :bookCover="book.book_cover"
                        :bookHash="book.book_hash" />
                </div>
                <div v-else>
                    <div v-loading="true" element-loading-text="Loading...">
                        <br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br>
                    </div>
                </div>

</template>

<script setup>
import { BookItem } from '#components'

let searchKeyWord = ref('')
let bookResult = ref([])
let books = ref([])

const config = useRuntimeConfig()


const getBooks = async () => {
    try {
        const result = await fetch(`${config.public.apiProxy}/book/search/?keyword=${encodeURIComponent(searchKeyWord.value)}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })

        if (!result.ok) {
            throw new Error(`HTTP error! Status: ${result.status}`)
        }

        bookResult.value = await result.json()
        books.value = bookResult.value.search_result || []
        // console.log(books.value)
    } catch (error) {
        ElNotification({
            title: 'Error',
            message: 'Search failed!',
            type: 'error',
        })
        console.error('There was a problem with the fetch operation:', error)
        bookResult.value = []
        books.value = []
    }
}

onMounted(() => {
    getBooks()
})
</script>

<style>
.book-item {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
}
</style>