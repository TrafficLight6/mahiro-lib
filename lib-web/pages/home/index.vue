<template>
    <div class="common-layout">
        <el-container>
            <el-header>
                <el-menu :default-active="activeIndex" class="el-menu-demo" mode="horizontal" :ellipsis="false">
                    <el-menu-item index="0">
                        <el-link href="/home">MahiroLib</el-link>
                    </el-menu-item>
                    <el-menu-item>All Book</el-menu-item>
                    <el-menu-item>
                        <el-input :prefix-icon="Search" clearable v-model="searchKeyWord" />
                        &ensp;&ensp;
                        <el-button @click="search">Search</el-button>
                    </el-menu-item>
                </el-menu>
            </el-header>
            <el-main>
                <div v-if="books.length > 0" class="book-item">
                    <BookItem v-for="book in books" :bookTitle="book.book_name" :bookCover="book.book_cover" />
                </div>
                <div v-else>Loading</div>
            </el-main>
            <el-footer>Footer</el-footer>
        </el-container>
    </div>
</template>

<script setup>
import { Search } from '@element-plus/icons-vue'
import { BookItem } from '#components'

let searchKeyWord = ref('')
let activeIndex = ref('0')
let bookResult = ref([])
let books = ref([])

const router = useRouter()
const config = useRuntimeConfig()

const search = () => {
    if (searchKeyWord.value == '') {
        ElNotification({
            title: 'Error',
            message: 'Please type your search key word!',
            type: 'error',
        })
    } else {
        router.push({
            path: '/search',
            query: {
                keyword: searchKeyWord.value
            }
        })
    }

}
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
        console.log(books.value)
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