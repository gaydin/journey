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
      <input
        :id="id"
        :type="type"
        :value="modelValue"
        class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
        @input="updateInput"
      />
    </div>
  </div>
</template>

<script>
export default {
  name: 'PostInput',
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
    updateInput(event) {
      this.$emit('update:modelValue', event.target.value)
    }
  }
}
</script>
