<template>
  <section class="vh-100">
    <div class="container-fluid">
      <div class="row">
        <div class="col-sm-6 text-black">
          <div class="px-5 ms-xl-4">
            <i class="fas fa-crow fa-2x me-3 pt-5 mt-xl-4" style="color: #709085;"></i>
            <span class="h1 fw-bold mb-0">Social Media App</span>
          </div>

          <div class="d-flex align-items-center h-custom-2 px-5 ms-xl-4 mt-5 pt-5 pt-xl-0 mt-xl-n5">
            <div class="p-4 border rounded shadow-sm bg-white" style="width: 100%; max-width: 23rem;">
              <!-- Login Form -->
              <form @submit.prevent="loginUser" aria-label="Login Form">
                <h3 class="fw-normal mb-3 pb-3 text-center" style="letter-spacing: 1px;">Log in</h3>

                <!-- Error Alert -->
                <div v-if="errorMessage" class="alert alert-danger text-center" role="alert" aria-live="polite">
                  {{ errorMessage }}
                </div>

                <div class="form-outline mb-4">
                  <input v-model="emailLogin" type="email" id="form2Example18" class="form-control form-control-lg" required />
                  <label class="form-label" for="form2Example18">Email address</label>
                </div>

                <div class="form-outline mb-4">
                  <input v-model="passwordLogin" type="password" id="form2Example28" class="form-control form-control-lg" required />
                  <label class="form-label" for="form2Example28">Password</label>
                </div>

                <div class="pt-1 mb-3">
                  <button class="btn btn-info btn-lg btn-block w-100 float-hover" type="submit" :disabled="loading" aria-label="Login Button">
                    <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                    <span v-if="!loading">Login</span>
                  </button>
                </div>

                <p class="text-center text-muted mb-3">
  or <a href="#" @click.prevent="microsoftLogin" class="link-primary">log in with your Microsoft account</a>
</p>


                <p class="small mb-2 text-center">
                  <a class="text-muted" href="#!">Forgot password?</a>
                </p>

                <p>
                  Don't have an account?
                  <a href="#" class="link-info" data-bs-toggle="modal" data-bs-target="#signupModal">Register here</a>
                </p>
              </form>
            </div>
          </div>
        </div>

        <div class="col-sm-6 px-0 d-none d-sm-block">
  <img :src="loginImg" alt="Login image"
    class="w-100 vh-100"
    style="object-fit: cover; object-position: left;" />
</div>

      </div>
    </div>
  </section>

  <!-- Sign Up Modal -->
  <div class="modal fade" id="signupModal" tabindex="-1" aria-labelledby="signupModalLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content p-4">
        <div class="modal-header">
          <h5 class="modal-title" id="signupModalLabel">Sign Up</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>

        <div class="modal-body">
          <!-- Sign Up Form -->
          <form @submit.prevent="registerUser" aria-label="Sign Up Form">
            <!-- Error Alert -->
            <div v-if="errorMessage" class="alert alert-danger text-center" role="alert" aria-live="polite">
              {{ errorMessage }}
            </div>

            <div class="form-outline mb-3">
              <input type="text" id="signupName" v-model="name" class="form-control" required />
              <label class="form-label" for="signupName">Full Name</label>
            </div>

            <div class="form-outline mb-3">
              <input type="email" id="signupEmail" v-model="email" class="form-control" required />
              <label class="form-label" for="signupEmail">Email</label>
            </div>

            <div class="form-outline mb-3">
              <input type="password" id="signupPassword" v-model="password" class="form-control" required />
              <label class="form-label" for="signupPassword">Password</label>
            </div>

            <div class="form-outline mb-3">
              <input type="password" id="signupConfirm" v-model="confirmPassword" class="form-control" required />
              <label class="form-label" for="signupConfirm">Confirm Password</label>
            </div>

            <button type="submit" class="btn btn-info w-100 float-hover" :disabled="loading" aria-label="Create Account Button">
              <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
              <span v-if="!loading">Create Account</span>
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import loginImg from '@/assets/agasagas.jpg'
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { Modal } from 'bootstrap'
import { PublicClientApplication, InteractionRequiredAuthError } from '@azure/msal-browser'


