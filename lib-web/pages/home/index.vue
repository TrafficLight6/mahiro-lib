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
                        <el-input :prefix-icon="Search" clearable v-model="searchKeyWord"/>
                        &ensp;&ensp;
                        <el-button @click="search">Search</el-button>
                    </el-menu-item>
                </el-menu>
            </el-header>
            <el-main>

            </el-main>

            <el-footer>Footer</el-footer>
        </el-container>
    </div>
</template>

<script setup>
import { Search } from '@element-plus/icons-vue'
let searchKeyWord = ref('')
let activeIndex = ref('0')

const router = useRouter()

const search = () => {
    if (searchKeyWord.value=='') {
        ElNotification({
            title: 'Error',
            message: 'Please type your search key word!',
            type: 'error',
        })
    }else{
        router.push({
            path: '/search',
            query: {
                keyword: searchKeyWord.value
            }
        })
    }
}
</script>

<style></style>