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
                <el-button type="primary" class="login-button" @click="login" :loading="loading">Login</el-button>
                <br><br>
            </el-form>
        </div>
    </div>
</template>

<script setup>
import Cookies from 'js-cookie'

let username = ref('')
let password = ref('')
let remember = ref(false)
let loading = ref(false)

onMounted(() => {
    if (Cookies.get('token') != undefined) {
        useRouter().push('/home')
    }
})

let login = () => {
    let config = useRuntimeConfig()

    loading.value = true
    if (username.value == '' || password.value == '') {
        ElNotification({
            title: 'Error',
            message: 'Please type your username and password!',
            type: 'error',
        })
        loading.value = false
        return
    } else {
        let xhr = new XMLHttpRequest();
        xhr.onreadystatechange = function () {
            if (xhr.readyState == 4 && xhr.status == 200) {
                let result = JSON.parse(xhr.responseText);
                if (result.success) {
                    if (remember.value == 'true') {
                        Cookies.set('token', result.token, { expires: 9999 })
                    } else {
                        Cookies.set('token', result.cookie)
                    }
                    ElNotification({
                        title: 'Success',
                        message: 'Login Success!',
                        type: 'success',
                    })
                    useRouter().push('/')
                } else {
                    ElNotification({
                        title: 'Error',
                        message: 'Username or password error!',
                        type: 'error',
                    })
                    loading.value = false
                }
            }
        }
        xhr.open('POST',config.public.apiProxy+"/user/login/"+"?username="+username.value+"&password="+password.value+"&remember="+remember.value , true);
        xhr.send();
    }
}

</script>

<!-- <script>
import Cookies from 'js-cookie'

if (Cookies.get('token')!=undefined){
    useRouter().push('/home')
}
</script> -->
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