// State
const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const emailLogin = ref('')
const passwordLogin = ref('')
const loading = ref(false)
const errorMessage = ref('')

const router = useRouter()

// MSAL Configuration
const msalConfig = {
  auth: {
    clientId: '381360bf-45da-44b6-af71-9fc5b992f109', // 🔁 Replace with your Azure AD App's client ID
    authority: 'https://login.microsoftonline.com/678a4840-b3b3-4cd7-8cee-1381def5b0b9', // Use tenant-specific if needed
    redirectUri: 'http://localhost:5173', // ✅ Replace with your app URL
  }
}

const msalInstance = new PublicClientApplication(msalConfig)

const loginRequest = {
  scopes: ['user.read'] // Change scopes based on what you need
}

// Microsoft Login Function
// Microsoft Login Function (Enhanced)
const microsoftLogin = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    await msalInstance.initialize()

    let account = msalInstance.getAllAccounts()[0]

    // If no account, initiate login
    if (!account) {
      const loginResponse = await msalInstance.loginPopup(loginRequest)
      account = loginResponse.account!
    }

    // Set account as active
    msalInstance.setActiveAccount(account)

    // Attempt to silently acquire an access token
    const tokenResponse = await msalInstance.acquireTokenSilent({
      ...loginRequest,
      account,
    }).catch(async (error) => {
      // If silent token acquisition fails, use interactive login
      if (error instanceof InteractionRequiredAuthError) {
        return await msalInstance.acquireTokenPopup(loginRequest)
      }
      throw error
    })

    const idToken = tokenResponse.idToken
    const accessToken = tokenResponse.accessToken

    // Send ID token to backend
    const res = await axios.post('http://localhost:8081/auth/microsoft', {
      token: idToken,
    })

    alert(res.data.message || 'Logged in with Microsoft successfully!')
    localStorage.setItem('user', JSON.stringify(res.data.user))
    localStorage.setItem('ms_access_token', accessToken)
    router.push('/dashboard')
  } catch (err: any) {
    console.error('Microsoft login failed:', err)
    errorMessage.value = err?.message || 'Microsoft login failed.'
  } finally {
    loading.value = false
  }
}



// Regular Email/Password Login
const loginUser = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await axios.post('http://localhost:8081/login', {
      email: emailLogin.value,
      password: passwordLogin.value
    })

    alert(res.data.message || 'Logged in successfully!')
    localStorage.setItem('user', JSON.stringify(res.data.user))
    router.push('/dashboard')
  } catch (err: unknown) {
    errorMessage.value = axios.isAxiosError(err)
      ? err.response?.data?.error || 'Login failed.'
      : 'Something went wrong.'
  } finally {
    loading.value = false
  }
}

// Registration
const registerUser = async () => {
  if (password.value !== confirmPassword.value) {
    errorMessage.value = "Passwords don't match!"
    return
  }

  loading.value = true
  errorMessage.value = ''
  try {
    const res = await axios.post('http://localhost:8081/signup', {
      name: name.value,
      email: email.value,
      password: password.value
    })

    alert(res.data.message || 'Account created successfully!')

    name.value = ''
    email.value = ''
    password.value = ''
    confirmPassword.value = ''

    const modal = document.getElementById('signupModal') as HTMLElement | null
    if (modal) {
      const bsModal = Modal.getInstance(modal) || new Modal(modal)
      bsModal.hide()
    }
  } catch (err: unknown) {
    errorMessage.value = axios.isAxiosError(err)
      ? err.response?.data?.error || 'Registration failed.'
      : 'Something went wrong.'
  } finally {
    loading.value = false
  }
}
</script>


<style scoped>
.float-hover {
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}
.float-hover:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
}
</style>
