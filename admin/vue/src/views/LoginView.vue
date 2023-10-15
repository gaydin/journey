<script setup>
import { Form, Field } from 'vee-validate'
import * as Yup from 'yup'
import axios from 'axios'

// const baseUrl = `${import.meta.env.VITE_API_URL}/users`;
const axiosInstance = axios.create({
  withCredentials: true
})

import useAuthStore from '@/stores/auth.store'
import { XCircleIcon } from '@heroicons/vue/20/solid'

const schema = Yup.object().shape({
  username: Yup.string().required('Username is required'),
  password: Yup.string().required('Password is required')
})

function onSubmit(values, { setErrors }) {
  const authStore = useAuthStore()
  const { username, password } = values

  axiosInstance
    .post('/admin/v1/api/auth/login', {
      login: username,
      password: password
    })
    .then((response) => {
      if (!response) {
        return
      }
      authStore.setAuth()
    })
    .catch(function (error) {
      if (error.response.status === 401) {
        setErrors({ apiError: 'Incorrect login or password' })
      } else {
        setErrors({ apiError: error })
      }
    })
}
</script>

<template>
  <div class="flex min-h-full flex-1 flex-col justify-center px-6 py-12 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-sm">
      <img class="mx-auto h-10 w-auto" src="@/assets/logo.svg" alt="Journey" />
      <h2 class="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">
        Sign in to your account
      </h2>
    </div>

    <div class="mt-10 sm:mx-auto sm:w-full sm:max-w-[480px]">
      <div class="bg-white px-6 py-12 shadow sm:rounded-lg sm:px-12">
        <Form
          v-slot="{ errors, isSubmitting }"
          :validation-schema="schema"
          class="space-y-6"
          @submit="onSubmit"
        >
          <div v-if="errors.apiError" class="rounded-md bg-red-50 p-4">
            <div class="flex">
              <div class="flex-shrink-0">
                <XCircleIcon class="h-5 w-5 text-red-400" aria-hidden="true" />
              </div>
              <div class="ml-3">
                <p class="text-sm font-medium text-red-800">{{ errors.apiError }}</p>
              </div>
            </div>
          </div>
          <div>
            <label for="login" class="block text-sm font-medium leading-6 text-gray-900">
              Login
            </label>
            <div class="mt-2">
              <Field
                id="username"
                name="username"
                type="text"
                default-value="{state.login}"
                class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
              />
            </div>
            <div class="mt-2 text-sm text-red-600">{{ errors.username }}</div>
          </div>

          <div>
            <div class="flex items-center justify-between">
              <label for="password" class="block text-sm font-medium leading-6 text-gray-900">
                Password
              </label>
            </div>
            <div class="mt-2">
              <Field
                id="password"
                name="password"
                type="password"
                auto-complete="current-password"
                class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                :class="{ 'is-invalid': errors.password }"
              />
              <p id="password-error" class="mt-2 text-sm text-red-600">{{ errors.password }}</p>
            </div>
          </div>

          <div>
            <button
              type="submit"
              class="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              :disabled="isSubmitting"
            >
              Sign in
            </button>
          </div>
        </Form>
      </div>
    </div>
  </div>
</template>
