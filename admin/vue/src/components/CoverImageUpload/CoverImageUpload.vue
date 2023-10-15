<template>
  <div class="space-y-2 px-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:space-y-0 sm:px-6 sm:py-5">
    <div>
      <label
        v-if="label"
        :for="id"
        class="block text-sm font-medium leading-6 text-gray-900 sm:mt-1.5"
      >
        {{ label }}
      </label>
    </div>

    <div class="sm:col-span-2">
      <template v-if="modelValue">
        <img
          id="post-cover"
          class="img-settings img-thumbnail img-settings"
          :src="modelValue"
          :alt="modelValue"
        />
        <input type="file" accept="image/jpeg" @change="upload" />
      </template>
      <template v-else>
        <div class="col-span-full">
          <label class="block text-sm font-medium leading-6 text-gray-900"> </label>
          <div
            class="mt-2 flex justify-center rounded-lg border border-dashed border-gray-900/25 px-6 py-10"
          >
            <div class="text-center">
              <PhotoIcon class="mx-auto h-12 w-12 text-gray-300" aria-hidden="true" />
              <div class="mt-4 flex text-sm leading-6 text-gray-600">
                <label
                  :for="id"
                  class="relative cursor-pointer rounded-md bg-white font-semibold text-indigo-600 focus-within:outline-none focus-within:ring-2 focus-within:ring-indigo-600 focus-within:ring-offset-2 hover:text-indigo-500"
                >
                  <span>Upload a file</span>
                  <input :id="id" type="file" class="sr-only" @change="upload" />
                </label>
                <p class="pl-1">or drag and drop</p>
              </div>
              <p class="text-xs leading-5 text-gray-600">PNG, JPG, GIF up to 10MB</p>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
import { PhotoIcon } from '@heroicons/vue/24/solid'
import axiosWrapper from '@/helpers/axios-wrapper'

export default {
  name: 'CoverImageUpload',
  components: {
    PhotoIcon
  },
  props: {
    id: {
      type: String,
      default: ''
    },
    label: {
      type: String,
      default: ''
    },
    modelValue: {
      type: [String, Number],
      default: ''
    },
    type: {
      type: String,
      default: 'text'
    }
  },
  emits: ['update:modelValue'],
  methods: {
    upload(event) {
      const image = event.target.files[0]
      let formData = new FormData()
      formData.append('multiplefiles', image, image.name)
      axiosWrapper
        .post('/admin/v1/api/upload', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        .then((response) => {
          this.$emit('update:modelValue', response.data[0])
        })
        .catch(function (reason) {
          console.log('upload error', reason)
        })
    }
  }
}
</script>
