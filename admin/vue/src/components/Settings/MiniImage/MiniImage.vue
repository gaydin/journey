<template>
  <div class="col-span-full">
    <label v-if="label" :for="id" class="block text-sm font-medium leading-6 text-gray-900">
      {{ label }}
    </label>
    <div class="mt-2">
      <div class="col-span-full flex items-center gap-x-8">
        <img
          v-if="modelValue"
          :src="modelValue"
          :alt="modelValue"
          class="h-24 w-24 flex-none rounded-lg bg-gray-800 object-cover"
        />
        <img
          v-else
          src="@/assets/user-image.jpg"
          alt="@/assets/user-image.jpg"
          class="h-24 w-24 flex-none rounded-lg bg-gray-800 object-cover"
        />
        <div>
          <label
            :for="id"
            class="relative cursor-pointer rounded-md bg-white font-semibold text-indigo-600 focus-within:outline-none focus-within:ring-2 focus-within:ring-indigo-600 focus-within:ring-offset-2 hover:text-indigo-500"
          >
            <span>Upload a file</span>
            <input :id="id" type="file" class="sr-only" @change="upload" />
          </label>
          <p class="mt-2 text-xs leading-5 text-gray-400">JPG, GIF or PNG. 1MB max.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axiosWrapper from '@/helpers/axios-wrapper'

export default {
  name: 'MiniImage',
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
