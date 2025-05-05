<template>
  <!-- Font Awesome CDN (for hand icons) -->
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.4/css/all.min.css">

  <div>
    <!-- Header: Navbar -->
    <header>
      <nav class="navbar navbar-expand-lg navbar-light bg-light shadow-sm">
        <div class="container-fluid">
          <a class="navbar-brand fw-bold" href="#">Social Media App</a>
          <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav"
            aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
          </button>
          <div class="collapse navbar-collapse" id="navbarNav">
            <ul class="navbar-nav ms-auto">
              <li class="nav-item">
                <button class="btn btn-outline-danger" @click="confirmLogout" aria-label="Logout">
                  Logout
                </button>
              </li>
            </ul>
          </div>
        </div>
      </nav>
    </header>

    <!-- Main Dashboard -->
    <main class="container my-5">
      <section v-if="user" class="mb-4">
        <h1>Welcome, {{ user.name }}!</h1>
        <p><strong>Email:</strong> {{ user.email }}</p>
      </section>

      <section v-else class="text-center my-5">
        <p class="text-danger fs-5">You need to log in to view your dashboard.</p>
      </section>




<!-- Active Users and Multiple Upload Button -->
<div class="d-flex justify-content-end mb-4">
  <button class="btn btn-primary shadow-sm px-4 py-2" @click="viewActiveUsers">
    View Active Users
  </button>
<!-- View Active Users Modal -->
<div class="modal fade" id="activeUsersModal" tabindex="-1" aria-labelledby="activeUsersModalLabel" aria-hidden="true">
  <div class="modal-dialog modal-fullscreen">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="activeUsersModalLabel">Active Users</h5>
        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
      </div>

      <div class="modal-body">
        <!-- Facebook-style Feed -->
        <div class="feed-container" v-if="users.length">
          <div
            class="feed-card mb-4 p-3 bg-white rounded shadow-sm"
            v-for="user in users"
            :key="user.id"
          >
            <!-- User Header: Avatar + Name -->
            <div class="d-flex align-items-center mb-3">
              <img
                :src="user.photo_url || '/default-avatar.png'"
                class="rounded-circle me-3"
                alt="User Avatar"
                width="50"
                height="50"
                style="object-fit: cover;"
              />
              <div>
                <h6 class="mb-0">{{ user.name }}</h6>
                <small class="text-muted">Active now</small>
              </div>
            </div>

            <!-- User Posts -->
            <div v-if="user.posts && user.posts.length">
              <div class="row g-2">
                <div class="col-4" v-for="(photo, index) in user.posts" :key="index">
                  <img
                    :src="photo.photo_url"
                    class="img-fluid rounded image-preview"
                    alt="Post Photo"
                    style="object-fit: cover; height: 100px; width: 100%; cursor: pointer;"
                    @click="openImagePreview(photo.photo_url)"
                  />

                  <!-- Like/Dislike Icons -->
                  <div class="d-flex justify-content-between align-items-center mt-2">
                    <button
                      class="btn btn-outline-primary btn-sm"
                      @click="likePost(photo.id, user.id)"
                      title="Like this post"
                    >
                      👍 Like
                    </button>
                    <button
                      class="btn btn-outline-danger btn-sm"
                      @click="dislikePost(photo.id, user.id)"
                      title="Dislike this post"
                    >
                      👎 Dislike
                    </button>
                  </div>

                  <!-- Like/Dislike Counters -->
                  <div class="text-center mt-1">
                    <small class="text-muted">
                      👍 {{ photo.likes || 0 }} | 👎 {{ photo.dislikes || 0 }}
                    </small>
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="text-muted text-center">No posts to show</div>

          </div>
        </div>

        <!-- Fallback: No Users -->
        <div v-else class="text-center py-4">
          <p class="text-danger">No active users found.</p>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- Image Preview Modal -->
<div class="modal fade" id="imagePreviewModal" tabindex="-1" aria-labelledby="imagePreviewModalLabel" aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="imagePreviewModalLabel">Image Preview</h5>
        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
      </div>
      <div class="modal-body">
        <img :src="selectedImage" class="img-fluid" alt="Selected Image" />
      </div>
    </div>
  </div>
</div>






  <!-- Multiple Uploads Button -->
  <button class="btn btn-success shadow-sm px-4 py-2 ms-3" @click="showUploadModal = true">
    Upload Multiple Photos
  </button>
