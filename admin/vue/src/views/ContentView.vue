<script>
import router from '@/router'
import axiosInstance from '@/helpers/axios-wrapper'
import LayoutAdmin from '@/layouts/defaultLayout.vue'
export default {
  name: 'App',
  components: { LayoutAdmin },
  data() {
    return {
      posts: []
    }
  },
  created() {
    axiosInstance
      .get('/admin/v1/api/posts/1')
      .then((response) => {
        this.posts = response.data
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
    truncate: function (str) {
      return str.length > 400 ? str.substring(0, 397) + '...' : str
    },
    openPost: function(id) {
      return router.push("/edit/" + id)
    }
  }
}
</script>

<template>
  <LayoutAdmin>
    <header class="bg-white shadow-sm">
      <div class="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:px-8">
        <h1 class="text-lg font-semibold leading-6 text-gray-900">Dashboard</h1>
      </div>
    </header>
    <main>
      <div class="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8">
        <ul role="list" class="space-y-3">
          <li
            v-for="post in posts"
            :key="post.Id"
            class="overflow-hidden bg-white px-4 py-4 shadow sm:rounded-md sm:px-6 hover:bg-slate-100"
            @click='openPost(post.Id)'
          >
            <div class="flex-auto">
              <div class="flex items-start gap-x-3">
                <router-link
                  class="text-sm font-semibold leading-6 text-gray-900"
                  :to="/edit/ + `${post.Id}`"
                  >{{ post.Title }}</router-link
                >
                <template v-if="post.IsPublished">
                  <p
                    class="text-green-700 bg-green-50 ring-green-600/20 rounded-md whitespace-nowrap mt-0.5 px-1.5 py-0.5 text-xs font-medium ring-1 ring-inset"
                  >
                    Published
                  </p>
                </template>
                <template v-else>
                  <p
                    class="text-yellow-800 bg-yellow-50 ring-yellow-600/20 rounded-md whitespace-nowrap mt-0.5 px-1.5 py-0.5 text-xs font-medium ring-1 ring-inset"
                  >
                    Draft
                  </p>
                </template>
              </div>
              <div class="mt-1 flex items-center gap-x-2 text-xs leading-5 text-gray-500">
                <p class="whitespace-nowrap">
                  Due on <time :datetime="post.Date">{{ post.Date }}</time>
                </p>
              </div>
              <p class="mt-1 line-clamp-2 text-sm leading-6 text-gray-600">
                {{ truncate(post.Markdown) }}
              </p>
            </div>
          </li>
        </ul>
      </div>
    </main>
  </LayoutAdmin>
</template>
