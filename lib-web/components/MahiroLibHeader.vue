<template>
    <el-affix :offset="0">
        <el-menu :default-active="activeIndex" class="el-menu-demo" mode="horizontal" :ellipsis="false" :offset="0">
            <el-menu-item index="0">
                <el-link href="/home">MahiroLib</el-link>
            </el-menu-item>
            <el-menu-item>
                <el-link href="/home">Home</el-link>
            </el-menu-item>
            <el-menu-item>
                <el-input :prefix-icon="Search" clearable v-model="searchKeyWord" />
                &ensp;&ensp;
                <el-button @click="search">Search</el-button>
            </el-menu-item>
        </el-menu>
    </el-affix>
</template>

<script setup>
import { Search } from '@element-plus/icons-vue'
const router = useRouter()

let searchKeyWord = ref(router.currentRoute.value.query.keyword ||  '')
let activeIndex = ref('0')



const search = () => {
    if (searchKeyWord.value == '') {
        ElNotification({
            title: 'Error',
            message: 'Please type your search key word!',
            type: 'error',
        })
    } else {
        router.push({
            path: '/home/search',
            query: {
                keyword: searchKeyWord.value
            }
        })
    }

}

</script>