</div>


<!-- Modal for Uploading Multiple Photos -->
<div v-if="showUploadModal" class="modal-overlay" @click="closeUploadModal">
  <div class="modal-content" @click.stop>
    <h5>Upload Multiple Photos</h5>
    <form @submit.prevent="uploadFiles">
      <!-- File Input -->
      <div class="mb-3">
        <label for="photos" class="form-label">Select Photos</label>
        <input type="file" id="photos" multiple @change="handleFileChange" class="form-control" />
      </div>

      <!-- Image Previews -->
      <div v-if="previewImages.length">
        <h6>Preview Selected Images</h6>
        <div class="d-flex flex-wrap">
          <img v-for="(image, index) in previewImages" :key="index" :src="image" class="img-thumbnail" width="100" />
        </div>
      </div>

      <div class="mt-3 d-flex justify-content-between">
        <button type="button" class="btn btn-secondary" @click="closeUploadModal">Cancel</button>
        <button type="submit" class="btn btn-primary">Upload</button>
      </div>
    </form>
  </div>
</div>



                <!-- About Section -->
      <section class="section about-section gray-bg py-5" id="about">
        <div class="container">
          <div class="row align-items-center flex-row-reverse">
            <!-- Left: Info -->
            <div class="col-lg-6 mb-4 mb-lg-0">
              <div class="about-text go-to">
                <h3 class="dark-color">About Me</h3>
                <p>{{ user?.description || 'No description provided yet.' }}</p>
                <div class="row about-list">
                  <div class="col-md-6">
                    <div class="media">
                      <label>Title:</label>
                      <p>{{ user?.title }}</p>
                    </div>
                    <div class="media">
                      <label>Birthday:</label>
                      <p>{{ user?.birthday }}</p>
                    </div>
                    <div class="media">
                      <label>Age:</label>
                      <p>{{ user?.age }}</p>
                    </div>
                    <div class="media">
                      <label>Address:</label>
                      <p>{{ user?.address }}</p>
                    </div>
                  </div>
                  <div class="col-md-6">
                    <div class="media">
                      <label>Email:</label>
                      <p>{{ user?.email }}</p>
                    </div>
                    <div class="media">
                      <label>Phone:</label>
                      <p>{{ user?.phone }}</p>
                    </div>
                    <div class="media">
                      <label>Skype:</label>
                      <p>{{ user?.skype }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

<!-- Right: Avatar and Button -->
<div class="col-lg-6 text-center">
  <div class="about-avatar mb-3">
    <!-- Display the avatar or a default image if none exists -->
    <img
      :src="user?.photo_url ? user.photo_url : '/default-avatar.png'"
      alt="User Avatar"
      class="profile-avatar"
    />
  </div>
  <button class="btn btn-primary" @click="showModal = true">
    Edit My Profile
  </button>
</div>

          </div>
        </div>

<!-- Photo Album Section -->
<!-- Show loading state -->
<div v-if="isLoading">Loading photos...</div>

<!-- No photos message -->
<div v-else-if="photos.length === 0">No photos available; Upload?.</div>

<!-- Thumbnails with delete button below -->
<div v-else class="photo-thumbnail-container d-flex flex-wrap gap-3">
  <div
    v-for="(photo, index) in photos"
    :key="index"
    class="photo-thumbnail text-center"
    style="position: relative; width: 200px;"
  >
    <img
      :src="photo"
      alt="User Photo"
      class="img-fluid rounded shadow-sm"
      style="width: 100%; height: auto; object-fit: cover;"
    />
    <button
      class="btn btn-outline-danger btn-sm mt-2"
      @click="deletePhoto(photo)"
    >
      <i class="bi bi-trash"></i> Delete
    </button>
  </div>
</div>


      </section>




    </main>


   <!-- Modal for Editing Profile -->
<div v-if="showModal" class="modal-overlay" @click="showModal = false">
  <div class="modal-content" @click.stop>
    <h5>Edit Profile</h5>
    <form @submit.prevent="saveProfile">
      <div class="mb-3">
        <label for="title" class="form-label">Title</label>
        <input v-model="profileData.title" type="text" class="form-control" id="title" required>
      </div>
      <div class="mb-3">
        <label for="birthday" class="form-label">Birthday</label>
        <input v-model="profileData.birthday" type="date" class="form-control" id="birthday" required>
      </div>
      <div class="mb-3">
        <label for="age" class="form-label">Age</label>
        <input v-model="profileData.age" type="number" class="form-control" id="age" required>
      </div>
      <div class="mb-3">
        <label for="address" class="form-label">Address</label>
        <input v-model="profileData.address" type="text" class="form-control" id="address" required>
      </div>
      <div class="mb-3">
        <label for="phone" class="form-label">Phone</label>
        <input v-model="profileData.phone" type="tel" class="form-control" id="phone" required>
      </div>
      <div class="mb-3">
        <label for="skype" class="form-label">Skype</label>
        <input v-model="profileData.skype" type="text" class="form-control" id="skype" required>
      </div>
      <!-- Description Field -->
      <div class="mb-3">
        <label for="description" class="form-label">Description</label>
        <textarea v-model="profileData.description" class="form-control" id="description" rows="3" required></textarea>
      </div>
      <div class="mb-3">
        <label for="avatar" class="form-label">Avatar</label>
        <input type="file" class="form-control" id="avatar" @change="previewImageHandler">
        <img v-if="previewImage" :src="previewImage" class="img-thumbnail mt-2" width="150" alt="Image Preview" />
      </div>
      <button type="submit" class="btn btn-primary">Save Changes</button>
    </form>
  </div>
</div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { Modal } from 'bootstrap'

// Interface for User
interface User {
  id: string;
  name: string;
  email: string;
  photo_url?: string;
  title?: string;
  description?: string;
  birthday?: string;
  age?: number;
  address?: string;
  phone?: string;
  skype?: string;
}

// Router
const router = useRouter()

// Refs and reactive data
const user = ref<User | null>(null)
const avatarUrl = ref<string>('/default-avatar.png')
const previewImage = ref<string>('')
const previewImages = ref<string[]>([])
const showModal = ref(false)
const showUploadModal = ref(false)
const selectedFiles = ref<File[]>([])

const users = ref([])

// Profile form data
const profileData = reactive({
  title: '',
  description: '',
  birthday: '',
  age: 0,
  address: '',
  phone: '',
  skype: ''
})

// Photo album data
const photos = ref<string[]>([])
const isLoading = ref(false) // Track loading state

// Fetch user data and photos on mount
onMounted(() => {
  const storedUser = localStorage.getItem('user')
  if (storedUser) {
    const parsedUser: User = JSON.parse(storedUser)
    user.value = parsedUser
    avatarUrl.value = parsedUser.photo_url || '/default-avatar.png'

    profileData.title = parsedUser.title || ''
    profileData.description = parsedUser.description || ''
    profileData.birthday = parsedUser.birthday || ''
    profileData.age = parsedUser.age || 0
    profileData.address = parsedUser.address || ''
    profileData.phone = parsedUser.phone || ''
    profileData.skype = parsedUser.skype || ''

    // Fetch photos for the user
    fetchPhotos()
  }
})

// Fetch photos for the user
const fetchPhotos = async () => {
  if (user.value?.email) {  // Ensure user email is available
    isLoading.value = true;  // Set loading state to true

    try {
      const encodedEmail = encodeURIComponent(user.value.email);  // Safely encode the email
      const response = await axios.get(`http://localhost:8081/user-posts/${encodedEmail}`);  // Fetch photos for the user

      // Check if the response contains photos
      if (response.data && response.data.photos) {
        photos.value = response.data.photos;  // Store the photos
      } else {
        photos.value = [];  // Set to empty if no photos found
      }
    } catch (err) {
      console.error('Failed to load photos', err);  // Log the error
      alert('Failed to load photos. Please try again later.');  // Notify user of failure
    } finally {
      isLoading.value = false;  // Set loading state to false after completion
    }
  }
};

// Handlers for profile update
const previewImageHandler = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = () => {
      previewImage.value = reader.result as string
    }
    reader.readAsDataURL(file)
  }
}

