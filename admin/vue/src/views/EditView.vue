<script>
import axios from 'axios'
import router from '@/router'
import { marked } from 'marked'
import { useRoute } from 'vue-router'
import LayoutAdmin from '@/layouts/defaultLayout.vue'
import ContentSlideOver from '@/components/ContentSlideOver/ContentSlideOver.vue'
import { Switch as HeadlessSwitch, SwitchGroup, SwitchLabel } from '@headlessui/vue'
import axiosInstance from '@/helpers/axios-wrapper'
import PostInput from '@/components/PostInput/PostInput.vue'
import CoverImageUpload from '@/components/CoverImageUpload/CoverImageUpload.vue'
import EditModalDelete from '@/components/EditModalDelete/EditModalDelete.vue'
import axiosWrapper from '@/helpers/axios-wrapper'

export default {
  name: 'App',
  components: {
    EditModalDelete,
    CoverImageUpload,
    PostInput,
    ContentSlideOver,
    LayoutAdmin,
    HeadlessSwitch,
    SwitchLabel,
    SwitchGroup
  },
  data() {
    return {
      postData: {
        Id: 0,
        Title: 'New Post',
        Slug: '',
        Markdown: '',
        HTML: '',
        IsFeatured: false,
        IsPage: false,
        IsPublished: false,
        Image: '',
        MetaDescription: '',
        Tags: ''
      },
      markdown: '# Write something!',
      canDelete: false,
      deleteModalOpen: false
    }
  },
  computed: {
    markdownToHtml: function () {
      return marked(this.markdown)
    }
  },
  created() {
    const route = useRoute()
    const id = route.params.id
    if (!id) {
      return
    }

    axios
      .get('/admin/v1/api/post/' + id)
      .then((response) => {
        this.markdown = response.data.Markdown
        this.postData = response.data
        this.canDelete = true
      })
      .catch(function (error) {
        if (error.response) {
          if (error.response.status === 401) {
            return router.push('/login')
          }

          console.log(error.response.data)
          console.log(error.response.status)
        }
      })
  },
  methods: {
    openModal() {
      this.$refs.ContentSlideOver.openModal()
    },
    openModalDelete() {
      this.deleteModalOpen = true
    },
    savePost(e) {
      e.preventDefault()
      let reqBody = {
        Id: this.postData.Id,
        Title: this.postData.Title,
        Slug: this.postData.Slug,
        Markdown: this.markdown,
        IsFeatured: this.postData.IsFeatured,
        IsPage: this.postData.IsPage,
        IsPublished: this.postData.IsPublished,
        Image: this.postData.Image,
        MetaDescription: this.postData.MetaDescription,
        Tags: this.postData.Tags
      }

      if (this.postData.Id === 0) {
        axiosInstance
          .post('/admin/v1/api/post', reqBody)
          .then((response) => {
            console.log(response.status, response.data)
          })
          .catch(function (error) {
            console.log(error)
          })
      } else {
        axiosInstance
          .patch('/admin/v1/api/post', reqBody)
          .then((response) => {
            console.log(response.status, response.data)
          })
          .catch(function (error) {
            console.log(error)
          })
      }
    },
    upload(e) {
      const image = e.target.files[0]
      let formData = new FormData()
      formData.append('multiplefiles', image, image.name)
      axiosWrapper
        .post('/admin/v1/api/upload', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        .then((response) => {
          this.markdown = this.markdown + '\n\n![](' + response.data[0] + ')\n\n';
        })
        .catch(function (reason) {
          console.log('upload error', reason)
        })
    }
  }
}
</script>

<template>
  <LayoutAdmin>
    <header class="bg-white shadow-sm">
      <div class="mx-auto max-w-7xl px-2 sm:px-4 lg:divide-y lg:divide-gray-200 lg:px-8">
        <nav class="hidden lg:flex lg:space-x-8 lg:py-2" aria-label="Global">
          <a
            class="bg-gray-100 text-gray-900 inline-flex items-center rounded-md py-2 px-3 text-sm font-medium"
            @click="openModal"
            >Settings</a
          >
          <a
            class="bg-gray-100 text-gray-900 inline-flex items-center rounded-md py-2 px-3 text-sm font-medium"
            @click="savePost"
            >Save</a
          >
          <label
            for="upload-image"
            class="bg-gray-100 text-gray-900 inline-flex items-center rounded-md py-2 px-3 text-sm font-medium"
          >
            <span>Upload image</span>
            <input id="upload-image" type="file" class="sr-only" @change="upload" />
          </label>
          <a
            v-if="canDelete"
            class="bg-red-100 text-gray-900 inline-flex items-center rounded-md py-2 px-3 text-sm font-medium"
            @click="openModalDelete"
            >Delete</a
          >
          <EditModalDelete v-model:is-open="deleteModalOpen" :post="postData.Id" />
        </nav>

        <ContentSlideOver ref="ContentSlideOver">
          <CoverImageUpload id="cover" v-model:modelValue="postData.Image" label="Cover" />
          <PostInput id="title" v-model:modelValue="postData.Title" type="text" label="Title" />
          <PostInput id="slug" v-model:modelValue="postData.Slug" type="text" label="Custom Slug" />
          <PostInput
            id="meta-description"
            v-model:modelValue="postData.MetaDescription"
            type="text"
            label="Meta Description"
          />
          <PostInput id="tags" v-model:modelValue="postData.Tags" type="text" label="Tags" />

          <div
            class="space-y-2 px-4 sm:grid sm:grid-cols-3 sm:items-start sm:gap-4 sm:space-y-0 sm:px-6 sm:py-5"
          >
            <div>
              <label
                for="meta-description"
                class="block text-sm font-medium leading-6 text-gray-900 sm:mt-1.5"
              >
                Settings
              </label>
            </div>
            <div class="space-y-5 sm:col-span-2">
              <div class="space-y-5 sm:mt-0">
                <div class="relative flex items-start">
                  <SwitchGroup as="div" class="flex items-center justify-between">
                    <span class="flex flex-grow flex-col">
                      <SwitchLabel
                        as="span"
                        class="text-sm font-medium leading-6 text-gray-900"
                        passive
                        >Feature Post</SwitchLabel
                      >
                    </span>
                    <HeadlessSwitch
                      v-model="postData.IsFeatured"
                      :class="[
                        postData.IsFeatured ? 'bg-indigo-600' : 'bg-gray-200',
                        'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2'
                      ]"
                    >
                      <span
                        aria-hidden="true"
                        :class="[
                          postData.IsFeatured ? 'translate-x-5' : 'translate-x-0',
                          'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out'
                        ]"
                      />
                    </HeadlessSwitch>
                  </SwitchGroup>
                </div>
                <div class="relative flex items-start">
                  <SwitchGroup as="div" class="flex items-center justify-between">
                    <span class="flex flex-grow flex-col">
                      <SwitchLabel
                        as="span"
                        class="text-sm font-medium leading-6 text-gray-900"
                        passive
                        >Static Page</SwitchLabel
                      >
                    </span>
                    <HeadlessSwitch
                      v-model="postData.IsPage"
                      :class="[
                        postData.IsPage ? 'bg-indigo-600' : 'bg-gray-200',
                        'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2'
                      ]"
                    >
                      <span
                        aria-hidden="true"
                        :class="[
                          postData.IsPage ? 'translate-x-5' : 'translate-x-0',
                          'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out'
                        ]"
                      />
                    </HeadlessSwitch>
                  </SwitchGroup>
                </div>
                <div class="relative flex items-start">
                  <SwitchGroup as="div" class="flex items-center justify-between">
                    <span class="flex flex-grow flex-col">
                      <SwitchLabel
                        as="span"
                        class="text-sm font-medium leading-6 text-gray-900"
                        passive
                        >Published</SwitchLabel
                      >
                    </span>
                    <HeadlessSwitch
                      v-model="postData.IsPublished"
                      :class="[
                        postData.IsPublished ? 'bg-indigo-600' : 'bg-gray-200',
                        ' relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2'
                      ]"
                    >
                      <span
                        aria-hidden="true"
                        :class="[
                          postData.IsPublished ? 'translate-x-5' : 'translate-x-0',
                          'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out'
                        ]"
                      />
                    </HeadlessSwitch>
                  </SwitchGroup>
                </div>
              </div>
            </div>
          </div>
        </ContentSlideOver>
      </div>
    </header>

    <main class="main">
      <div class="h-screen grid grid-cols-2 gap-2">
        <textarea v-model="markdown" class="block w-full border-0" />

        <div class="prose lg:prose-xl">
          <div v-html="markdownToHtml"></div>
        </div>
      </div>
    </main>
  </LayoutAdmin>
</template>
