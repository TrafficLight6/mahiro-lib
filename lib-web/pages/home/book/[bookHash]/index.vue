<template>
    <el-container v-if="bookInfo.value != {}">
        <el-row :gutter="20">
            <el-col :span="12">
                <div class="left-content" style="height: 50%; width: 50%;">
                    <el-image fit="cover" :src="bookCover" />
                </div>
            </el-col>
            <el-col :span="12">
                <div class="right-content">
                    <h2>{{ bookTitle }}</h2>
                    <!-- add book info-->

                    <div class="chapter">
                        <p>Chapter List</p>
                        <div v-if="BookChapterList == null || BookChapterList.length == 0">
                            <h1>Sorry,no chapter here</h1>
                        </div>
                        <div v-else>
                            <div v-for="(chapter, index) in BookChapterList">
                                <br>
                                <el-link :href="`/home/book/${bookHash}/read/${chapter.chapter_hash}`"><el-button
                                        style="width: 600px;">{{ index + 1 }}{{ " . " + chapter.chapter_name
                                        }}</el-button>
                                </el-link>
                            </div>
                        </div>
                    </div>
                </div>
            </el-col>
        </el-row>
    </el-container>
    <div v-else v-loading="ture">
        <br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br>
    </div>
</template>

<style scoped>
.left-content,
.right-content {
    padding: 20px;
}

.left-content {
    margin: 0 auto;
}
</style>

<script setup>
const router = useRoute()
const bookHash = router.params.bookHash
const config = useRuntimeConfig()

let bookInfoResult = ref({})
let bookInfo = ref({})

let bookCover = ref('')
let bookTitle = ref('')

let BookChapterListResult = ref({})
let BookChapterList = ref([])

if (!bookHash) {
    useRouter().push('/')
}

const getBookInfo = async () => {
    try {
        const result = await fetch(`${config.public.apiProxy}/book/get?book_hash=${encodeURIComponent(bookHash)}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })

        if (!result.ok) {
            throw new Error(`HTTP error! Status: ${result.status}`)
        }
        bookInfoResult.value = await result.json()
        bookInfo.value = bookInfoResult.value || {}

        bookCover.value = bookInfo.value.book.book_cover.replaceAll('{proxy}', config.public.apiProxy)
        bookTitle.value = bookInfo.value.book.book_name
        // console.log(bookInfo.value)
        // console.log(bookCover.value)
        // console.log(bookTitle.value)
    } catch (error) {
        ElNotification({
            title: 'Error',
            message: 'Get book info failed!',
            type: 'error',
        })
        console.error('There was a problem with the fetch operation:', error)
        bookInfoResult.value = {}
        bookInfo.value = {}
    }
}
const getBookChapterList = async () => {
    try {
        const result = await fetch(`${config.public.apiProxy}/book/chapter/getlist?book_hash=${encodeURIComponent(bookHash)}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })

        if (!result.ok) {
            throw new Error(`HTTP error! Status: ${result.status}`)
        }
        BookChapterListResult.value = await result.json() || {}
        BookChapterList.value = BookChapterListResult.value.result || []
        console.log(BookChapterList.value)
    } catch (error) {
        ElNotification({
            title: 'Error',
            message: 'Get book chapter list failed!',
            type: 'error',
        })
        console.error('There was a problem with the fetch operation:', error)
        BookChapterListResult.value = {}
        BookChapterList.value = []
    }
}

const handleChapterClick = (chapter) => {
    try {
        route.push(`/home/book/${bookHash}/chapter/read/${chapter.chapter_id}`)
    } catch (error) {
        ElNotification({
            title: 'Error',
            message: 'Failed to navigate to chapter!',
            type: 'error',
        })
        console.error('There was a problem with navigation:', error)
    }
}

onMounted(() => {
    getBookInfo()
    getBookChapterList()
})
</script>