<template>
    <div v-if="books.length > 0" class="book-item">
        <BookItem v-for="book in books" :bookTitle="book.book_name" :bookCover="book.book_cover"
            :bookHash="book.book_hash" />
    </div>
    <div v-else class="nothing">
            <br><br><br><br><br><br><br><br><br><br><br><br><br><br><h1>Sorry,Nothing here...</h1><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br>
    </div>

</template>

<script setup>
import { BookItem } from '#components'

let bookResult = ref([])
let books = ref([])

const config = useRuntimeConfig()
const router =useRouter()
const routerForUpdate = useRouter()

useHead({
    title: 'Mahiro Lib - Search:'+router.currentRoute.value.query.keyword
})

let searchKeyWord = ref(router.currentRoute.value.query.keyword)
console.log(searchKeyWord.value)
const getBooks = async () => {
    try {
        const result = await fetch(`${config.public.apiProxy}/book/search/?key=${searchKeyWord.value}`, {
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
        // console.log(bookResult.value)
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

watch(
    () => routerForUpdate.currentRoute.value.query.keyword,
    (newKeyword) => {
        searchKeyWord.value = newKeyword || ''
        getBooks()
    }
)
</script>

<style>
.book-item {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
}

.nothing{
    text-align: center;
}
</style>