// Handlers for file uploads
const handleFileChange = (event: Event) => {
  const files = (event.target as HTMLInputElement).files
  if (files) {
    selectedFiles.value = Array.from(files)
    previewImages.value = selectedFiles.value.map(file => URL.createObjectURL(file))
    console.log('Selected files:', selectedFiles.value)
  }
}

// Profile update functionality
const saveProfile = async () => {
  if (!user.value) return

  try {
    const formData = new FormData()
    formData.append('name', user.value.name)
    formData.append('email', user.value.email)
    formData.append('title', profileData.title)
    formData.append('description', profileData.description)
    formData.append('birthday', profileData.birthday)
    formData.append('age', String(profileData.age))
    formData.append('address', profileData.address)
    formData.append('phone', profileData.phone)
    formData.append('skype', profileData.skype)

    const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement
    if (fileInput?.files?.length) {
      formData.append('photo', fileInput.files[0])
    }

    const res = await axios.post(`http://localhost:8081/update-profile/${user.value.id}`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    user.value = res.data.user
    avatarUrl.value = res.data.user.photo_url || '/default-avatar.png'
    localStorage.setItem('user', JSON.stringify(user.value))
    showModal.value = false
    alert('Profile updated successfully!')
  } catch (err) {
    console.error(err)
    alert('Failed to update profile.')
  }
}

// Upload photos to the album
const uploadFiles = async () => {
  if (!user.value?.email) {
    console.error('Error: No user email available')
    alert('You must be logged in to upload files.')
    return
  }

  if (!selectedFiles.value.length) {
    console.error('Error: No files selected for upload')
    alert('Please select at least one file.')
    return
  }

  const formData = new FormData()
  formData.append('email', user.value.email)
  selectedFiles.value.forEach(file => {
    formData.append('photos', file)
  })

  try {
    const response = await axios.post('http://localhost:8081/upload-multiple-photos', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    console.log('Upload success:', response.data)
    alert('Files uploaded successfully!')
    fetchPhotos() // Refresh photos after upload
    closeUploadModal()
  } catch (err) {
    const error = err as any
    console.error('Upload failed:', error.response?.data || error.message)
    alert('Failed to upload. Please try again.')
  }
}

// Close the upload modal
const closeUploadModal = () => {
  showUploadModal.value = false
  previewImages.value = []
}

// Logout confirmation
const confirmLogout = () => {
  if (confirm('Are you sure you want to log out?')) {
    logout()
  }
}

const logout = () => {
  localStorage.removeItem('user')
  router.push('/')
}

// Delete a photo
const deletePhoto = async (photoUrl: string) => {
  if (!user.value?.email) {
    alert('You must be logged in to delete a photo.');
    return;
  }

  const confirmDelete = confirm('Are you sure you want to delete this photo?');
  if (!confirmDelete) return;

  try {
    const params = new URLSearchParams({
      email: user.value.email,
      photo_url: photoUrl,
    });

    const response = await axios.post(`http://localhost:8081/delete-photo?${params.toString()}`);

    if (response.data && response.data.message === 'Photo deleted successfully') {
      photos.value = photos.value.filter((photo) => photo !== photoUrl);
      alert('Photo deleted successfully!');
    } else {
      alert('Failed to delete the photo. Please try again.');
    }
  } catch (err) {
    console.error('Error deleting photo:', err);
    alert('Failed to delete photo. Please try again later.');
  }
};


// Placeholder for active users functionality
const viewActiveUsers = async () => {
  try {
    const response = await fetch('http://localhost:8081/api/active-users');
    const data = await response.json();
    users.value = data;

    const modalEl = document.getElementById('activeUsersModal');
    if (modalEl) {
      const modal = new Modal(modalEl); // Modal is now properly imported
      modal.show();
    }
  } catch (error) {
    console.error('Error fetching active users:', error);
    alert('Failed to load active users.');
  }
}

// Like a post
const likePost = async (photoId: string) => {
  if (!user.value?.id) {
    alert('You must be logged in to like a post.');
    return;
  }

  try {
    const response = await axios.post('http://localhost:8081/api/like', {
      userId: user.value.id,
      photoId: photoId,
    });

    if (response.data.success) {
      alert('Post liked!');
      // Optional: Refresh UI if like counts are dynamic
      viewActiveUsers(); // Reload active users/posts
    } else {
      alert(response.data.message || 'Could not like the post.');
    }
  } catch (error) {
    console.error('Error liking post:', error);
    alert('Failed to like the post. Please try again later.');
  }
};

// Dislike a post
const dislikePost = async (photoId: string) => {
  if (!user.value?.id) {
    alert('You must be logged in to dislike a post.');
    return;
  }

  try {
    const response = await axios.post('http://localhost:8081/api/dislike', {
      userId: user.value.id,
      photoId: photoId,
    });

    if (response.data.success) {
      alert('Post disliked!');
      // Optional: Refresh UI if dislike counts are dynamic
      viewActiveUsers(); // Reload active users/posts
    } else {
      alert(response.data.message || 'Could not dislike the post.');
    }
  } catch (error) {
    console.error('Error disliking post:', error);
    alert('Failed to dislike the post. Please try again later.');
  }
};

</script>


<style scoped>
/* Navbar */
.navbar {
  margin-bottom: 2rem;
}

.navbar-brand {
  font-size: 1.5rem;
  color: #343a40;
}

.btn-outline-danger {
  color: #dc3545;
  border-color: #dc3545;
  transition: all 0.3s ease;
}

.btn-outline-danger:hover {
  background-color: #dc3545;
  color: white;
}

/* Section Styles */
.section {
  padding: 60px 0;
}

.gray-bg {
  background-color: #f7f7f7;
}

.about-section .media {
  margin-bottom: 1rem;
}

.about-section label {
  font-weight: 600;
  display: block;
  color: #555;
}

.about-section p {
  margin: 0;
  font-size: 0.95rem;
  color: #333;
}

/* Avatar */
.about-avatar {
  display: flex;
  justify-content: center;
  align-items: center;
}

.profile-avatar {
  width: 400px;
  height: 400px;
  object-fit: cover;
  border-radius: 50%;
  border: 4px solid #fff;
  box-shadow: 0px 4px 12px rgba(0, 0, 0, 0.3);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.profile-avatar:hover {
  transform: scale(1.05);
  box-shadow: 0px 8px 16px rgba(0, 0, 0, 0.4);
}

/* Buttons */
button {
  transition: all 0.3s ease;
}

button:focus {
  outline: 2px dashed #007bff;
  outline-offset: 2px;
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5); /* dark background */
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  overflow: auto; /* Enable scrolling if needed */
  padding: 1rem;
}

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 10px;
  width: 500px;
  max-width: 90%;
  max-height: 90vh; /* Limit height */
  overflow-y: auto; /* Enable vertical scrolling inside the modal */
  box-shadow: 0 5px 15px rgba(0,0,0,0.3);
}
.modal-content h5 {
  margin-bottom: 20px;
  font-weight: 600;
  font-size: 1.25rem;
}

/* Form Controls */
.form-label {
  font-weight: 500;
  margin-bottom: 0.5rem;
}

.form-control:focus {
  border-color: #007bff;
  box-shadow: 0 0 0 0.2rem rgba(0,123,255,.25);
}

/* Image Preview */
.img-thumbnail {
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Photo Thumbnail Grid */
.photo-thumbnail-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 15px;
  overflow-x: auto;
}

.photo-thumbnail {
  position: relative;
  overflow: hidden;
  cursor: pointer;
}

.photo-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

/* Hover Effect */
.photo-thumbnail:hover img {
  transform: scale(1.05);
}

.photo-thumbnail:hover .btn-danger {
  opacity: 1;
}

.btn-danger {
  opacity: 0;
  transition: opacity 0.3s ease;
}

/* Adjusting the appearance of the delete button */
.btn-danger {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 10;
  background-color: rgba(255, 0, 0, 0.7);
  border: none;
  padding: 5px 10px;
  border-radius: 50%;
}

.btn-danger:hover {
  background-color: rgba(255, 0, 0, 1);
}
.feed-container {
  max-height: 70vh;
  overflow-y: auto;
}

.feed-card {
  border: 1px solid #eaeaea;
  transition: box-shadow 0.3s ease;
}

.feed-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
/* Hover Effect to Enlarge Image */
.photo-hover:hover {
  transform: scale(1.1);
  transition: transform 0.3s ease;
}

/* Styling for Like/Dislike Hand Icons */
.like-dislike-icons {
  display: flex;
  flex-direction: column;
  gap: 5px;
  top: 10px;
  right: 10px;
  z-index: 10;
}

/* Optional: Change cursor when hovering over the image */
.photo-container:hover {
  cursor: pointer;
}
.image-preview {
  transition: transform 0.3s ease;
}

.image-preview:hover {
  transform: scale(1.1);
}


</style>
