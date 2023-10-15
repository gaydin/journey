<script setup>
import LayoutAdmin from '@/layouts/defaultLayout.vue'
import MiniImage from '@/components/Settings/MiniImage/MiniImage.vue'
import CoverImage from '@/components/Settings/CoverImage/CoverImage.vue'
import NavigationBar from '@/components/Settings/NavigationBar/NavigationBar.vue'
import FormInput from '@/components/Settings/FormInput/FormInput.vue'
</script>

<script>
import router from '@/router'
import axiosInstance from '@/helpers/axios-wrapper'

export default {
  name: 'App',
  data() {
    return {
      userData: {
        Id: 0,
        Name: '',
        Slug: '',
        Email: '',
        Image: '',
        Cover: '',
        Bio: '',
        Website: '',
        Location: '',
        Password: '',
        PasswordRepeated: ''
      }
    }
  },
  created() {
    axiosInstance
      .get('/admin/v1/api/user')
      .then((response) => {
        this.userData = response.data
      })
      .catch(function (error) {
        if (error.response) {
          if (error.response.status === 401) {
            return router.push('/login')
          }
        }
      })
  },
  methods: {
    handleSubmit(e) {
      e.preventDefault()
      axiosInstance
        .patch('/admin/v1/api/user', this.userData)
        .then((response) => {
          console.log(response.status, response.data)
        })
        .catch(function (error) {
          console.log(error)
        })
    }
  }
}
</script>

<template>
  <LayoutAdmin>
    <NavigationBar />
    <main>
      <div class="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8">
        <div class="space-y-10 divide-y divide-gray-900/10">
          <div class="grid grid-cols-1 gap-x-8 gap-y-8 md:grid-cols-3">
            <div class="px-4 sm:px-0">
              <h2 class="text-base font-semibold leading-7 text-gray-900">
                User {{ userData.Name }}
              </h2>
            </div>

            <form
              class="bg-white shadow-sm ring-1 ring-gray-900/5 sm:rounded-xl md:col-span-2"
              @submit="handleSubmit"
            >
              <div class="px-4 py-6 sm:p-8">
                <div class="grid max-w-2xl grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
                  <MiniImage
                    id="user-avatar"
                    v-model:modelValue="userData.Image"
                    label="Avatar"
                  ></MiniImage>

                  <FormInput id="user-name" v-model:modelValue="userData.Name" label="Name" />
                  <FormInput id="user-slug" v-model:modelValue="userData.Slug" label="Slug" />
                  <CoverImage id="blog-cover" v-model:modelValue="userData.Cover" label="Cover" />
                  <FormInput id="user-email" v-model:modelValue="userData.Email" label="Email" />
                  <FormInput
                    id="user-location"
                    v-model:modelValue="userData.Location"
                    label="Location"
                  />
                  <FormInput
                    id="user-website"
                    v-model:modelValue="userData.Website"
                    label="Website"
                  />
                  <div class="sm:col-span-4">
                    <label for="about" class="block text-sm font-medium leading-6 text-gray-900">
                      Bio
                    </label>
                    <div class="mt-2">
                      <textarea
                        id="about"
                        v-model="userData.Bio"
                        rows="{3}"
                        class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                      />
                    </div>
                  </div>
                  <div class="sm:col-span-4">
                    <label
                      for="user-password"
                      class="block text-sm font-medium leading-6 text-gray-900"
                    >
                      New Password
                    </label>
                    <div class="mt-2">
                      <div
                        class="flex rounded-md shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600 sm:max-w-md"
                      >
                        <input
                          id="user-password"
                          v-model="userData.Password"
                          type="password"
                          class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                        />
                      </div>
                    </div>
                  </div>
                  <div class="sm:col-span-4">
                    <label
                      for="user-password-confirm"
                      class="block text-sm font-medium leading-6 text-gray-900"
                    >
                      Repeat New Password
                    </label>
                    <div class="mt-2">
                      <div
                        class="flex rounded-md shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600 sm:max-w-md"
                      >
                        <input
                          id="user-password-confirm"
                          v-model="userData.PasswordRepeated"
                          type="password"
                          class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                        />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div
                class="flex items-center justify-end gap-x-6 border-t border-gray-900/10 px-4 py-4 sm:px-8"
              >
                <button
                  type="submit"
                  class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                >
                  Save
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </main>
  </LayoutAdmin>
</template>
