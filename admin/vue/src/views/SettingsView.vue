<script>
import router from '@/router'
import axiosInstance from '@/helpers/axios-wrapper'
import LayoutAdmin from '@/layouts/defaultLayout.vue'
import MiniImage from '@/components/Settings/MiniImage/MiniImage.vue'
import CoverImage from '@/components/Settings/CoverImage/CoverImage.vue'
import NavigationBar from '@/components/Settings/NavigationBar/NavigationBar.vue'
import FormInput from '@/components/Settings/FormInput/FormInput.vue'

export default {
  name: 'App',
  components: { FormInput, LayoutAdmin, MiniImage, CoverImage, NavigationBar },
  data() {
    return {
      blogData: {
        Url: '',
        Title: '',
        Description: '',
        Logo: '',
        Cover: '',
        Themes: [],
        ActiveTheme: '',
        PostsPerPage: 0,
        NavigationItems: []
      },
      disableSubmit: false
    }
  },

  created() {
    axiosInstance
      .get('/admin/v1/api/blog')
      .then((response) => {
        this.blogData = response.data
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
        .patch('/admin/v1/api/blog', this.blogData)
        .then((response) => {
          console.log(response.status, response.data)
        })
        .catch(function (error) {
          console.log(error)
        })
    },
    addNavigationItem() {
      this.blogData.NavigationItems.push({ label: 'Home', url: this.blogData.Url })
    },
    deleteNavigationItem(index) {
      this.blogData.NavigationItems.splice(index, 1)
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
          <form
            @submit="handleSubmit">
          <div class="grid grid-cols-1 gap-x-8 gap-y-8 md:grid-cols-3">
            <div class="px-4 sm:px-0">
              <h2 class="text-base font-semibold leading-7 text-gray-900">Blog settings</h2>
            </div>

            <div
              class="bg-white shadow-sm ring-1 ring-gray-900/5 sm:rounded-xl md:col-span-2"
            >
              <div class="px-4 py-6 sm:p-8">
                <div class="grid max-w-2xl grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
                  <MiniImage id="blog-logo" v-model:modelValue="blogData.Logo" label="Logotype" />
                  <FormInput id="blog-title" v-model:modelValue="blogData.Title" label="Title" />
                  <FormInput
                    id="blog-description"
                    v-model:modelValue="blogData.Description"
                    label="Description"
                  />

                  <CoverImage id="blog-cover" v-model:modelValue="blogData.Cover" label="Cover" />

                  <div class="sm:col-span-4">
                    <label
                      for="blog-postsperpage"
                      class="block text-sm font-medium leading-6 text-gray-900"
                    >
                      Posts per page
                    </label>
                    <div class="mt-2">
                      <div
                        class="flex rounded-md shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600 sm:max-w-md"
                      >
                        <input
                          id="blog-postsperpage"
                          v-model="blogData.PostsPerPage"
                          type="number"
                          class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                        />
                      </div>
                    </div>
                  </div>
                  <div class="sm:col-span-4">
                    <label for="location" class="block text-sm font-medium leading-6 text-gray-900"
                      >Theme</label
                    >
                    <select
                      id="location"
                      v-model="blogData.ActiveTheme"
                      class="mt-2 block w-full rounded-md border-0 py-1.5 pl-3 pr-10 text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-indigo-600 sm:text-sm sm:leading-6"
                    >
                      <option v-for="item in blogData.Themes" :key="item">{{ item }}</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-1 gap-x-8 gap-y-8 pt-10 md:grid-cols-3">
            <div class="px-4 sm:px-0">
              <h2 class="text-base font-semibold leading-7 text-gray-900">Navigation</h2>
            </div>

            <div class="bg-white shadow-sm ring-1 ring-gray-900/5 sm:rounded-xl md:col-span-2">
              <div class="px-4 py-6 sm:p-8">
                <div
                  v-for="(item, index) in blogData.NavigationItems"
                  :key="index"
                  class="grid max-w-4xl grid-cols-9 gap-x-6 gap-y-8"
                >
                  <div class="sm:col-span-4">
                    <div class="relative">
                      <label
                        for="navigation-label"
                        class="absolute -top-2 left-2 inline-block bg-white px-1 text-xs font-medium text-gray-900"
                        >Label</label
                      >
                      <input
                        id="navigation-label"
                        v-model="blogData.NavigationItems[index].label"
                        type="text"
                        class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                      />
                    </div>
                  </div>
                  <div class="sm:col-span-4 ">
                    <div class="relative ">
                      <label
                        for="navigation-url"
                        class="absolute -top-2 left-2 inline-block bg-white px-1 text-xs font-medium text-gray-900"
                        >URL</label
                      >
                      <input
                        id="navigation-url"
                        v-model="blogData.NavigationItems[index].url"
                        type="text"
                        class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                      />
                    </div>
                  </div>

                  <div class="sm:col-span-1">
                    <button
                      type="button"
                      class="rounded-md bg-red-600 px-3 py-3 text-sm font-semibold text-white shadow-sm hover:bg-red-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                      @click="deleteNavigationItem(index)"
                    >
                      Delete
                    </button>

                  </div>
                </div>
              </div>
              <div class="flex items-center justify-end gap-x-6 px-4 py-4 sm:px-8">
                <button
                  type="button"
                  class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                  @click="addNavigationItem"
                >
                  Add
                </button>
              </div>

              <div
                class="flex items-center justify-end gap-x-6 border-t border-gray-900/10 px-4 py-4 sm:px-8"
              >
                <button
                  type="submit"
                  class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                  :disabled="disableSubmit"
                >
                  Save
                </button>
              </div>
            </div>
          </div></form>
        </div>
      </div>
    </main>
  </LayoutAdmin>
</template>
