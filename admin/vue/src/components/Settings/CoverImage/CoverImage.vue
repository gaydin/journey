<template>
  <div class="col-span-full">
    <label v-if="label" :for="id" class="block text-sm font-medium leading-6 text-gray-900">
      {{ label }}
    </label>
    <div class="mt-2">
      <label for="image">
        <input id="image" type="file" style="display: none" @change="upload" />
        <img
          v-if="modelValue"
          class="h-64 w-full flex-none rounded-lg bg-gray-800 object-cover"
          :src="modelValue"
          :alt="modelValue"
        />
        <img
          v-else
          id="blog-cover"
          class="img-settings img-thumbnail img-settings"
          src="/src/assets/no-image.png"
          alt="No image"
        />
      </label>
    </div>
  </div>
</template>

<script>
import axiosWrapper from '@/helpers/axios-wrapper'

export default {
  name: 'CoverImage',

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
