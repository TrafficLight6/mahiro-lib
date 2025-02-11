<template>
    <template>
        <div class="container">
            <el-container>
                <el-aside width="600px"></el-aside>
                <el-main>
                    <div v-if="bookPictureList != {} && bookType.value != ''">
                        <div v-if="bookType == 'picture'">
                            <ReadingPicture :book-picture-list="bookPictureList" />
                        </div>
                        <div v-else-if="bookType == 'text'">
                            <h1>text,here</h1>
                        </div>
                    </div>
                </el-main>
                <el-aside width="600px"></el-aside>
            </el-container>
        </div>
    </template>
</template>

<script setup>
const route = useRoute()
const router = useRouter()
const config = useRuntimeConfig()

let bookPictureListResult = ref({})
let bookPictureList = ref([])
let chapterHash = ref(route.params.chapterHash)

let bookInfoResult = ref({})
let bookType = ref('')
let bookHash = ref(route.params.bookHash)

const getBookInfo = async () => {
    try {
        const result = await fetch(`${config.public.apiProxy}/book/get?book_hash=${encodeURIComponent(bookHash.value)}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })

        if (!result.ok) {
            throw new Error(`HTTP error! Status: ${result.status}`)
        }
        bookInfoResult.value = await result.json()
        bookType.value = bookInfoResult.value.book.book_type || ''
        console.log(bookInfoResult.value)
    } catch (error) {
        ElNotification({
            title: 'Error',
            message: 'Get book info failed!',
            type: 'error',
        })
        console.error('There was a problem with the fetch operation:', error)
        bookInfoResult.value = {}
        bookType.value = ''
    }
}

const getBookChapterList = async () => {
    try {
        const result = await fetch(`${config.public.apiProxy}/book/chapter/get?chapter_hash=${encodeURIComponent(chapterHash.value)}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })

        if (!result.ok) {
            throw new Error(`HTTP error! Status: ${result.status}`)
        }
        bookPictureListResult.value = await result.json() || {}
        bookPictureList.value = bookPictureListResult.value.chapter.chapter_file_list || []
        // console.log(bookPictureListResult.value)
        // console.log(bookPictureList.value)
    } catch (error) {
        ElNotification({
            title: 'Error',
            message: 'Get chapter context failed!',
            type: 'error',
        })
        console.error('There was a problem with the fetch operation:', error)
        bookPictureListResult.value = {}
        bookPictureList.value = []
    }
}

onMounted(() => {
    getBookInfo()
    getBookChapterList()
})
</script>

<style>
.container {
    text-align: center;
}
</style>