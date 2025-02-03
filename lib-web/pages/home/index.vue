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
                        <el-link type="primary" :href="'/search?key_word='+searchKeyWord" @click="searchCheck">Search</el-link>
                    </el-menu-item>
                </el-menu>
            </el-header>

            <el-main></el-main>

            <el-footer>Footer</el-footer>
        </el-container>
    </div>
</template>

<script setup>
import { Search } from '@element-plus/icons-vue'

const activeIndex = ref('1')
let searchKeyWord = ref('')
</script>

<style>
el-header {
    background-color: rgb(102.2, 177.4, 255);
}
</style>