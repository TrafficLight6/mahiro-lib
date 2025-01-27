<template>
    <div class="container">
        <div class="centered-content">
            <h1>Login the Mahiro lib!</h1>
            <el-form>
                <el-input v-model="username" style="width: 240px" placeholder="Please type your username" clearable />
                <br><br>
                <el-input v-model="password" style="width: 240px" type="password"
                    placeholder="Please type your password" show-password clearable />
                <br></br>
                <el-checkbox v-model="remember">Remember me</el-checkbox>
                <br><br>
                <el-button type="primary" class="login-button">Login</el-button>
                <br><br>
            </el-form>
        </div>
    </div>
</template>

<script setup>
let username = ref('')
let password = ref('')
let remember = ref(false)
</script>

<style scoped>
.container {
    background: #ffffff;
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
}
.centered-content {
    text-align: center;
    /* border: 2px dashed #000; */
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    /* 添加阴影边框 */
}
</style>