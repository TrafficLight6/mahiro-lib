<template>
    <div class="common-layout">
        <el-container>
            <el-header>
                <MahiroLibHeader/>
            </el-header>
            <el-main>
                <NuxtPage/>
            </el-main>
            <el-footer class="footer">
                <MahiroLibFooter/>
            </el-footer>
        </el-container>
    </div>
</template>

<script setup>
import { MahiroLibHeader } from '#components';
import { MahiroLibFooter } from '#components';
</script>

<style>

.footer {
    background-color:#ff8da6;
    text-align: center;
}

html body {
    margin: 0px;
    padding: 0px;
}
</style>